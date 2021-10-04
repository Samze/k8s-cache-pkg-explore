package main

import (
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
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

	informerFact := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, time.Minute*20)

	informer := informerFact.ForResource(gvr)

	stopChan := make(chan struct{})
	go informerFact.Start(stopChan)
	informerFact.WaitForCacheSync(stopChan)

	objs, err := informer.Lister().List(labels.NewSelector())
	if err != nil {
		panic(err)
	}
	fmt.Printf("list using created sharedInformer %+v\n", objs)

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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

	//run forever
	<-stopChan
}
