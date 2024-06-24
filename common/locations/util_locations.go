package locations

// For now the relative paths are hardcoded, there may be a case to make this
// more generic and data driven.

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

func locationExists(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetBuildInfoFile returns the path to build_info.json if one exists.
// build_info.json is typically part of the install-bundle
func GetBuildInfoFile() (string, error) {
	filePath := path.Clean(e2e_config.GetConfig().MayastorRootDir + "/scripts/../build_info.json")
	_, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	return filePath, err
}

// ParseVersionFile returns unmarshalled data of version.json if one exists.
// version.json is typically part of the install-bundle
func ParseVersionFile() (map[string]string, error) {
	filePath := path.Clean(e2e_config.GetConfig().MayastorRootDir + "/version.json")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	fileContent, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()
	byteResult, _ := io.ReadAll(fileContent)
	version := make(map[string]string)
	err = json.Unmarshal(byteResult, &version)
	return version, err
}

func GetMayastorScriptsDir() (string, error) {
	return locationExists(path.Clean(e2e_config.GetConfig().MayastorRootDir + "/scripts"))
}

func GetControlPlaneScriptsDir() (string, error) {
	return locationExists(path.Clean(e2e_config.GetConfig().MayastorRootDir + "/mcp/scripts"))
}

// GetGeneratedYamlsDir return the path to where Mayastor yaml files are generated this is a generated directory, so may not exist yet.
func GetGeneratedYamlsDir() string {
	return path.Clean(e2e_config.GetConfig().SessionDir + "/install-yamls")
}

// GetControlPlaneGeneratedYamlsDir return the path to where Mayastor yaml files are generated this is a generated directory, so may not exist yet.
func GetControlPlaneGeneratedYamlsDir() string {
	return path.Clean(e2e_config.GetConfig().SessionDir + "/install-yamls-control-plane")
}

// GetE2EAgentPath return the path e2e-agent yaml file
func GetE2EAgentPath() string {
	return path.Clean(e2e_config.GetConfig().E2eRootDir + "/openebs-e2e/tools/e2e-agent")
}

// GetE2EProcyPath return the path e2e-proxy yaml file
func GetE2EProxyPath() string {
	return path.Clean(e2e_config.GetConfig().E2eRootDir + "/openebs-e2e/tools/e2e-proxy")
}

// GetE2EServiceMonitorPath return the path e2e-agent yaml file
func GetE2EServiceMonitorPath() string {
	return path.Clean(e2e_config.GetConfig().E2eRootDir + "/openebs-e2e/configurations")
}

// GetE2EScriptsPath return the path e2e-agent yaml file
func GetE2EScriptsPath() string {
	return path.Clean(e2e_config.GetConfig().E2eRootDir + "/scripts")
}

func GetHelmCrdDir() (string, error) {
	return locationExists(path.Clean(e2e_config.GetConfig().MayastorRootDir + "/extensions/chart/crds"))
}

func GetHelmChartsDir() (string, error) {
	return locationExists(path.Clean(e2e_config.GetConfig().MayastorRootDir + "/extensions/chart"))
}

func GetGeneratedHelmYamlsDir() string {
	return path.Clean(e2e_config.GetConfig().SessionDir + "/charts/generated-yamls")
}

func GenerateLicenseDir() (string, error) {
	// Get the session directory from the configuration
	sessionDir := e2e_config.GetConfig().SessionDir

	// Concatenate with the subdirectory "/license"
	licenseDir := path.Clean(sessionDir + "/license")

	var err error
	// Check if the directory exists
	if _, err := os.Stat(licenseDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.MkdirAll(licenseDir, 0755); err != nil {
			return "", err
		}
	}

	return licenseDir, err

}
