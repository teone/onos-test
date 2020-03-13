// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"github.com/onosproject/onos-test/pkg/onit/api/resource"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type DaemonSetsReader interface {
	Get(name string) (*DaemonSet, error)
	List() ([]*DaemonSet, error)
}

func NewDaemonSetsReader(client resource.Client) DaemonSetsReader {
	return &daemonSetsReader{
		Client: client,
	}
}

type daemonSetsReader struct {
	resource.Client
}

func (c *daemonSetsReader) Get(name string) (*DaemonSet, error) {
	daemonSet := &appsv1.DaemonSet{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DaemonSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(daemonSet)
	if err != nil {
		return nil, err
	}
	return NewDaemonSet(daemonSet, c.Client), nil
}

func (c *daemonSetsReader) List() ([]*DaemonSet, error) {
	list := &appsv1.DaemonSetList{}
	err := c.Clientset().
		AppsV1().
		RESTClient().
		Get().
		Namespace(c.Namespace()).
		Resource(DaemonSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*DaemonSet, len(list.Items))
	for i, daemonSet := range list.Items {
		results[i] = NewDaemonSet(&daemonSet, c.Client)
	}
	return results, nil
}
