package k8stest

import (
	"context"
	"fmt"
	"time"

	errors "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Deployment is the wrapper over k8s deployment Object
type Deployment struct {
	// kubernetes deployment instance
	object *appsv1.Deployment
}

// Builder enables building an instance of
// deployment
type DeploymentBuilder struct {
	deployment *Deployment // kubernetes deployment instance
	errors     []error
}

// NewBuilder returns a new instance of builder meant for deployment
func NewDeploymentBuilder() *DeploymentBuilder {
	return &DeploymentBuilder{
		deployment: &Deployment{
			object: &appsv1.Deployment{},
		},
	}
}

// WithName sets the Name field of deployment with provided value.
func (b *DeploymentBuilder) WithName(name string) *DeploymentBuilder {
	if len(name) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment: missing name"),
		)
		return b
	}
	b.deployment.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of deployment with provided value.
func (b *DeploymentBuilder) WithNamespace(namespace string) *DeploymentBuilder {
	if len(namespace) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment: missing namespace"),
		)
		return b
	}
	b.deployment.object.Namespace = namespace
	return b
}

// WithLabelsNew resets existing labels if any with
// ones that are provided here
func (b *DeploymentBuilder) WithLabelsNew(labels map[string]string) *DeploymentBuilder {
	if len(labels) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: no new labels"),
		)
		return b
	}

	// copy of original map
	newlbls := map[string]string{}
	for key, value := range labels {
		newlbls[key] = value
	}

	// override
	b.deployment.object.Labels = newlbls
	return b
}

// WithSelectorMatchLabels merges existing matchlabels if any
// with the ones that are provided here
func (b *DeploymentBuilder) WithSelectorMatchLabels(matchlabels map[string]string) *DeploymentBuilder {
	if len(matchlabels) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: missing matchlabels"),
		)
		return b
	}

	if b.deployment.object.Spec.Selector == nil {
		return b.WithSelectorMatchLabelsNew(matchlabels)
	}

	for key, value := range matchlabels {
		b.deployment.object.Spec.Selector.MatchLabels[key] = value
	}
	return b
}

// WithSelectorMatchLabelsNew resets existing matchlabels if any with
// ones that are provided here
func (b *DeploymentBuilder) WithSelectorMatchLabelsNew(matchlabels map[string]string) *DeploymentBuilder {
	if len(matchlabels) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: no new matchlabels"),
		)
		return b
	}

	// copy of original map
	newmatchlabels := map[string]string{}
	for key, value := range matchlabels {
		newmatchlabels[key] = value
	}

	newselector := &metaV1.LabelSelector{
		MatchLabels: newmatchlabels,
	}

	// override
	b.deployment.object.Spec.Selector = newselector
	return b
}

// WithNodeSelector Sets the node selector with the provided argument.
func (b *DeploymentBuilder) WithNodeSelector(selector map[string]string) *DeploymentBuilder {
	if len(selector) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: no node selector"),
		)
		return b
	}
	if b.deployment.object.Spec.Template.Spec.NodeSelector == nil {
		return b.WithNodeSelectorNew(selector)
	}

	for key, value := range selector {
		b.deployment.object.Spec.Template.Spec.NodeSelector[key] = value
	}
	return b
}

// WithNodeSelector Sets the node selector with the provided argument.
func (b *DeploymentBuilder) WithNodeSelectorNew(selector map[string]string) *DeploymentBuilder {
	if len(selector) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: no new node selector"),
		)
		return b
	}

	b.deployment.object.Spec.Template.Spec.NodeSelector = selector
	return b
}

// WithReplicas sets the replica field of deployment
func (b *DeploymentBuilder) WithReplicas(replicas *int32) *DeploymentBuilder {

	if replicas == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment object: nil replicas"),
		)
		return b
	}

	newreplicas := *replicas

	if newreplicas < 0 {
		b.errors = append(
			b.errors,
			errors.Errorf(
				"failed to build deployment object: invalid replicas {%d}",
				newreplicas,
			),
		)
		return b
	}

	b.deployment.object.Spec.Replicas = &newreplicas
	return b
}

