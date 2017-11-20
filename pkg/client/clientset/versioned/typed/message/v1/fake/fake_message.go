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

package fake

import (
	message_v1 "github.com/yasker/example-crd/apis/message/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMessages implements MessageInterface
type FakeMessages struct {
	Fake *FakeMessageV1
	ns   string
}

var messagesResource = schema.GroupVersionResource{Group: "message", Version: "v1", Resource: "messages"}

var messagesKind = schema.GroupVersionKind{Group: "message", Version: "v1", Kind: "Message"}

// Get takes name of the message, and returns the corresponding message object, and an error if there is any.
func (c *FakeMessages) Get(name string, options v1.GetOptions) (result *message_v1.Message, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(messagesResource, c.ns, name), &message_v1.Message{})

	if obj == nil {
		return nil, err
	}
	return obj.(*message_v1.Message), err
}

// List takes label and field selectors, and returns the list of Messages that match those selectors.
func (c *FakeMessages) List(opts v1.ListOptions) (result *message_v1.MessageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(messagesResource, messagesKind, c.ns, opts), &message_v1.MessageList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &message_v1.MessageList{}
	for _, item := range obj.(*message_v1.MessageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested messages.
func (c *FakeMessages) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(messagesResource, c.ns, opts))

}

// Create takes the representation of a message and creates it.  Returns the server's representation of the message, and an error, if there is any.
func (c *FakeMessages) Create(message *message_v1.Message) (result *message_v1.Message, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(messagesResource, c.ns, message), &message_v1.Message{})

	if obj == nil {
		return nil, err
	}
	return obj.(*message_v1.Message), err
}

// Update takes the representation of a message and updates it. Returns the server's representation of the message, and an error, if there is any.
func (c *FakeMessages) Update(message *message_v1.Message) (result *message_v1.Message, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(messagesResource, c.ns, message), &message_v1.Message{})

	if obj == nil {
		return nil, err
	}
	return obj.(*message_v1.Message), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMessages) UpdateStatus(message *message_v1.Message) (*message_v1.Message, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(messagesResource, "status", c.ns, message), &message_v1.Message{})

	if obj == nil {
		return nil, err
	}
	return obj.(*message_v1.Message), err
}

// Delete takes name of the message and deletes it. Returns an error if one occurs.
func (c *FakeMessages) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(messagesResource, c.ns, name), &message_v1.Message{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMessages) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(messagesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &message_v1.MessageList{})
	return err
}

// Patch applies the patch and returns the patched message.
func (c *FakeMessages) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *message_v1.Message, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(messagesResource, c.ns, name, data, subresources...), &message_v1.Message{})

	if obj == nil {
		return nil, err
	}
	return obj.(*message_v1.Message), err
}
