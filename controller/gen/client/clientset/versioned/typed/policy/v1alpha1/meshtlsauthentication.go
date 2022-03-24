/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/linkerd/linkerd2/controller/gen/apis/policy/v1alpha1"
	scheme "github.com/linkerd/linkerd2/controller/gen/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MeshTLSAuthenticationsGetter has a method to return a MeshTLSAuthenticationInterface.
// A group's client should implement this interface.
type MeshTLSAuthenticationsGetter interface {
	MeshTLSAuthentications(namespace string) MeshTLSAuthenticationInterface
}

// MeshTLSAuthenticationInterface has methods to work with MeshTLSAuthentication resources.
type MeshTLSAuthenticationInterface interface {
	Create(ctx context.Context, meshTLSAuthentication *v1alpha1.MeshTLSAuthentication, opts v1.CreateOptions) (*v1alpha1.MeshTLSAuthentication, error)
	Update(ctx context.Context, meshTLSAuthentication *v1alpha1.MeshTLSAuthentication, opts v1.UpdateOptions) (*v1alpha1.MeshTLSAuthentication, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.MeshTLSAuthentication, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.MeshTLSAuthenticationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MeshTLSAuthentication, err error)
	MeshTLSAuthenticationExpansion
}

// meshTLSAuthentications implements MeshTLSAuthenticationInterface
type meshTLSAuthentications struct {
	client rest.Interface
	ns     string
}

// newMeshTLSAuthentications returns a MeshTLSAuthentications
func newMeshTLSAuthentications(c *PolicyV1alpha1Client, namespace string) *meshTLSAuthentications {
	return &meshTLSAuthentications{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the meshTLSAuthentication, and returns the corresponding meshTLSAuthentication object, and an error if there is any.
func (c *meshTLSAuthentications) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MeshTLSAuthentication, err error) {
	result = &v1alpha1.MeshTLSAuthentication{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MeshTLSAuthentications that match those selectors.
func (c *meshTLSAuthentications) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MeshTLSAuthenticationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MeshTLSAuthenticationList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested meshTLSAuthentications.
func (c *meshTLSAuthentications) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a meshTLSAuthentication and creates it.  Returns the server's representation of the meshTLSAuthentication, and an error, if there is any.
func (c *meshTLSAuthentications) Create(ctx context.Context, meshTLSAuthentication *v1alpha1.MeshTLSAuthentication, opts v1.CreateOptions) (result *v1alpha1.MeshTLSAuthentication, err error) {
	result = &v1alpha1.MeshTLSAuthentication{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(meshTLSAuthentication).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a meshTLSAuthentication and updates it. Returns the server's representation of the meshTLSAuthentication, and an error, if there is any.
func (c *meshTLSAuthentications) Update(ctx context.Context, meshTLSAuthentication *v1alpha1.MeshTLSAuthentication, opts v1.UpdateOptions) (result *v1alpha1.MeshTLSAuthentication, err error) {
	result = &v1alpha1.MeshTLSAuthentication{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		Name(meshTLSAuthentication.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(meshTLSAuthentication).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the meshTLSAuthentication and deletes it. Returns an error if one occurs.
func (c *meshTLSAuthentications) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *meshTLSAuthentications) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched meshTLSAuthentication.
func (c *meshTLSAuthentications) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MeshTLSAuthentication, err error) {
	result = &v1alpha1.MeshTLSAuthentication{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("meshtlsauthentications").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