// WithPodTemplateSpecBuilder sets the template field of the deployment
func (b *DeploymentBuilder) WithPodTemplateSpecBuilder(
	tmplbuilder *PodtemplatespecBuilder,
) *DeploymentBuilder {
	if tmplbuilder == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build deployment: nil templatespecbuilder"),
		)
		return b
	}

	templatespecObj, err := tmplbuilder.Build()

	if err != nil {
		b.errors = append(
			b.errors,
			errors.Wrap(
				err,
				"failed to build deployment",
			),
		)
		return b
	}

	b.deployment.object.Spec.Template = *templatespecObj.Object
	return b
}

type deploymentBuildOption func(*Deployment)

// NewForAPIObject returns a new instance of builder
// for a given deployment Object
func NewForAPIObject(
	obj *appsv1.Deployment,
	opts ...deploymentBuildOption,
) *Deployment {
	d := &Deployment{object: obj}
	for _, o := range opts {
		o(d)
	}
	return d
}

// Build returns a deployment instance
func (b *DeploymentBuilder) Build() (*appsv1.Deployment, error) {
	err := b.validate()
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to build a deployment: %s",
			b.deployment.object,
		)
	}
	return b.deployment.object, nil
}

func (b *DeploymentBuilder) validate() error {
	if len(b.errors) != 0 {
		return errors.Errorf(
			"failed to validate: build errors were found: %+v",
			b.errors,
		)
	}
	return nil
}

// IsTerminationInProgress checks for older replicas are waiting to
// terminate or not. If Status.Replicas > Status.UpdatedReplicas then
// some of the older replicas are in running state because newer
// replicas are not in running state. It waits for newer replica to
// come into running state then terminate.
func (d *Deployment) IsTerminationInProgress() bool {
	return d.object.Status.Replicas > d.object.Status.UpdatedReplicas
}

// VerifyReplicaStatus verifies whether all the replicas
// of the deployment are up and running
func (d *Deployment) VerifyReplicaStatus() error {
	if d.object.Spec.Replicas == nil {
		return errors.New("failed to verify replica status for deployment: nil replicas")
	}
	if d.object.Status.ReadyReplicas != *d.object.Spec.Replicas {
		return errors.Errorf(d.object.Name+" deployment pods are not in running state expected: %d got: %d",
			*d.object.Spec.Replicas, d.object.Status.ReadyReplicas)
	}
	return nil
}

// IsNotSyncSpec compare generation in status and spec and check if
// deployment spec is synced or not. If Generation <= Status.ObservedGeneration
// then deployment spec is not updated yet.
func (d *Deployment) IsNotSyncSpec() bool {
	return d.object.Generation > d.object.Status.ObservedGeneration
}

// IsUpdateInProgress Checks if all the replicas are updated or not.
// If Status.AvailableReplicas < Status.UpdatedReplicas then all the
// older replicas are not there but there are less number of availableReplicas
func (d *Deployment) IsUpdateInProgress() bool {
	return d.object.Status.AvailableReplicas < d.object.Status.UpdatedReplicas
}

// CreateDeployment creates deployment with provided deployment object
func CreateDeployment(obj *appsv1.Deployment) error {
	deployApi := gTestEnv.KubeInt.AppsV1().Deployments
	_, createErr := deployApi(obj.Namespace).Create(context.TODO(), obj, metaV1.CreateOptions{})
	return createErr
}

// DeleteDeployment deletes the deployment
func DeleteDeployment(name string, namespace string) error {
	deployApi := gTestEnv.KubeInt.AppsV1().Deployments
	err := deployApi(namespace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return err
}

// set the deployment replication to zero until all pods are removed, then reset to the original value.
func BounceDeployment(deploymentName string, namespace string, timeout_seconds int) (int32, error) {
	orig_replicas, err := ZeroDeploymentReplicas(deploymentName, namespace, timeout_seconds)
	if err != nil {
		return orig_replicas, fmt.Errorf("failed to zero replicas, error: %v", err)
	}
	err = RestoreDeploymentReplicas(deploymentName, namespace, timeout_seconds, orig_replicas)
	if err != nil {
		return orig_replicas, fmt.Errorf("failed to restore replicas, error: %v", err)
	}
	return orig_replicas, err
}

func GetDeploymentSpecReplicas(deploymentName string, namespace string) (int32, error) {
	deployment, err := gTestEnv.KubeInt.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metaV1.GetOptions{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get deployment, error: %v", err)
	}
	return *deployment.Spec.Replicas, err
}

func GetDeploymentStatusReplicas(deploymentName string, namespace string) (int32, error) {
	deployment, err := gTestEnv.KubeInt.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get deployment", "error", err)
		return 0, err
	}
	return deployment.Status.Replicas, err
}

