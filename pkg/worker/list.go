package worker

import (
	"os"

	log "github.com/sirupsen/logrus"
	// kubernetes

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HandleList function
func HandleList(c Config) ListResponse {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(token)

	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	list := fetchTags(token, filteredCatalog)

	return list
}

// HandleListRepos function
func HandleListRepos(c Config) Catalog {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(token)

	return catalog
}

func getConfigMap(clientset *kubernetes.Clientset, name string, namespace string) (*v1.ConfigMap, error) {

	res, err := clientset.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.WithFields(log.Fields{"configmap": name, "namespace": namespace}).Error("configmap not found")
		return nil, err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.WithFields(log.Fields{"configmap": name, "namespace": namespace, "error": statusError.ErrStatus.Message}).Error("configmap not found")
		return nil, err
	} else if err != nil {
		panic(err.Error())
	}

	log.WithFields(log.Fields{"configmap": name, "namespace": namespace}).Info("configmap found")

	return res, nil
}

// HandleListCluster function
func HandleListCluster(c Config) Catalog {

	token := getToken(c.CredsFile)

	// kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	// configmapName := "deploy-snapshot"
	// namespace := "kube-system"
	// property := "mysql"

	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	log.Fatalf("could not build kubernetes config: %s", err)
	// }

	// // create the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Fatalf("could not create clientset: %s", err)
	// }

	// res, err := getConfigMap(clientset, configmapName, namespace)
	// if err != nil {
	// 	log.Fatalf("could not get configmap: %s", err)
	// }

	// if _, ok := res.Data[property]; !ok {
	// 	log.
	// 		WithField("property", property).
	// 		Fatalf("property does not exist")
	// }

	// fmt.Println(res.Data[property])

	catalog := fetchCatalog(token)

	return catalog
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
