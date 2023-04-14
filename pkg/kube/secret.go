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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ErrSecretNotFound = fmt.Errorf("service account's secret not found")

func GetLastServiceAccountSecrets(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.Secret, error) {
	ss, err := GetServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, ErrSecretNotFound
	}

	fs := ss[0]
	for _, s := range ss[1:] {
		if s.GetCreationTimestamp().Compare(fs.GetCreationTimestamp().Time) >= 0 {
			fs = s
		}
	}

	return &fs, nil
}

func GetServiceAccountSecrets(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) ([]corev1.Secret, error) {
	ss, err := cli.CoreV1().Secrets(namespace).List(ctx, mv1.ListOptions{})
	if err != nil {
		return nil, err
	}

	fss := []corev1.Secret{}
	for _, s := range ss.Items {
		if n, ok := s.Annotations[corev1.ServiceAccountNameKey]; ok && n == name {
			fss = append(fss, s)
		}
	}
	return fss, nil
}

func CreateServiceAccountSecret(ctx context.Context, cli kubernetes.Clientset, name string, namespace string, sa *corev1.ServiceAccount) (*corev1.Secret, error) {
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: sa.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "ServiceAccount",
					Name:       sa.Name,
					UID:        sa.UID,
				},
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	return cli.CoreV1().Secrets(namespace).Create(ctx, s, mv1.CreateOptions{})
}

func DeleteServiceAccountSecret(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) error {
	return cli.CoreV1().Secrets(namespace).Delete(ctx, name, mv1.DeleteOptions{})
}
