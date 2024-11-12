package e2e_config

import (
	"fmt"
	"os"
	"path"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/ilyakaznacheev/cleanenv"
)

// const ConfigDir = "../../configurations"
// const PlatformConfigDir = "../../configurations/platforms/"
type configurationType int

const (
	baseCfg configurationType = iota
	productCfg
	platformCfg
)

type ConfigurationContext int

const (
	Library ConfigurationContext = iota
	E2eTesting
)

type ProductSpec struct {
	AgentCoreContainerName            string            `yaml:"agentCoreContainerName" env-default:"agent-core"`
	AlertManagerPodPrefix             string            `yaml:"alertManagerPodPrefix" env-default:"alertmanager"`
	ControlPlaneAgent                 string            `yaml:"controlPlaneAgent" env-default:"core-agents"`
	ControlPlaneCoreAgent             string            `yaml:"controlPlaneCoreAgent" env-default:"agent-core"`
	ControlPlaneCsiController         string            `yaml:"controlPlaneCsiController" env-default:"csi-controller"`
	ControlPlaneEtcd                  string            `yaml:"controlPlaneEtcd" env-default:"mayastor-etcd"`
	ControlPlanePoolOperator          string            `yaml:"controlPlanePoolOperator" env-default:"msp-operator"`
	ControlPlaneRestServer            string            `yaml:"controlPlaneRestServer" env-default:"rest"`
	ControlPlaneLocalpvProvisioner    string            `yaml:"controlPlaneLocalpvProvisioner" env-default:"mayastor-localpv-provisioner"`
	ControlPlaneObsCallhome           string            `yaml:"controlPlaneObsCallhome" env-default:"mayastor-obs-callhome"`
	CpuCount                          string            `yaml:"cpuCount" env-default:"2"`
	CrdGroupName                      string            `yaml:"crdGroupName" env-default:"openebs.io"`
	CrdPoolsResourceName              string            `yaml:"crdPoolsResourceName" env-default:"mayastorpools"`
	CsiDaemonsetName                  string            `yaml:"csiDaemonsetName" env-default:"mayastor-csi"`
	CsiNodeServiceAppLabel            string            `yaml:"csiNodeServiceAppLabel" env-default:"csi-node"`
	CsiNodeServiceDaemonset           string            `yaml:"csiNodeServiceDaemonset" env-default:"mayastor-csi-node"`
	CsiNodeContainerName              string            `yaml:"csiNodeContainerName" env-default:"csi-node"`
	CsiProvisioner                    string            `yaml:"csiProvisioner" env-default:"io.openebs.csi-mayastor"`
	DaemonsetName                     string            `yaml:"daemonsetName" env-default:"mayastor"`
	DataPlaneNats                     string            `yaml:"dataPlaneNats" env-default:"nats"`
	DockerOrganisation                string            `yaml:"dockerOrganisation" env-default:"openebs"`
	DockerSecretName                  string            `yaml:"dockerSecretName" env-default:""`
	EngineLabel                       string            `yaml:"engineLabel" env-default:"openebs.io/engine"`
	EngineLabelValue                  string            `yaml:"engineLabelValue" env-default:"mayastor"`
	EtcdYaml                          string            `yaml:"etcdYaml" env-default:"etcd"`
	EventBusNatsSts                   string            `yaml:"eventBusNatsSts" env-default:"mayastor-nats"`
	HaNodeAgentDs                     string            `yaml:"haNodeAgentDs" env-default:"mayastor-agent-ha-node"`
	HaNodeAgentPodPrefix              string            `yaml:"haNodeAgentPodPrefix" env-default:"mayastor-agent-ha-node"`
	HelmReleaseName                   string            `yaml:"helmReleaseName" env-default:"mayastor"`
	OpenEBSHelmReleaseName            string            `yaml:"openEBSHelmReleaseName"`
	IOEnginePodLabelValue             string            `yaml:"ioEnginePodLabelValue" env-default:"io-engine"`
	IOEnginePodName                   string            `yaml:"ioEnginePodName"`
	JaegersCrdName                    string            `yaml:"jaegersCrdName" env-default:"jaegers.jaegertracing.io"`
	KubectlPluginName                 string            `yaml:"kubectlPluginName" env-default:"kubectl-mayastor" env:"e2e_kc_plugin"`
	KubectlPluginPort                 int               `yaml:"kubectlPluginPort" env-default:"30011"`
	LogConfigResources                []string          `yaml:"logConfigResources"`
	LogDumpCsiAttacherName            string            `yaml:"logDumpCsiAttacherName" env-default:"csi-attacher"`
	LogDumpCsiDriverRegistrarName     string            `yaml:"logDumpCsiDriverRegistrarName" env-default:"csi-driver-registrar"`
	LogDumpCsiProvisionerName         string            `yaml:"logDumpCsiProvisionerName" env-default:"csi-provisioner"`
	LogDumpCsiResizerName             string            `yaml:"logDumpCsiResizerName" env-default:"csi-resizer"`
	LogDumpCsiSnapshotControllerName  string            `yaml:"logDumpCsiSnapshotControllerName" env-default:"csi-snapshot-controller"`
	LogDumpCsiSnapshotterName         string            `yaml:"logDumpCsiSnapshotterName" env-default:"csi-snapshotter"`
	LogDumpDirs                       []string          `yaml:"logDumpDirs"`
	LogDumpEngineLabel                string            `yaml:"logDumpEngineLabel"`
	LogDumpHaClusterName              string            `yaml:"logDumpHaClusterName"`
	LogDumpMetricsExporterLabel       string            `yaml:"logDumpMetricsExporterLabel"`
	LoggingLabel                      string            `yaml:"loggingLabel" env-default:"openebs.io/logging"`
	LogLevel                          string            `yaml:"logLevel" env-default:"debug"`
	LokiStatefulset                   string            `yaml:"lokiStatefulset" env-default:"mayastor-loki"`
	MetricsPollingInterval            string            `yaml:"metricsPollingInterval" env-default:"30s"`
	MongoAuthDatabase                 string            `yaml:"mongoAuthDatabase" env-default:"test"`
	MongoAuthPassword                 string            `yaml:"mongoAuthPassword" env-default:"admin123"`
	MongoAuthRootPassword             string            `yaml:"mongoAuthRootPassword" env-default:"r00tAdmin"`
	MongoAuthUsername                 string            `yaml:"mongoAuthUsername" env-default:"admin"`
	MongoDatabasePort                 int               `yaml:"mongoDatabasePort" env-default:"27017"`
	MongoDefaultChartVersion          string            `yaml:"mongoDefaultChartVersion" env-default:"14.5.0"`
	OpenEBSHelmChartName              string            `yaml:"openEBSHelmChartName"`
	OpenEBSHelmRepoName               string            `yaml:"openEBSHelmRepoName"`
	OpenEBSHelmRepoUrl                string            `yaml:"openEBSHelmRepoUrl"`
	MongoHelmRepo                     string            `yaml:"mongoHelmRepo" env-default:"bitnami/mongodb"`
	MongoHelmRepoName                 string            `yaml:"mongoHelmRepoName" env-default:"bitnami"`
	MongoHelmRepoUrl                  string            `yaml:"mongoHelmRepoUrl" env-default:"https://charts.bitnami.com/bitnami"`
	MongoReleaseName                  string            `yaml:"mongoReleaseName" env-default:"ms-mongo"`
	NatsPort                          string            `yaml:"natsPort" env-default:"4222"`
	NvmeControllerModel               string            `yaml:"nvmeControllerModel" env-default:"Mayastor NVMe controller"`
	PartialRebuildCpTimeout           string            `yaml:"partialRebuildCpTimeout" env-default:"600s"`
	PgBenchImage                      string            `yaml:"pgBenchImage" env-default:"postgres:16"`
	PodLabelKey                       string            `yaml:"podLabelKey" env-default:"app"`
	PoolCrdName                       string            `yaml:"poolCrdName" env-default:"mayastorpools.openebs.io"`
	PostgresDatabaseName              string            `yaml:"postgresDatabaseName"`
	PostgresAuthPassword              string            `yaml:"postgresAuthPassword"`
	PostgresAuthRootPassword          string            `yaml:"postgresAuthRootPassword"`
	PostgresAuthUsername              string            `yaml:"postgresAuthUsername"`
	PostgresDatabasePort              int               `yaml:"postgresDatabasePort"`
	PostgresDefaultChartVersion       string            `yaml:"postgresDefaultChartVersion"`
	PostgresHelmRepo                  string            `yaml:"postgresHelmRepo"`
	PostgresHelmRepoName              string            `yaml:"postgresHelmRepoName"`
	PostgresHelmRepoUrl               string            `yaml:"postgresHelmRepoUrl"`
	PostgresReleaseName               string            `yaml:"postgresReleaseName"`
	PostgresK8sLabelName              string            `yaml:"postgresK8sLabelName"`
	PostgresK8sLabelValue             string            `yaml:"postgresK8sLabelValue"`
	ProductName                       string            `yaml:"productName" env-default:"mayastor"`
	ProductNamespace                  string            `yaml:"productNamespace" env-default:"mayastor"`
	OpenEBSProductNamespace           string            `yaml:"openEBSProductNamespace"`
	PrometheusPodPrefix               string            `yaml:"prometheusPodPrefix" env-default:"prometheus"`
	PrometheusPort                    int               `yaml:"prometheusPort" env-default:"30090"`
	PromtailDaemonsetName             string            `yaml:"promtailDaemonsetName" env-default:"mayastor-promtail"`
	RestApiPort                       int               `yaml:"restApiPort" env-default:"30011"`
	RestApiService                    string            `yaml:"restApiService" env-default:"mayastor-api-rest"`
	StatsConfigMapName                string            `yaml:"statsConfigMap" env-default:"mayastor-event-store"`
	StatsDeployment                   string            `yaml:"statsDeployment" env-default:"mayastor-obs-callhome"`
	StatsPort                         string            `yaml:"statsPort" env-default:"9090"`
	StatsService                      string            `yaml:"statsService" env-default:"mayastor-obs-callhome-stats"`
	UpgradePodLabelValue              string            `yaml:"upgradePodLabelValue" env-default:"upgrade"`
	DiskPoolAPIVersionMap             map[string]string `yaml:"diskPoolAPIVersionMap"`
	ChartName                         string            `yaml:"chartName"`
	ChartVersion                      string            `yaml:"chartVersion" env-default:"0.0.0-main"`
	LocalPVContainerName              string            `yaml:"localPVContainerName" env-default:"mayastor-localpv-provisioner"`
	LocalEngineComponentPodLabelKey   string            `yaml:"localEngineComponentPodLabelKey"`
	LvmEngineComponentDsPodLabelValue string            `yaml:"lvmEngineComponentDsPodLabelValue"`
	ZfsEngineComponentDsPodLabelValue string            `yaml:"zfsEngineComponentDsPodLabelValue"`
	LvmEngineDaemonSetName            string            `yaml:"lvmEngineDaemonSetName"`
	ZfsEngineDaemonSetName            string            `yaml:"zfsEngineDaemonSetName"`
	ZfsEngineProvisioner              string            `yaml:"zfsEngineProvisioner"`
	LvmEngineProvisioner              string            `yaml:"lvmEngineProvisioner"`
	HostPathEngineProvisioner         string            `yaml:"hostPathEngineProvisioner"`
	LvmEngineControllerDeploymentName string            `yaml:"lvmEngineControllerDeploymentName"`
	LvmEngineLeaseName                string            `yaml:"lvmEngineLeaseName"`
	LvmEnginePluginContainerName      string            `yaml:"lvmEnginePluginContainerName"`
	LvmEnginePluginDriverName         string            `yaml:"lvmEnginePluginDriverName"`
	UmbrellaOpenebsHelmChartName      string            `yaml:"umbrellaOpenebsHelmChartName"`
	UseUmbrellaOpenEBSChart           bool              `yaml:"useUmbrellaOpenEBSChart" env:"e2e_use_umbrella_openebs_chart" env-default:"false"`
	DisableEngines                    []string          `yaml:"disableEngines"`
	PrometheusNodeExporterServicePort int               `yaml:"prometheusNodeExporterServicePort" env-default:"10100"`
	ZfsEnginePluginContainerName      string            `yaml:"zfsEnginePluginContainerName"`
	ZfsEnginePluginDriverName         string            `yaml:"zfsEnginePluginDriverName"`
	ZfsEngineControllerDeploymentName string            `yaml:"zfsEngineControllerDeploymentName"`
	ZfsEngineLeaseName                string            `yaml:"zfsEngineLeaseName"`
	MayastorCsiControllerLogLevel     string            `yaml:"mayastorCsiControllerLogLevel" env-default:"debug"`
}

