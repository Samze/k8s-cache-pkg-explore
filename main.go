package main

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	res, err := client.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("listing resources %+v", res)
}
