package worker

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	// kubernetes

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

// HandleList function
func HandleList(c Config) ListResponse {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)

	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	list := fetchTags(c, token, filteredCatalog)

	return list
}

// HandleListRepos function
func HandleListRepos(c Config) Catalog {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(c, token)

	return catalog
}

func getAllPods(clientset *kubernetes.Clientset) (*v1.PodList, error) {

	res, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		return nil, err
	} else if _, isStatus := err.(*errors.StatusError); isStatus {
		return nil, err
	} else if err != nil {
		panic(err.Error())
	}

	return res, nil
}

// HandleListCluster function
func HandleListCluster(c Config) Catalog {

	cat := Catalog{}

	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("could not build kubernetes config: %s", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("could not create clientset: %s", err)
	}

	res, err := getAllPods(clientset)
	if err != nil {
		log.Fatalf("could not get pod: %s", err)
	}

	for _, pod := range res.Items {
		for _, container := range pod.Spec.Containers {

			// add only relevant images to our GCR
			if strings.HasPrefix(container.Image, c.RegistryURL) {
				cat.Repositories = append(cat.Repositories, container.Image)
			}
		}
	}

	return cat
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
