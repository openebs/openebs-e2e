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

// LVM
const LvmShared = "shared"
const LvmVgPattern = "vgpattern"
const LvmVgVolGroup = "volgroup"
const LvmThinProvision = "thinProvision"
const LvmStorage = "storage"

// ZFS
const ZfsShared = "shared"
const ZfsRecordSize = "recordsize"
const ZfsCompression = "compression"
const ZfsThinProvision = "thinprovision"
const ZfsDedUp = "dedup"
const ZfsPoolName = "poolname"
const ZfsVolBlockSize = "volblocksize"

//  These variables match the settings used in fsx pod definition

const FsxBlockFileName = "/dev/sdm"

const StrFioCriticalFailure = "fio Critical Failure"
