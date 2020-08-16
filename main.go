package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("cannot get homedir")
		}
		os.Exit(1)
	}
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pod, err := clientset.CoreV1().Pods("book").Get("example", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pod)

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(new interface{}) { fmt.Println("Added") },
			UpdateFunc: func(old, new interface{}) { fmt.Println("Updated") },
			DeleteFunc: func(obj interface{}) { fmt.Println("Deleted") },
		},
	)
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	pod, err = podInformer.Lister().Pods("programming-kubernetes").Get("client-go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pod)
}
