package v1beta2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *DiskPool) DeepCopyInto(out *DiskPool) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

func (in *DiskPool) DeepCopy() *DiskPool {
	if in == nil {
		return nil
	}
	out := new(DiskPool)
	in.DeepCopyInto(out)
	return out
}

func (in *DiskPool) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *DiskPoolList) DeepCopyInto(out *DiskPoolList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DiskPool, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *DiskPoolList) DeepCopy() *DiskPoolList {
	if in == nil {
		return nil
	}
	out := new(DiskPoolList)
	in.DeepCopyInto(out)
	return out
}

func (in *DiskPoolList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
