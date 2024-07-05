package types

type Platform interface {
	PowerOnNode(node string) error
	PowerOffNode(node string) error
	RebootNode(node string) error
	GetNodeStatus(node string) (string, error)
	DetachVolume(volName string, node string) error
	AttachVolume(volName, node string) error
}
