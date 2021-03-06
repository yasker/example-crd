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

package client

import (
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"

	messagev1 "github.com/yasker/example-crd/apis/message/v1"
)

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := messagev1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &messagev1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}

	return client, scheme, nil
}

func WaitForMessageInstanceProcessed(client *rest.RESTClient, name string) error {
	return wait.Poll(100*time.Millisecond, 10*time.Second, func() (bool, error) {
		var message messagev1.Message
		err := client.Get().
			Resource(messagev1.MessageResourcePlural).
			Namespace(apiv1.NamespaceDefault).
			Name(name).
			Do().Into(&message)

		if err == nil && message.Status.State == messagev1.MessageStateBroadcasted {
			return true, nil
		}

		return false, err
	})
}
