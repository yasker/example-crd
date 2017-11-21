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
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	messagev1 "github.com/yasker/example-crd/apis/message/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type MessageController struct {
	MessageClient *rest.RESTClient
	MessageScheme *runtime.Scheme
}

// Run starts an Message resource controller
func (c *MessageController) Run(ctx context.Context) error {
	fmt.Print("Watch Message objects\n")

	// Watch Message objects
	_, err := c.watchMessages(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Message resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *MessageController) watchMessages(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.MessageClient,
		messagev1.MessageResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&messagev1.Message{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *MessageController) onAdd(obj interface{}) {
	message := obj.(*messagev1.Message)
	fmt.Printf("[CONTROLLER] OnAdd %s\n", message.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	messageCopy := message.DeepCopy()
	messageCopy.Status = messagev1.MessageStatus{
		State: messagev1.MessageStateBroadcasted,
	}

	err := c.MessageClient.Put().
		Name(message.ObjectMeta.Name).
		Namespace(message.ObjectMeta.Namespace).
		Resource(messagev1.MessageResourcePlural).
		Body(messageCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", messageCopy)
	}
}

func (c *MessageController) onUpdate(oldObj, newObj interface{}) {
	oldMessage := oldObj.(*messagev1.Message)
	newMessage := newObj.(*messagev1.Message)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %s\n", oldMessage.ObjectMeta.SelfLink)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %s\n", newMessage.ObjectMeta.SelfLink)
}

func (c *MessageController) onDelete(obj interface{}) {
	message := obj.(*messagev1.Message)
	fmt.Printf("[CONTROLLER] OnDelete %s\n", message.ObjectMeta.SelfLink)
}
