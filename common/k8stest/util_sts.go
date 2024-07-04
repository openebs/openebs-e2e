package k8stest

// This file contains support functions for stateful set on k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"

	errors "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Statefulset is the wrapper over k8s statefulset Object
type Statefulset struct {
	// kubernetes statefulset instance
	object *appsv1.StatefulSet
}

// Builder enables building an instance of
// statefulset
type StatefulsetBuilder struct {
	statefulset *Statefulset // kubernetes statefulset instance
	errors      []error
}

// NewBuilder returns a new instance of builder meant for statefulset
func NewStatefulsetBuilder() *StatefulsetBuilder {
	return &StatefulsetBuilder{
		statefulset: &Statefulset{
			object: &appsv1.StatefulSet{},
		},
	}
}

// WithName sets the Name field of statefulset with provided value.
func (b *StatefulsetBuilder) WithName(name string) *StatefulsetBuilder {
	if len(name) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset: missing name"),
		)
		return b
	}
	b.statefulset.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of statefulset with provided value.
func (b *StatefulsetBuilder) WithNamespace(namespace string) *StatefulsetBuilder {
	if len(namespace) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset: missing namespace"),
		)
		return b
	}
	b.statefulset.object.Namespace = namespace
	return b
}

// WithLabels sets the label field of statefulset with provided value.
func (b *StatefulsetBuilder) WithLabels(labels map[string]string) *StatefulsetBuilder {
	if len(labels) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: no new labels"),
		)
		return b
	}

	// copy of original map
	newLabels := map[string]string{}
	for key, value := range labels {
		newLabels[key] = value
	}

	if b.statefulset.object.Labels != nil {
		for k, v := range labels {
			b.statefulset.object.Labels[k] = v
		}
	} else {
		b.statefulset.object.Labels = newLabels
	}
	return b
}

// WithSelectorMatchLabels set the selector labels field of statefulset with provided value.
func (b *StatefulsetBuilder) WithSelectorMatchLabels(matchlabels map[string]string) *StatefulsetBuilder {
	if len(matchlabels) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: missing matchlabels"),
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

	if b.statefulset.object.Spec.Selector != nil {
		for k, v := range matchlabels {
			b.statefulset.object.Spec.Selector.MatchLabels[k] = v
		}
	} else {
		b.statefulset.object.Spec.Selector = newselector
	}

	return b
}

// WithNodeSelector Sets the node selector with the provided argument.
func (b *StatefulsetBuilder) WithNodeSelector(selector map[string]string) *StatefulsetBuilder {
	if len(selector) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: no node selector"),
		)
		return b
	}

	// copy of original map
	newselector := map[string]string{}
	for key, value := range selector {
		newselector[key] = value
	}

	if b.statefulset.object.Spec.Template.Spec.NodeSelector != nil {
		for k, v := range selector {
			b.statefulset.object.Spec.Template.Spec.NodeSelector[k] = v
		}
	} else {
		b.statefulset.object.Spec.Template.Spec.NodeSelector = newselector
	}

	return b
}

// WithReplicas sets the replica field of statefulset
func (b *StatefulsetBuilder) WithReplicas(replicas *int32) *StatefulsetBuilder {

	if replicas == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: nil replicas"),
		)
		return b
	}

	newreplicas := *replicas

	if newreplicas < 0 {
		b.errors = append(
			b.errors,
			errors.Errorf(
				"failed to build statefulset object: invalid replicas {%d}",
				newreplicas,
			),
		)
		return b
	}

	b.statefulset.object.Spec.Replicas = &newreplicas
	return b
}

// WithPodManagementPolicy sets the spec.podManagementPolicy field of statefulset
func (b *StatefulsetBuilder) WithPodManagementPolicy(policy appsv1.PodManagementPolicyType) *StatefulsetBuilder {

	if len(policy) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: nil pod management policy"),
		)
		return b
	}

	if policy != appsv1.ParallelPodManagement && policy != appsv1.OrderedReadyPodManagement {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset object: wrong value of podManagementPolicy"),
		)
		return b
	}

	b.statefulset.object.Spec.PodManagementPolicy = policy
	return b
}

// WithVolumeClaimTemplate sets the volumeClaimTemplate in statefulset to create volume claims
// according to no of statefulset replicas
func (b *StatefulsetBuilder) WithVolumeClaimTemplate(name, storage, storageClassName string, volType common.VolumeType) *StatefulsetBuilder {
	// Create a new PersistentVolumeClaim object
	volumeClaim := v1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(storage),
				},
			},
			StorageClassName: &storageClassName,
			VolumeMode:       getVolumeMode(volType),
		},
	}

	b.statefulset.object.Spec.VolumeClaimTemplates = append(b.statefulset.object.Spec.VolumeClaimTemplates, volumeClaim)
	return b

}

func getVolumeMode(volType common.VolumeType) *v1.PersistentVolumeMode {
	if volType == common.VolRawBlock {
		mode := v1.PersistentVolumeBlock
		return &mode
	} else {
		mode := v1.PersistentVolumeFilesystem
		return &mode
	}
}

// WithPodTemplateSpecBuilder sets the template field of the statefulset
func (b *StatefulsetBuilder) WithPodTemplateSpecBuilder(
	tmplbuilder *PodtemplatespecBuilder,
) *StatefulsetBuilder {
	if tmplbuilder == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build statefulset: nil templatespecbuilder"),
		)
		return b
	}

	templatespecObj, err := tmplbuilder.Build()

	if err != nil {
		b.errors = append(
			b.errors,
			errors.Wrap(
				err,
				"failed to build statefulset",
			),
		)
		return b
	}

	b.statefulset.object.Spec.Template = *templatespecObj.Object
	return b
}

