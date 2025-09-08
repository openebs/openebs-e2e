package common

const NSE2EAgent = "e2e-agent"
const NSE2EPrefix = "e2e-maya"
const NSDefault = "default"
const DefaultIOTimeout = 60
const DefaultVolumeSizeMb = 312
const DefaultFioSizeMb = 250

const SmallClaimSizeMb = 312
const LargeClaimSizeMb = 1024

//  These variables match the settings used in createFioPodDef

const FioFsMountPoint = "/volume"
const FioBlockFilename = "/dev/sdm"

var XFSTestsBlockFilenames = []string{"/dev/test", "/dev/scratch"}

const FioFsFile = "fiotestfile"
const FioFsFilename = FioFsMountPoint + "/" + FioFsFile
const FioFsBlocksPercent = "availblockspercent"
const FioFsBlocksLessBy = "availblockslessby"

// ConfigDir  Relative path to the configuration directory WRT e2e root.
// See common/e2e_config/e2e_config.go

// DefaultConfigFileRelPath  Relative path to default configuration file.
// See common/e2e_config/e2e_config.go

// Storageclass parameter keys
const ScProtocol = "protocol"
const ScFsType = "fsType"
const ScReplicas = "repl"
const ScNodeAffinityTopologyLabel = "nodeAffinityTopologyLabel"
const ScNodeSpreadTopologyKey = "nodeSpreadTopologyKey"
const ScNodeHasTopologyKey = "nodeHasTopologyKey"
const ScPoolAffinityTopologyLabel = "poolAffinityTopologyLabel"
const ScPoolHasTopologyKey = "poolHasTopologyKey"
const ScIOTimeout = "ioTimeout"
const ScNvmeCtrlLossTmo = "nvmeCtrlLossTmo"
const ScThinProvisioning = "thin"
const ScStsAffinityGroup = "stsAffinityGroup"
const ScCloneFsIdAsVolumeId = "cloneFsIdAsVolumeId"
const ScMaxSnapshots = "maxSnapshots"
const ScEncrypted = "encrypted"
const ScPoolClusterSize = "poolClusterSize"

// LVM
const ScLvmShared = "shared"
const ScLvmVgPattern = "vgpattern"
const ScLvmVgVolGroup = "volgroup"
const ScLvmThinProvision = "thinProvision"
const ScLvmStorage = "storage"

// ZFS
const ScZfsShared = "shared"
const ScZfsRecordSize = "recordsize"
const ScZfsCompression = "compression"
const ScZfsThinProvision = "thinProvision"
const ScZfsDeDup = "dedup"
const ScZfsPoolName = "poolname"
const ScZfsVolBlockSize = "volblocksize"
const ScZfsQuotaType = "quotatype"

//  These variables match the settings used in fsx pod definition

const FsxBlockFileName = "/dev/sdm"

const StrFioCriticalFailure = "fio Critical Failure"

// Upgrade constants
type upgradeFlags string

const (
	SkipDataPlaneRestartFlag         upgradeFlags = "--skip-data-plane-restart"
	SkipSingleReplicaValidationFlag  upgradeFlags = "--skip-single-replica-volume-validation"
	SkipReplicaRebuildFlag           upgradeFlags = "--skip-replica-rebuild"
	SkipCordonNodeValidationFlag     upgradeFlags = "--skip-cordoned-node-validation"
	AllowUpgradeToUnstableBranchFlag upgradeFlags = "--allow-unstable"
	SkipUpgradePathValidationFlag    upgradeFlags = "--skip-upgrade-path-validation-for-unsupported-version"
	DisablePartialRebuild            upgradeFlags = "agents.core.rebuild.partial.enabled=false"
)

type userPromptMessages string

const (
	RebuildWarning                           userPromptMessages = "The cluster is rebuilding replica of some volumes"
	SkipSingleReplicaVolumeWarning           userPromptMessages = "These single replica volumes may not be accessible during upgrade"
	SkipSingleReplicaVolumeWarningForOpenEBS userPromptMessages = "These single-replica volumes may not be accessible during upgrade"
	CordonedNodeWarning                      userPromptMessages = "One or more nodes in this cluster are in a Mayastor cordoned state"
)

const (
	VolSizeMb                        = 8192 // in Mb
	DefTimeoutSecs                   = 300  // in seconds
	WaitForRebuildTriggerTimeoutSecs = 60   // in seconds
	UpgradeJobCompletionTimeOutSecs  = 1800 // in seconds
	SmallPollIntervalSecs            = 5    // in seconds
	LargePollIntervalSecs            = 30   // in seconds
	DefRebuildTimeoutSecs            = 600  // in seconds
	SleepTime                        = 3    // in seconds
	ToLocalpvProvisionerImage        = "4.3.0"
)

const (
	VaultNamespace  = "vault"
	VaultStoreName  = "vault-secretstore"
	VaultSecretName = "vault-root-token"
	VaultServer     = "http://vault.vault.svc.cluster.local:8200"
	VaultPath       = "secret"
	VaultKVVersion  = "v2"
)

const (
	HashicorpRepoName     = "hashicorp"
	HashicorpRepoURL      = "https://helm.releases.hashicorp.com"
	HashicorpNamespace    = "vault"
	HashicorpReleaseName  = "vault"
	HashicorpChart        = "hashicorp/vault"
	HashicorpVersion      = "" // empty means latest
	HashicorpDevMode      = true
	HashicorpPodName      = "vault-0"
	HashicorpAppLabel     = "app.kubernetes.io/name=vault"
	HashicorpSecretPath   = "secret"
	HashicorpKVVersion    = "v2"
	HashicorpRootTokenEnv = "VAULT_DEV_ROOT_TOKEN_ID"
)

const (
	SecretStoreGroup      = "external-secrets.io"
	SecretStoreVersion    = "v1"
	SecretStoreResource   = "secretstores"
	SecretStoreAPIVersion = "external-secrets.io/v1"
	SecretStoreKind       = "SecretStore"
	SecretStoreTokenKey   = "token"
)