// E2EConfig is an application configuration structure
type E2EConfig struct {
	ConfigName  string `yaml:"configName" env-default:"default"`
	ConfigPaths struct {
		ConfigFile         string `yaml:"configFile" env:"e2e_config_file" env-default:""`
		PlatformConfigFile string `yaml:"platformConfigFile" env:"e2e_platform_config_file" env-default:""`
		ProductConfigFile  string `yaml:"productConfigFile" env:"e2e_product_config_file" env-default:""`
	} `yaml:"configPaths"`
	Platform struct {
		// E2ePlatform indicates where the e2e is currently being run from
		Name string `yaml:"name" env-default:"default"`
		// Add HostNetwork: true to the spec of test pods.
		HostNetworkingRequired bool `yaml:"hostNetworkingRequired" env-default:"false"`
		// Some deployments use a different namespace
		MayastorNamespace string `yaml:"mayastorNamespace" env-default:"mayastor"`
		// Some deployments use a different namespace
		FilteredMayastorPodCheck int `yaml:"filteredMayastorPodCheck" env-default:"0"`
		// data plane interface
		IOEngineTargetNvmfIface string `yaml:"ioengineTargetNvmfIface" env-default:"" env:"e2e_ioengine_target_nvmf_iface"`
	} `yaml:"platform"`
	Product ProductSpec `yaml:"product"`

	TestControlNodeLabel              string `yaml:"testControlNodeLabel" env-default:"openebs-test-control"`
	AppLabelControlPlaneCsiController string `yaml:"appLabelControlPlaneCsiController" env-default:"csi-controller"`
	AppLabelControlPlanePoolOperator  string `yaml:"appLabelControlPlanePoolOperator" env-default:"operator-diskpool"`
	AppLabelControlPlaneRestServer    string `yaml:"appLabelControlPlaneRestServer" env-default:"api-rest"`
	AppLabelControlPlaneCoreAgent     string `yaml:"appLabelControlPlaneCoreAgent" env-default:"agent-core"`

	ReactorFreezeDetect            bool `yaml:"reactorFreezeDetect" env:"e2e_reactor_freeze_detect" env-default:"false"`
	SetNexusRebuildVerifyOnInstall bool `yaml:"setNexusRebuildVerifyOnInstall" env-default:"true"`

	// gRPC connection to the mayastor is mandatory for the test run
	// With few exceptions, all CI configurations MUST set this to true
	GrpcMandated bool   `yaml:"grpcMandated" env:"e2e_grpc" env-default:"false"`
	GrpcVersion  string `yaml:"grpcVersion" env:"e2e_grpc_version" env-default:""`
	// Generic configuration files used for CI and automation should not define MayastorRootDir and E2eRootDir
	MayastorRootDir   string `yaml:"mayastorRootDir" env:"e2e_mayastor_root_dir"`
	E2eRootDir        string `yaml:"e2eRootDir" env:"e2e_root_dir"`
	OpenEbsE2eRootDir string `yaml:"openEbsE2eRootDir" env:"openebs_e2e_root_dir"`
	SessionDir        string `yaml:"sessionDir" env:"e2e_session_dir"`
	MayastorVersion   string `yaml:"mayastorVersion" env:"e2e_mayastor_version"`
	KubectlPluginDir  string `yaml:"kubectlPluginDir" env:"e2e_kubectl_plugin_dir"`
	MaasOauthApiToken string `yaml:"maasOauthApiToken" env:"e2e_maas_api_token"`
	MaasEndpoint      string `yaml:"maasEndpoint" env:"e2e_maas_endpoint"`
	ReplicatedEngine  bool   `yaml:"replicatedEngine" env:"replicatedEngine"`

	// Operational parameters
	Cores int `yaml:"cores,omitempty"`
	// Registry from where mayastor images are retrieved
	DockerCache                  string `yaml:"dockercache" env:"e2e_docker_cache" env-default:""`
	ImageTag                     string `yaml:"imageTag" env:"e2e_image_tag"`
	ImagePullPolicy              string `yaml:"imagePullPolicy" env-default:"IfNotPresent" env:"e2e_image_pull_policy"`
	InstallLoki                  bool   `yaml:"installLoki" env-default:"true" env:"install_loki"`
	LokiStatefulsetOnControlNode bool   `yaml:"lokiOnControlNode" env-default:"true" env:"loki_on_control_node"`
	E2eFioImage                  string `yaml:"e2eFioImage" env-default:"openebs/e2e-fio:v3.38-e2e-0" env:"e2e_fio_image"`
	SetSafeMountAlways           bool   `yaml:"setSafeMountAlways" env-default:"false" env:"safe_mount_always"`
	// This is an advisory setting for individual tests
	// If set to true - typically during test development - tests with multiple 'It' clauses should defer asserts till after
	// resources have been cleaned up . This behaviour makes it possible to have useful runs for all 'It' clauses.
	// Typically, set to false for CI test execution - no cleanup after first failure, as a result subsequent 'It' clauses
	// in the test will fail the BeforeEach check, rendering post-mortem checks on the cluster more useful.
	// It may be set to true for when we want maximum test coverage, and post-mortem analysis is a secondary requirement.
	// NOTE: Only some tests support this feature.
	DeferredAssert bool `yaml:"deferredAssert" env-default:"false" env:"e2e_defer_asserts"`
	// TODO: for now using a simple boolean for a specific behaviour suffices, a more sophisticated approach using a policy for test runs may be required.

	// Default replica count, used by tests which do not have a config section.
	DefaultReplicaCount int `yaml:"defaultReplicaCount" env-default:"2" env:"e2e_default_replica_count"`
	// Default provisioning type
	DefaultThinProvisioning bool `yaml:"defaultThinProvisioning" env-default:"false" env:"e2e_default_thin_provisioning"`
	// Restart Mayastor on failure in a prior AfterEach or ResourceCheck
	BeforeEachCheckAndRestart bool `yaml:"beforeEachCheckAndRestart" env-default:"false"`
	// Fail  quickly after failure of a prior AfterEach, overrides BeforeEachCheckAndRestart
	FailQuick bool `yaml:"failQuick" env-default:"false" env:"e2e_fail_quick"`

	// Network interface , HZ: eth0 and GCP: ens4
	NetworkInterface string `yaml:"networkInterface" env-default:"eth0" env:"e2e_default_network_interface"`

	// Run configuration
	ReportsDir string `yaml:"reportsDir" env:"e2e_reports_dir"`
	SelfTest   bool   `yaml:"selfTest" env:"e2e_self_test" env-default:"false"`

	// Boolean value which indicates whether to apply crds or not
	InstallCrds string `yaml:"installCrds" env:"e2e_install_crds" env-default:"false"`

	IOEngineNvmeTimeout int `yaml:"ioEngineNvmeTimeout" env-default:"0"`

	// Individual Test parameters
}

