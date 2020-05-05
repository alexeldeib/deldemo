/*
Copyright 2020 Alexander Eldeib.

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

package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	// +kubebuilder:scaffold:scheme
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	kubeconfig, err := ctrl.GetConfig()
	if err != nil {
		fmt.Printf("failed to get kubeconfig: %#+v", err)
		os.Exit(1)
	}

	kubeclient, err := client.New(kubeconfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		fmt.Printf("failed to create kubeclient: %#+v", err)
		os.Exit(1)
	}

	setupLog.Info("deleting specific pod in default namespace")
	err = kubeclient.Delete(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "foo",
		},
	})
	if client.IgnoreNotFound(err) != nil {
		fmt.Printf("failed to create kubeclient: %#+v", err)
		os.Exit(1)
	}

	toBeDeleted := []string{"foo", "bar", "baz"}

	for i := range toBeDeleted {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      toBeDeleted[i],
			},
		}

		setupLog.WithValues("name", pod.Name).Info("deleting pod")

		err = kubeclient.Delete(context.Background(), pod)
		if client.IgnoreNotFound(err) != nil {
			fmt.Printf("failed to create kubeclient: %#+v", err)
			os.Exit(1)
		}
	}

	setupLog.Info("deleting all pods in default namespace")
	err = kubeclient.DeleteAllOf(context.Background(), &corev1.Pod{}, client.InNamespace("default"))
	if err != nil {
		fmt.Printf("failed to create kubeclient: %#+v", err)
		os.Exit(1)
	}
}
