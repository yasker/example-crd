/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.  */

package main

import (
	"context"
	"flag"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	kubeclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	messagev1 "github.com/yasker/example-crd/apis/message/v1"
	"github.com/yasker/example-crd/client"
	messageClientset "github.com/yasker/example-crd/pkg/client/clientset/versioned"

	controller "github.com/yasker/example-crd/controller"
)

func main() {
	masterURL := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use masterURL and kubeconfig if given, otherwise assume in-cluster.
	config, err := clientcmd.BuildConfigFromFlags(*masterURL, *kubeconfig)
	if err != nil {
		panic(err)
	}

	kubeclient, err := kubeclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize custom resource using a CustomResourceDefinition if it does not exist
	crd, err := client.CreateCustomResourceDefinition(kubeclient)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}

	if apierrors.IsAlreadyExists(err) {
		fmt.Println("CRD existed")
	} else {
		fmt.Printf("CRD %v registered\n", crd.ObjectMeta.Name)
	}

	messageClient, messageScheme, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}

	// start a controller on instances of our custom resource
	controller := controller.MessageController{
		MessageClient: messageClient,
		MessageScheme: messageScheme,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	crClient, err := messageClientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	var result *messagev1.Message
	// Create an instance of our custom resource
	firstMessage := &messagev1.Message{
		ObjectMeta: metav1.ObjectMeta{
			Name: "firstmessage",
		},
		Spec: messagev1.MessageSpec{
			Context: "First message",
			Urgent:  false,
		},
		Status: messagev1.MessageStatus{
			State: messagev1.MessageStateCreated,
		},
	}
	result, err = crClient.MessageV1().Messages(apiv1.NamespaceDefault).Create(firstMessage)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	secondMessage := &messagev1.Message{
		ObjectMeta: metav1.ObjectMeta{
			Name: "secondmessage",
		},
		Spec: messagev1.MessageSpec{
			Context: "Second message",
			Urgent:  true,
		},
		Status: messagev1.MessageStatus{
			State: messagev1.MessageStateCreated,
		},
	}
	result, err = crClient.MessageV1().Messages(apiv1.NamespaceDefault).Create(secondMessage)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	// Poll until Message object is handled by controller and gets status
	// updated to "Broadcasted"
	err = client.WaitForMessageInstanceProcessed(messageClient, "firstmessage")
	if err != nil {
		panic(err)
	}
	fmt.Print("BROADCASTED first message!\n")

	err = client.WaitForMessageInstanceProcessed(messageClient, "secondmessage")
	if err != nil {
		panic(err)
	}
	fmt.Print("BROADCASTED second message!\n")

	// Fetch a list of our CRs
	messageList := messagev1.MessageList{}
	err = messageClient.Get().Resource(messagev1.MessageResourcePlural).Do().Into(&messageList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", messageList)
}
