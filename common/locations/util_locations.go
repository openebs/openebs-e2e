package locations

// For now the relative paths are hardcoded, there may be a case to make this
// more generic and data driven.

import (
	"fmt"
	"os"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

var openebsRootDir = e2e_config.GetConfig().OpenEbsE2eRootDir

// GetE2EAgentPath return the path e2e-agent yaml file
func GetE2EAgentPath() string {
	return locationExists(openebsRootDir + "/tools/e2e-agent")
}

// GetE2EProcyPath return the path e2e-proxy yaml file
func GetE2EProxyPath() string {
	return locationExists(openebsRootDir + "/tools/e2e-proxy")
}

// GetE2EServiceMonitorPath return the path e2e-agent yaml file
func GetE2EServiceMonitorPath() string {
	return locationExists(openebsRootDir + "/configurations")
}

// GetE2EScriptsPath return the path e2e-agent yaml file
func GetE2EScriptsPath() string {
	return locationExists(openebsRootDir + "/scripts")
}

func locationExists(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("directory %s not , error: %v", path, err)
		panic("Error: directory not found")
	}
	return path
}