// Get deployment object for "deploymentName" deployment
func GetDeployment(deploymentName string, namespace string) (*appsv1.Deployment, error) {
	deployment, err := gTestEnv.KubeInt.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metaV1.GetOptions{},
	)
	if err != nil {
		return nil, err
	}
	return deployment, err
}

// Get pods for "deploymentName" deployment in "namespace"
func GetDeploymentPods(deploymentName string, namespace string) (*coreV1.PodList, error) {
	deployment, err := GetDeployment(deploymentName, namespace)
	if err != nil {
		return nil, err
	}
	// Get the pod template labels for deployment
	labels := deployment.Spec.Template.Labels

	return ListPodsWithLabel(namespace, labels)
}

// Adjust the number of replicas in the deployment
func SetDeploymentReplication(deploymentName string, namespace string, replicas *int32) error {
	depAPI := gTestEnv.KubeInt.AppsV1().Deployments
	var err error

	// this is to cater for a race condition, occasionally seen,
	// when the deployment is changed between Get and Update
	for attempts := 0; attempts < 10; attempts++ {
		deployment, err := depAPI(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
				deploymentName,
				namespace,
				err)
		}
		deployment.Spec.Replicas = replicas
		_, err = depAPI(namespace).Update(context.TODO(), deployment, metaV1.UpdateOptions{})
		if err == nil {
			break
		}
		logf.Log.Info("Re-trying update attempt due to error", "error", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to set replication to deployment, name: %s, namespace: %s, replication: %d, error: %v",
			deploymentName,
			namespace,
			*replicas,
			err)
	}
	return nil
}

func ZeroDeploymentReplicas(deploymentName string, namespace string, timeout_seconds int) (int32, error) {
	orig_replicas, err := GetDeploymentSpecReplicas(deploymentName, namespace)
	if err != nil {
		return orig_replicas, fmt.Errorf("failed to get deployment replicas, error: %v", err)
	}
	if orig_replicas == 0 {
		return orig_replicas, fmt.Errorf("replication unexpectedly zero")
	}

	log.Log.Info("Setting replica count to zero")
	var zero_replicas int32 = 0
	err = SetDeploymentReplication(deploymentName, namespace, &zero_replicas)
	if err != nil {
		return orig_replicas, fmt.Errorf("failed to set deployment replicas, error: %v", err)
	}

	endTime := time.Now().Add(time.Duration(timeout_seconds) * time.Second)
	var replicas int
	for ; time.Now().Before(endTime); time.Sleep(time.Second * 2) {
		pods, err := ListPodsByPrefix(namespace, deploymentName)
		if err != nil {
			err = fmt.Errorf("failed to list pods, error: %v", err)
			return orig_replicas, err
		}
		replicas = len(pods)
		if replicas == 0 {
			break
		}
	}
	if replicas != 0 {
		return orig_replicas, fmt.Errorf("timed out waiting for pods to be zero, podcount: %d", replicas)
	}
	return orig_replicas, err
}

func RestoreDeploymentReplicas(deploymentName string, namespace string, timeout_seconds int, orig_replicas int32) error {
	log.Log.Info("Enabling the original replica count")
	err := SetDeploymentReplication(deploymentName, namespace, &orig_replicas)
	if err != nil {
		return fmt.Errorf("failed to set deployment replicas, error: %v", err)
	}

	var replicas int32
	endTime := time.Now().Add(time.Duration(timeout_seconds) * time.Second)
	for ; time.Now().Before(endTime); time.Sleep(time.Second * 2) {
		replicas, err = GetDeploymentStatusReplicas(deploymentName, namespace)
		if err != nil {
			return fmt.Errorf("failed to get status replicas, error: %v", err)
		}
		if replicas == orig_replicas {
			break
		}
	}
	if replicas != orig_replicas {
		return fmt.Errorf("timed out waiting for pods to be restored, podcount: %d", replicas)
	}
	return err
}

func VerifyDeploymentReadyReplicaCount(name, namespace string, expectedReplicas int) (bool, error) {
	deployment, err := GetDeployment(name, namespace)
	if err != nil {
		return false, err
	}
	logf.Log.Info("Get Deployment", "Ready Replicas", deployment.Status.ReadyReplicas, "Replicas Expected", expectedReplicas)
	if deployment.Status.ReadyReplicas == int32(expectedReplicas) {
		return true, nil
	}
	return false, nil
}
