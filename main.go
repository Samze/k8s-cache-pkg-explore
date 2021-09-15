package main

import (
	"context"
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamiclister"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Create client
	restCfg := ctrl.GetConfigOrDie()
	dynamicClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		panic(err)
	}
	gvr := schema.GroupVersionResource{Group: "rabbitmq.com", Version: "v1beta1", Resource: "rabbitmqclusters"}
	client := dynamicClient.Resource(gvr)

	// Create a list/watcher
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return client.List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return client.Watch(context.TODO(), options)
		},
	}

	// Created Shared Informer
	sharedInformer := cache.NewSharedIndexInformer(
		lw,
		&unstructured.Unstructured{},
		time.Minute*20,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	stopChan := make(chan struct{})
	go sharedInformer.Run(stopChan)
	cache.WaitForCacheSync(stopChan, sharedInformer.HasSynced)

	lister := dynamiclister.NewRuntimeObjectShim(dynamiclister.New(sharedInformer.GetIndexer(), gvr))

	objs, err := lister.List(labels.NewSelector())
	if err != nil {
		panic(err)
	}
	fmt.Printf("list using created sharedInformer %+v\n", objs)

	sharedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("event add %+v", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("event update %+v", newObj)

		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("event delete %+v", obj)
		},
	})

	time.Sleep(time.Minute)
}