var once sync.Once
var e2eConfig E2EConfig
var configContext ConfigurationContext

// SetContext set execution context to be an e2e ginkgo test run context
// must be called before first invocation of GetConfig to be effective
func SetContext(ctx ConfigurationContext) {
	configContext = ctx
}

// getConfigFilePath given a configuration filename search for its existence
// in the expected locations
//   - current directory
//   - e2e_root_dir (if defined in environment)
//   - openebs_e2e_root_dir (if defined in environment)
//
// Returns the fully qualified path to the file if found.
// Note: if the config filename is a fully qualified path and the file exists
// then that path will be returned
func getConfigFilePath(filename string, cfgType configurationType) (string, error) {
	var subPath = map[configurationType]string{
		baseCfg:     "/configurations",
		productCfg:  "/configurations/product",
		platformCfg: "/configurations/platforms",
	}

	subDir := subPath[cfgType]
	configFile := path.Clean(filename)
	e2eRootDir, haveE2ERootDir := os.LookupEnv("e2e_root_dir")
	openebsE2eRootDir, haveOpenebsE2eRootDir := os.LookupEnv("openebs_e2e_root_dir")
	info, err := os.Stat(configFile)
	if os.IsNotExist(err) && haveE2ERootDir {
		configFile = path.Clean(e2eRootDir + subDir + "/" + filename)
		info, err = os.Stat(configFile)
	}
	if os.IsNotExist(err) && haveOpenebsE2eRootDir {
		configFile = path.Clean(openebsE2eRootDir + subDir + "/" + filename)
		info, err = os.Stat(configFile)
	}
	if err == nil && info.IsDir() {
		err = fmt.Errorf("%v is not a file", configFile)
	}
	return configFile, err
}

