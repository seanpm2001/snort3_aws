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

	configv1 "github.com/snort3_aws/api/v1"
	"github.com/snort3_aws/ipspolicy"
	"github.com/snort3_aws/message"
)

const (
	apiAddr = ":60011"
)

// IpsPolicyReconciler reconciles a IpsPolicy object
type IpsPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=config.github.com,resources=ipspolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.github.com,resources=ipspolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.github.com,resources=ipspolicies/finalizers,verbs=update

func (r *IpsPolicyReconciler) reloadPolicy(policy *configv1.IpsPolicy) error {
	if err := ipspolicy.ValidatePolicyName(policy.Spec.PolicyName); err != nil {
		return errors.New("invalid policy name in crd")
	}
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(apiAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial "+apiAddr)
	}
	defer conn.Close()
	client := message.NewMessageClient(conn)
	ips := message.IpsPolicy{PolicyName: policy.Spec.PolicyName}
	response, err := client.ReloadIpsPolicy(context.Background(), &ips)
	if err != nil {
		return errors.Wrap(err, "failed to reload ips policy")
	}
	if response.Status != "ok" {
		return errors.New("failed to reload ips policy " + response.Status)
	}
	return nil
}

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *IpsPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rLog := log.FromContext(ctx)
	var ipsPolicy configv1.IpsPolicy
	if err := r.Get(ctx, req.NamespacedName, &ipsPolicy); err != nil {
		rLog.Info("No IpsPolicy found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	rLog.Info("Found IpsPolicy", "policy ", ipsPolicy)
	if err := r.reloadPolicy(&ipsPolicy); err != nil {
		rLog.Error(err, "failed to reload ips policy")
		return ctrl.Result{}, err
	}
	rLog.Info("successfully reloaded ips policy")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IpsPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1.IpsPolicy{}).
		Complete(r)
}
