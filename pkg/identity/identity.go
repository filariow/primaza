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

package identity

import (
	"context"
	"fmt"

	"github.com/primaza/primaza/pkg/kube"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

type Instance struct {
	Namespace      string   `json:"namespace"`
	ServiceAccount string   `json:"serviceAccount"`
	Secrets        []string `json:"secrets"`
}

func Create(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*Instance, error) {
	sa, err := kube.CreateServiceAccountIfNotExists(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	sn := createSecretName(name)
	if _, err := kube.CreateServiceAccountSecret(ctx, cli, sn, namespace, sa); err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, err
	}

	return &Instance{
		Namespace:      namespace,
		ServiceAccount: name,
		Secrets:        []string{sn},
	}, nil
}

func DeleteIfExists(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) error {
	err := kube.DeleteServiceAccount(ctx, cli, name, namespace)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

func createSecretName(sa string) string {
	return fmt.Sprintf("tkn-%s", sa)
}
