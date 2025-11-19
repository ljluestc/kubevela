/*
Copyright 2025 The KubeVela Authors.

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

package workflowrun

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	workflowv1alpha1 "github.com/kubevela/workflow/api/v1alpha1"
)

func buildScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	require.NoError(t, workflowv1alpha1.AddToScheme(s))
	return s
}

func TestLabelerAddsLabels(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	scheme := buildScheme(t)

	wr := &workflowv1alpha1.WorkflowRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core.oam.dev/v1alpha1",
			Kind:       "WorkflowRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wr-cli-sample",
			Namespace: "default",
		},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(wr).
		Build()

	r := &labeler{
		Client: cl,
		Scheme: scheme,
		log:    logr.Discard(),
	}

	_, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
		Namespace: "default",
		Name:      "wr-cli-sample",
	}})
	require.NoError(t, err)

	var got workflowv1alpha1.WorkflowRun
	require.NoError(t, cl.Get(ctx, types.NamespacedName{Namespace: "default", Name: "wr-cli-sample"}, &got))

	require.Equal(t, "wr-cli-sample", got.Labels["app.oam.dev/name"])
	require.Equal(t, "default", got.Labels["app.oam.dev/namespace"])
	require.Equal(t, "vela-cli", got.Labels["app.kubernetes.io/managed-by"])
	require.Equal(t, "cli", got.Labels["workflow.oam.dev/source"])
}

func TestLabelerNoUpdateWhenLabelsPresent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	scheme := buildScheme(t)

	wr := &workflowv1alpha1.WorkflowRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core.oam.dev/v1alpha1",
			Kind:       "WorkflowRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wr-existing",
			Namespace: "default",
			Labels: map[string]string{
				"app.oam.dev/name":               "custom-name",
				"app.oam.dev/namespace":          "custom-ns",
				"app.kubernetes.io/managed-by":   "custom-manager",
				"workflow.oam.dev/source":        "custom-source",
			},
		},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(wr).
		Build()

	r := &labeler{
		Client: cl,
		Scheme: scheme,
		log:    logr.Discard(),
	}

	_, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
		Namespace: "default",
		Name:      "wr-existing",
	}})
	require.NoError(t, err)

	var got workflowv1alpha1.WorkflowRun
	require.NoError(t, cl.Get(ctx, types.NamespacedName{Namespace: "default", Name: "wr-existing"}, &got))

	require.Equal(t, "custom-name", got.Labels["app.oam.dev/name"])
	require.Equal(t, "custom-ns", got.Labels["app.oam.dev/namespace"])
	require.Equal(t, "custom-manager", got.Labels["app.kubernetes.io/managed-by"])
	require.Equal(t, "custom-source", got.Labels["workflow.oam.dev/source"])
}


