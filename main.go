package main

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

func main() {

	// Create client
	restCfg := ctrl.GetConfigOrDie()
	dynamicClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		panic(err)
	}
	gvk := schema.GroupVersionResource{Group: "rabbitmq.com", Version: "v1beta1", Resource: "rabbitmqclusters"}
	client := dynamicClient.Resource(gvk)

	res, err := client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("listing resources %+v", res)

	// Create a list/watcher

	lw := cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return client.List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return client.Watch(context.TODO(), options)
		},
	}

	obj, err := lw.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("listing resources using listwatcher %+v", obj)
}
