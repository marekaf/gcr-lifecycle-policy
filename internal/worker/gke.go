package worker

import (
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

func fetchImagesFromCluster(c Config) Catalog {

	cat := Catalog{}

	if c.KubeconfigPath == "" {
		// no kubeconfig provided, return empty catalog
		return cat
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.KubeconfigPath)
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

			// add only relevant images from our GCR
			if strings.HasPrefix(container.Image, c.RegistryURL) {

				repo := extractRepositoryFromImage(container.Image)
				cat.Repositories = append(cat.Repositories, repo)
			}
		}
	}

	return cat
}