// Build returns a statefulset instance
func (b *StatefulsetBuilder) Build() (*appsv1.StatefulSet, error) {
	err := b.validate()
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to build a statefulset: %s",
			b.statefulset.object,
		)
	}
	return b.statefulset.object, nil
}

func (b *StatefulsetBuilder) validate() error {
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
func (d *Statefulset) IsTerminationInProgress() bool {
	return d.object.Status.Replicas > d.object.Status.UpdatedReplicas
}

// VerifyReplicaStatus verifies whether all the replicas
// of the statefulset are up and running
func (d *Statefulset) VerifyReplicaStatus() error {
	if d.object.Spec.Replicas == nil {
		return errors.New("failed to verify replica status for statefulset: nil replicas")
	}
	if d.object.Status.ReadyReplicas != *d.object.Spec.Replicas {
		return errors.Errorf(d.object.Name+" statefulset pods are not in running state expected: %d got: %d",
			*d.object.Spec.Replicas, d.object.Status.ReadyReplicas)
	}
	return nil
}

// IsNotSyncSpec compare generation in status and spec and check if
// statefulset spec is synced or not. If Generation <= Status.ObservedGeneration
// then statefulset spec is not updated yet.
func (d *Statefulset) IsNotSyncSpec() bool {
	return d.object.Generation > d.object.Status.ObservedGeneration
}

// IsUpdateInProgress Checks if all the replicas are updated or not.
// If Status.AvailableReplicas < Status.UpdatedReplicas then all the
// older replicas are not there but there are less number of availableReplicas
func (d *Statefulset) IsUpdateInProgress() bool {
	return d.object.Status.AvailableReplicas < d.object.Status.UpdatedReplicas
}

// CreateStatefulset creates statefulset with provided statefulset object
func CreateStatefulset(obj *appsv1.StatefulSet) error {
	stsApi := gTestEnv.KubeInt.AppsV1().StatefulSets
	_, createErr := stsApi(obj.Namespace).Create(context.TODO(), obj, metaV1.CreateOptions{})
	return createErr
}

// DeleteStatefulset deletes the statefulset
func DeleteStatefulset(name string, namespace string) error {
	stsApi := gTestEnv.KubeInt.AppsV1().StatefulSets
	err := stsApi(namespace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return err
}

// Add a node selector to the statefuleset spec and apply
func ApplyNodeSelectorToStatefulset(stsName string, namespace string, label string, value string) error {
	stsApi := gTestEnv.KubeInt.AppsV1().StatefulSets
	sts, err := stsApi(namespace).Get(context.TODO(), stsName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get statefuleset %s : ns: %s : Error: %v", stsName, namespace, err)
	}
	if sts.Spec.Template.Spec.NodeSelector == nil {
		sts.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	sts.Spec.Template.Spec.NodeSelector[label] = value
	_, err = stsApi(namespace).Update(context.TODO(), sts, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply node selector to statefuleset %s : ns: %s : Error: %v", stsName, namespace, err)
	}
	return nil
}

// Adjust the number of replicas in the statefulset
func SetStatefulsetReplication(statefulsetName string, namespace string, replicas *int32) error {
	stsAPI := gTestEnv.KubeInt.AppsV1().StatefulSets
	var err error

	// this is to cater for a race condition, occasionally seen,
	// when the deployment is changed between Get and Update
	for attempts := 0; attempts < 10; attempts++ {
		sts, err := stsAPI(namespace).Get(context.TODO(), statefulsetName, metaV1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get statefulset, name: %s, namespace: %s, error: %v",
				statefulsetName,
				namespace,
				err)
		}
		sts.Spec.Replicas = replicas
		_, err = stsAPI(namespace).Update(context.TODO(), sts, metaV1.UpdateOptions{})
		if err == nil {
			break
		}
		logf.Log.Info("Re-trying update attempt due to error", "error", err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to set replication to deployment, name: %s, namespace: %s, replication: %d, error: %v",
			statefulsetName,
			namespace,
			*replicas,
			err)
	}
	return nil
}

func StatefulSetReady(statefulSetName string, namespace string) bool {
	statefulSet, err := gTestEnv.KubeInt.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		statefulSetName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get daemonset", "error", err)
		return false
	}
	status := statefulSet.Status
	logf.Log.Info("StatefulSet "+statefulSetName, "status", status)
	return status.Replicas == status.ReadyReplicas &&
		status.ReadyReplicas == status.CurrentReplicas && status.ReadyReplicas != 0
}

func GetStsSpecReplicas(statefulSetName string, namespace string) (int32, error) {
	statefulSet, err := gTestEnv.KubeInt.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		statefulSetName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get statefulset", "error", err)
		return 0, err
	}
	return *statefulSet.Spec.Replicas, err
}

func GetStsStatusReplicas(statefulSetName string, namespace string) (int32, error) {
	statefulSet, err := gTestEnv.KubeInt.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		statefulSetName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get statefulset", "error", err)
		return 0, err
	}
	return statefulSet.Status.Replicas, err
}

func StsExists(statefulSetName string, namespace string) (bool, error) {
	statefulSetList, err := gTestEnv.KubeInt.AppsV1().StatefulSets(namespace).List(
		context.TODO(),
		metaV1.ListOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get statefulset list", "error", err)
		return false, err
	}
	for _, sts := range statefulSetList.Items {
		if sts.Name == statefulSetName {
			return true, err
		}
	}
	return false, err
}
