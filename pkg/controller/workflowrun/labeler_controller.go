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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workflowv1alpha1 "github.com/kubevela/workflow/api/v1alpha1"
)

// +kubebuilder:rbac:groups=core.oam.dev,resources=workflowruns,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=core.oam.dev,resources=workflowruns/status,verbs=get

// Setup registers the WorkflowRun labeler with the manager.
func Setup(mgr ctrl.Manager, log logr.Logger) error {
	return (&labeler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		log:    log.WithValues("controller", "workflowrun-labeler"),
	}).SetupWithManager(mgr)
}

type labeler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
}

func (r *labeler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var wr workflowv1alpha1.WorkflowRun
	if err := r.Get(ctx, req.NamespacedName, &wr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	changed := false
	if wr.Labels == nil {
		wr.Labels = map[string]string{}
	}
	if wr.Labels["app.oam.dev/name"] == "" {
		wr.Labels["app.oam.dev/name"] = wr.Name
		changed = true
	}
	if wr.Labels["app.oam.dev/namespace"] == "" {
		wr.Labels["app.oam.dev/namespace"] = wr.Namespace
		changed = true
	}
	if wr.Labels["app.kubernetes.io/managed-by"] == "" {
		wr.Labels["app.kubernetes.io/managed-by"] = "vela-cli"
		changed = true
	}
	if wr.Labels["workflow.oam.dev/source"] == "" {
		wr.Labels["workflow.oam.dev/source"] = "cli"
		changed = true
	}
	if !changed {
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, r.Update(ctx, &wr)
}

func (r *labeler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workflowv1alpha1.WorkflowRun{}).
		Complete(r)
}


