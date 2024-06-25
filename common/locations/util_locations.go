package locations

// For now the relative paths are hardcoded, there may be a case to make this
// more generic and data driven.

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

// GetE2EAgentPath return the path e2e-agent yaml file
func GetE2EAgentPath() string {
	return GetAbsolutePath("../../tools/e2e-agent")
}

// GetE2EProcyPath return the path e2e-proxy yaml file
func GetE2EProxyPath() string {
	return GetAbsolutePath("../../tools/e2e-proxy")
}

// GetE2EServiceMonitorPath return the path e2e-agent yaml file
func GetE2EServiceMonitorPath() string {
	return GetAbsolutePath("../../configurations")
}

// GetE2EScriptsPath return the path e2e-agent yaml file
func GetE2EScriptsPath() string {
	return path.Clean(e2e_config.GetConfig().E2eRootDir + "/scripts")
}

func GetAbsolutePath(relativePath string) string {
	// Get the directory of the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error: unable to get current file location")
		panic("Error: unable to get current file location")
	}

	currentDir := filepath.Dir(currentFile)

	// Construct the full path relative to the current file's directory
	configFile := filepath.Join(currentDir, relativePath)

	// Clean the path
	configFile = filepath.Clean(configFile)
	fmt.Println("Cleaned Path:", configFile)

	// Convert to absolute path
	absConfigFile, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Println("Error converting to absolute path:", err)
		panic("Error converting to absolute path")
	}
	fmt.Println("Absolute Path:", absConfigFile)

	return absConfigFile
}
