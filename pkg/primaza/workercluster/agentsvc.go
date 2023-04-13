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

package workercluster

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

func DeleteServiceAgent(ctx context.Context, cli *kubernetes.Clientset, namespace string, clusterEnvironment string) error {
	s := runtime.NewScheme()
	if err := appsv1.AddToScheme(s); err != nil {
		return fmt.Errorf("decoder error: %w", err)
	}
	decode := serializer.NewCodecFactory(s).UniversalDeserializer().Decode

	ceDep := strings.ReplaceAll(agentSvcDeployment, "primaza-controller-agentsvc", "pmz-svc-"+clusterEnvironment)
	obj, _, err := decode([]byte(ceDep), nil, nil)
	if err != nil {
		return fmt.Errorf("decoder error: %w", err)
	}

	dep := obj.(*appsv1.Deployment)
	if err := cli.AppsV1().Deployments(namespace).Delete(ctx, dep.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("error deleting deployment: %w", err)
	}

	return nil
}

func PushServiceAgent(ctx context.Context, cli *kubernetes.Clientset, namespace string, clusterEnvironment string) error {
	if _, err := createAgentSvcDeployment(ctx, cli, namespace, clusterEnvironment); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func createAgentSvcDeployment(ctx context.Context, cli *kubernetes.Clientset, namespace string, clusterEnvironment string) (*appsv1.Deployment, error) {
	s := runtime.NewScheme()
	if err := appsv1.AddToScheme(s); err != nil {
		return nil, fmt.Errorf("decoder error: %w", err)
	}
	decode := serializer.NewCodecFactory(s).UniversalDeserializer().Decode

	ceDep := strings.ReplaceAll(agentSvcDeployment, "primaza-controller-agentsvc", "pmz-svc-"+clusterEnvironment)
	obj, _, err := decode([]byte(ceDep), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("decoder error: %w", err)
	}

	dep := obj.(*appsv1.Deployment)
	r, err := cli.AppsV1().Deployments(namespace).Create(ctx, dep, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating deployment: %w", err)
	}
	return r, nil
}

const agentSvcDeployment string = `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/part-of: primaza
    control-plane: primaza-controller-agentsvc
  name: primaza-controller-agentsvc
spec:
  replicas: 1
  selector:
    matchLabels:
      control-plane: primaza-controller-agentsvc
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/default-container: manager
      labels:
        control-plane: primaza-controller-agentsvc
    spec:
      containers:
      - args:
        - --leader-elect
        command:
        - /manager
        env:
        - name: WATCH_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: agentsvc:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        name: manager
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 10m
            memory: 64Mi
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
        volumeMounts:
        - mountPath: /tmp/k8s-webhook-server/serving-certs
          name: cert
          readOnly: true
      securityContext:
        runAsNonRoot: true
      serviceAccountName: primaza-controller-agentsvc
      terminationGracePeriodSeconds: 10
      volumes:
      - name: cert
        secret:
          defaultMode: 420
          secretName: webhook-server-cert
`
