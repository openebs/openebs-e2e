package types

type Platform interface {
	PowerOnNode(node string) error
	PowerOffNode(node string) error
	RebootNode(node string) error
	GetNodeStatus(node string) (string, error)
	DetachVolume(volName string, node string) error
	AttachVolume(volName, node string) error
	DetachVolumeFromNode(node string) error
	AttachVolumeToNode(node string) error
	ResizeVolume(volName string, newSizeGB int) error
	// ExtractVolumeIdFromDevicePath parses a provider-specific device path and returns the volume id
	ExtractVolumeIdFromDevicePath(devicePath string) (string, error)
}