// GetConfig This function is called early from junit and various bits have not been initialised yet,
// we cannot use logf or Expect instead we use fmt.Print... and panic.
func GetConfig() E2EConfig {
	once.Do(func() {
		// Initialise the configuration
		_ = cleanenv.ReadEnv(&e2eConfig)
		// We absorb the complexity of locating the configuration file here
		// so that scripts invoking the tests can be simpler
		// - if OS envvar e2e_config is defined
		//		- if it is a path to a file then that file is used as the config file
		//		- else try to use a file of the same name in the configuration directory
		if e2eConfig.ConfigPaths.ConfigFile == "" {
			if configContext == E2eTesting {
				_, _ = fmt.Fprintln(os.Stderr, "Configuration file not specified, using defaults.")
				_, _ = fmt.Fprintln(os.Stderr, "	Use environment variable \"e2e_config_file\" to specify configuration file.")
			}
		} else {
			configFile, err := getConfigFilePath(e2eConfig.ConfigPaths.ConfigFile, baseCfg)
			if err != nil {
				panic(fmt.Sprintf("Unable to access configuration file %v, %v", e2eConfig.ConfigPaths.ConfigFile, err))
			}
			e2eConfig.ConfigPaths.ConfigFile = configFile
			_, _ = fmt.Fprintf(os.Stderr, "Using configuration file %s\n", configFile)
			err = cleanenv.ReadConfig(configFile, &e2eConfig)
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
		}

		if e2eConfig.ConfigPaths.PlatformConfigFile == "" {
			if configContext == E2eTesting {
				_, _ = fmt.Fprintln(os.Stderr, "Platform configuration file not specified, using defaults.")
				_, _ = fmt.Fprintln(os.Stderr, "	Use environment variable \"e2e_platform_config_file\" to specify platform configuration.")
			}
		} else {
			platformCfg, err := getConfigFilePath(e2eConfig.ConfigPaths.PlatformConfigFile, platformCfg)
			if err != nil {
				panic(fmt.Sprintf("Unable to access configuration file %v, %v", e2eConfig.ConfigPaths.PlatformConfigFile, err))
			}
			e2eConfig.ConfigPaths.PlatformConfigFile = platformCfg
			_, _ = fmt.Fprintf(os.Stderr, "Using platform configuration file %s\n", platformCfg)
			err = cleanenv.ReadConfig(platformCfg, &e2eConfig)
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
		}

		if e2eConfig.ConfigPaths.ProductConfigFile == "" {
			if configContext == E2eTesting {
				_, _ = fmt.Fprintln(os.Stderr, "Product configuration file not specified, using defaults.")
				_, _ = fmt.Fprintln(os.Stderr, "	Use environment variable \"e2e_product_config_file\" to specify product configuration.")
			}
		} else {
			productCfg, err := getConfigFilePath(e2eConfig.ConfigPaths.ProductConfigFile, productCfg)
			if err != nil {
				panic(fmt.Sprintf("Unable to access configuration file %v, %v", e2eConfig.ConfigPaths.ProductConfigFile, err))
			}
			e2eConfig.ConfigPaths.ProductConfigFile = productCfg
			_, _ = fmt.Fprintf(os.Stderr, "Using product configuration file %s\n", productCfg)
			err = cleanenv.ReadConfig(productCfg, &e2eConfig)
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
		}

		// MayastorRootDir is either set from the environment variable
		// e2e_mayastor_root_dir or is set in the configuration file.
		if e2eConfig.MayastorRootDir == "" {
			if configContext == E2eTesting {
				_, _ = fmt.Fprintln(os.Stderr, "WARNING: mayastor directory not specified, install and uninstall tests will fail!")
			}
		}

		e2eRootDir, haveE2ERootDir := os.LookupEnv("e2e_root_dir")
		if configContext == E2eTesting {
			artifactsDir := ""
			// if e2e root dir was specified record this in the configuration
			if haveE2ERootDir {
				e2eConfig.E2eRootDir = e2eRootDir
				// and setup the artifacts directory
				artifactsDir = path.Clean(e2eRootDir + "/artifacts")
			} else {
				// use the tmp directory for artifacts
				artifactsDir = path.Clean("/tmp/openebs-e2e")
			}
			_, _ = fmt.Fprintf(os.Stderr, "artifacts directory is %s\n", artifactsDir)

			if e2eConfig.SessionDir == "" {
				// The session directory is required for install and uninstall tests
				// create and use the default one.
				e2eConfig.SessionDir = artifactsDir + "/sessions/default"
				err := os.MkdirAll(e2eConfig.SessionDir, os.ModeDir|os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
			_, _ = fmt.Fprintf(os.Stderr, "session directory is %s\n", e2eConfig.SessionDir)

			if e2eConfig.ReportsDir == "" {
				_, _ = fmt.Fprintln(os.Stderr, "junit report files will not be generated.")
				_, _ = fmt.Fprintln(os.Stderr, "		Use environment variable \"e2e_reports_dir\" to specify a path for the report directory")
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "reports directory is %s\n", e2eConfig.ReportsDir)
			}
			saveConfig()
		} else {
			fmt.Fprintln(os.Stdout, "WARNING: not recording configuration in use!")
		}
		if e2eConfig.SessionDir == "" {
			e2eConfig.SessionDir = "/tmp"
		}
	})

	return e2eConfig
}

func saveConfig() {
	if configContext != E2eTesting {
		return
	}
	cfgBytes, _ := yaml.Marshal(e2eConfig)
	cfgUsedFile := path.Clean(e2eConfig.SessionDir + "/resolved-configuration-" + e2eConfig.ConfigName + "-" + e2eConfig.Platform.Name + ".yaml")
	err := os.WriteFile(cfgUsedFile, cfgBytes, 0644)
	if err == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Resolved config written to %s\n", cfgUsedFile)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Resolved config not written to %s\n%v\n", cfgUsedFile, err)
	}
}

// SetControlPlane sets the control plane configuration if it is unset (i.e. empty) and writes it out if changed.
// If config setting matches  the existing value no action.
// Returns true it the config control plane value matches the input value
func SetControlPlane(controlPlane string) bool {
	_ = GetConfig()
	if e2eConfig.MayastorVersion == "" || e2eConfig.MayastorVersion == controlPlane {
		e2eConfig.MayastorVersion = controlPlane
		saveConfig()
		return true
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to override config control plane from '%s' to '%s'",
			e2eConfig.MayastorVersion, controlPlane)
	}
	return false
}

// GetProductsSpecsMap return a maps of product specs keyed on product name
func GetProductsSpecsMap() map[string]*ProductSpec {
	_ = GetConfig()
	return map[string]*ProductSpec{
		e2eConfig.Product.ProductName: &e2eConfig.Product,
	}
}
