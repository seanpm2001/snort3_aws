/*
Copyright 2021.

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

package controllers

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	policyv1 "github.com/snort3_aws/api/v1"
	"github.com/snort3_aws/message"
)

// TalosSpdReconciler reconciles a TalosSpd object
type TalosSpdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=policy.github.com,resources=talosspds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy.github.com,resources=talosspds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=policy.github.com,resources=talosspds/finalizers,verbs=update

func (r *TalosSpdReconciler) reloadSpd(lspConfig *policyv1.TalosSpd) error {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(apiAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial "+apiAddr)
	}
	defer conn.Close()
	client := message.NewMessageClient(conn)
	lsp := message.ReloadLsp{LspVersion: lspConfig.Spec.Version}
	response, err := client.ReloadTalosLsp(context.Background(), &lsp)
	if err != nil {
		return errors.Wrap(err, "failed to reload talos lsp")
	}
	if response.Status != "ok" {
		return errors.New("failed to reload talos lsp " + response.Status)
	}
	return nil
}

func (r *TalosSpdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rLog := log.FromContext(ctx)

	var spd policyv1.TalosSpd
	if err := r.Get(ctx, req.NamespacedName, &spd); err != nil {
		rLog.Info("No spd found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	rLog.Info("Found new talos lightspd version", "version", spd)
	if err := r.reloadSpd(&spd); err != nil {
		rLog.Error(err, "failed to reload talos lsp")
		return ctrl.Result{}, err
	}
	rLog.Info("successfully reloaded talos lsp")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TalosSpdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyv1.TalosSpd{}).
		Complete(r)
}
