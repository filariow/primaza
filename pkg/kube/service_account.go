/*
Copyright 2023 The Primaza Authors.

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

package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateServiceAccountIfNotExists(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.ServiceAccount, error) {
	c := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	o := mv1.CreateOptions{}
	sa, err := cli.CoreV1().ServiceAccounts(namespace).Create(ctx, c, o)
	if err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}

		return GetServiceAccount(ctx, cli, name, namespace)
	}

	return sa, nil
}

func GetServiceAccount(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.ServiceAccount, error) {
	return cli.CoreV1().ServiceAccounts(namespace).Get(ctx, name, mv1.GetOptions{})
}

func DeleteServiceAccount(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) error {
	return cli.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, mv1.DeleteOptions{})
}
