package platform

import (
	gcpClient "github.com/openebs/openebs-e2e/common/platform/gcp"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	hcloudClient "github.com/openebs/openebs-e2e/common/platform/hcloud"
	maasClient "github.com/openebs/openebs-e2e/common/platform/maas"
	types "github.com/openebs/openebs-e2e/common/platform/types"
)

func Create() types.Platform {
	cfg := e2e_config.GetConfig()
	switch cfg.Platform.Name {
	case "Hetzner":
		return hcloudClient.New()
	case "Maas":
		return maasClient.New()
	case "Gcp":
		return gcpClient.New()
	}
	return nil
}
