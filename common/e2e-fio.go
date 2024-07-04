package common

import (
	"fmt"
	"strings"
)

type FioAppArgsSet int

var fioRandReadWriteParams = []string{
	"--random_generator=tausworthe64",
	"--rw=randrw",
	"--ioengine=libaio",
	"--iodepth=16",
	"--verify_fatal=1",
	"--verify=crc32",
	"--verify_async=2",
}

var fioDefaultParams = fioRandReadWriteParams

var fioRandWriteParams = []string{
	"--random_generator=tausworthe64",
	"--rw=randwrite",
	"--ioengine=libaio",
	"--iodepth=16",
	"--verify_fatal=1",
	"--verify=crc32",
	"--verify_async=2",
	"--verify_pattern=%o",
}

var fioRandReadParams = []string{
	"--random_generator=tausworthe64",
	"--rw=randread",
	"--ioengine=libaio",
	"--iodepth=16",
	"--verify_fatal=1",
	"--verify=crc32",
	"--verify_async=2",
}

var fioPerfSeqReadParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=read",
}

var fioPerfSeqWriteParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=write",
}

var fioPerfSeqMixedParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=rw",
}

var fioPerfRandReadParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=randread",
}

var fioPerfRandWriteParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=randwrite",
}

var fioPerfRandMixedParams = []string{
	"--numjobs=1",
	"--ioengine=libaio",
	"--iodepth=16",
	"--rw=randrw",
}

const (
	DefaultFioArgs FioAppArgsSet = iota
	RandWriteFioArgs
	RandReadFioArgs
	CustomFioArgs
	NoFioArgs
	PerfSeqReadFioArgs
	PerfSeqWriteFioArgs
	PerfSeqMixedFioArgs
	PerfRandReadFioArgs
	PerfRandWriteFioArgs
	PerfRandMixedFioArgs
	RandReadWriteFioArgs
)

func (f FioAppArgsSet) String() string {
	switch f {
	case NoFioArgs:
		return "None"
	case DefaultFioArgs:
		return "Default"
	case RandWriteFioArgs:
		return "RandWrite"
	case RandReadFioArgs:
		return "RandRead"
	case CustomFioArgs:
		return "CustomFioArgs"
	case PerfSeqReadFioArgs:
		return "PerfSeqRead"
	case PerfSeqWriteFioArgs:
		return "PerfSeqWrite"
	case PerfSeqMixedFioArgs:
		return "PerfSeqMixed"
	case PerfRandReadFioArgs:
		return "PerfRandRead"
	case PerfRandWriteFioArgs:
		return "PerfRandWrite"
	case PerfRandMixedFioArgs:
		return "PerfRandMixed"
	case RandReadWriteFioArgs:
		return "RandReadWrite"
	default:
		return "Unknown"
	}
}

type FioFsAllocation int

const (
	FioFsAllocDefault      FioFsAllocation = iota
	FioFsAllocLessByBlocks FioFsAllocation = iota
	FioFsAllocPercentage   FioFsAllocation = iota
)

type fioTarget struct {
	// fio "filename"
	target string
	// device used for e2e-fio liveness and makefile
	device string
	// filename used for e2e-fio makefile
	fsFile string
	// required if only target is set
	targetSize uint64
	// fs file allocation strategy
	fsAllocStrategy FioFsAllocation
	// one of
	//  - number of blocks less the available reported by the file system
	//  - percentage of available bytes reported by the file system to use for fio file.
	fsAllocUnit uint
}

type FioTargetSizeRecord struct {
	Path *string `json:"path"`
	Size *uint64 `json:"fio_target_size"`
}
type FioExitRecord struct {
	ExitValue   *int    `json:"exit_value"`
	ElapsedSecs *uint64 `json:"elapsed_seconds"`
}

type FioJsonRecords struct {
	ExitValues  []FioExitRecord
	TargetSizes []FioTargetSizeRecord
}

type E2eFioPodLogSynopsis struct {
	Err             error
	CriticalFailure bool
	Text            []string
	JsonRecords     FioJsonRecords
}

func (s *E2eFioPodLogSynopsis) String() string {
	cfStr := ""
	if s.CriticalFailure {
		cfStr = StrFioCriticalFailure + "\n"
	}
	return fmt.Sprintf("%s%s\n",
		cfStr,
		strings.Join(s.Text, "\n"),
	)
}

type E2eFioPodOutputMonitor struct {
	Completed bool
	Synopsis  E2eFioPodLogSynopsis
}

type E2eFioArgsBuilder struct {
	err                  error
	customArgs           []string
	additionalArgs       []string
	targets              []fioTarget
	duration             int
	liveness             bool
	livenessReadInterval uint
	livenessTimeout      uint
	exitValue            int
	zerofill             bool
	argsSet              FioAppArgsSet
	blockSize            uint
	loops                uint
	// declare inverse of direct, we want direct to be defaulted to ON
	indirectIO bool
}

// NewE2eFioArgsBuilder returns an instance of e2e fio args builder
func NewE2eFioArgsBuilder() *E2eFioArgsBuilder {
	obj := E2eFioArgsBuilder{
		livenessReadInterval: 1,
		livenessTimeout:      60,
		blockSize:            4096,
	}
	return &obj
}

func (e *E2eFioArgsBuilder) WithArgumentSet(set FioAppArgsSet) *E2eFioArgsBuilder {
	if e.argsSet == DefaultFioArgs {
		e.argsSet = set
	} else {
		e.err = fmt.Errorf("overwriting previously set argument %s,%s; %v", e.argsSet, set, e.err)
	}
	return e
}

// WithDefaultArgs use "standard" fio arguments
func (e *E2eFioArgsBuilder) WithDefaultArgs() *E2eFioArgsBuilder {
	return e.WithArgumentSet(DefaultFioArgs)
}

// WithRandWrite only writes to all blocks of target volume
func (e *E2eFioArgsBuilder) WithRandWrite() *E2eFioArgsBuilder {
	return e.WithArgumentSet(RandWriteFioArgs)
}

// WithRandRead verifies data written with previous RandWrite
func (e *E2eFioArgsBuilder) WithRandRead() *E2eFioArgsBuilder {
	return e.WithArgumentSet(RandReadFioArgs)
}

// WithRandReadWrite random reads and writes of target volume
func (e *E2eFioArgsBuilder) WithRandReadWrite() *E2eFioArgsBuilder {
	return e.WithArgumentSet(RandReadWriteFioArgs)
}

// WithCustomArgs use custom fio arguments
func (e *E2eFioArgsBuilder) WithCustomArgs(customArgs []string) *E2eFioArgsBuilder {
	e.argsSet = CustomFioArgs
	e.customArgs = customArgs[:]
	return e
}

// WithAdditionalArgs use additional fio arguments
func (e *E2eFioArgsBuilder) WithAdditionalArgs(args []string) *E2eFioArgsBuilder {
	e.additionalArgs = append(e.additionalArgs, args[:]...)
	return e
}

// WithAdditionalArg use an additional fio arguments
func (e *E2eFioArgsBuilder) WithAdditionalArg(arg string) *E2eFioArgsBuilder {
	e.additionalArgs = append(e.additionalArgs, arg)
	return e
}

// WithFsFile add a fio target on a filesystem
// the target file is created explicitly by e2e-fio before fio is launched
func (e *E2eFioArgsBuilder) WithFsFile(fsMount string, fileRelPath string) *E2eFioArgsBuilder {
	// defaults to blockslessby 1
	// TODO: add facility to choose %blocks, blockslessby or size in bytes
	e.targets = append(e.targets, fioTarget{
		target: fsMount + "/" + fileRelPath,
		device: fsMount,
		fsFile: fileRelPath,
	})
	return e
}

// WithDefaultFile add the "standard" fio target on a filesystem
// the target file is created explicitly by e2e-fio before fio is launched
func (e *E2eFioArgsBuilder) WithDefaultFile() *E2eFioArgsBuilder {
	// defaults to blockslessby 1
	// TODO: add facility to choose %blocks, blockslessby or size in bytes
	return e.WithFsFile(FioFsMountPoint, FioFsFile)
}

// WithFsFileExt add a fio target on a filesystem
// the target file is created explicitly by e2e-fio before fio is launched
func (e *E2eFioArgsBuilder) WithFsFileExt(fsMount string, fileRelPath string, fsAllocStrategy FioFsAllocation, units uint) *E2eFioArgsBuilder {
	tgt := fioTarget{
		target:          fsMount + "/" + fileRelPath,
		device:          fsMount,
		fsFile:          fileRelPath,
		fsAllocStrategy: fsAllocStrategy,
		fsAllocUnit:     units,
	}
	e.targets = append(e.targets, tgt)
	return e
}

// WithDefaultFileExt add the "standard" fio target on a filesystem
// the target file is created explicitly by e2e-fio before fio is launched
func (e *E2eFioArgsBuilder) WithDefaultFileExt(fsAllocStrategy FioFsAllocation, units uint) *E2eFioArgsBuilder {
	return e.WithFsFileExt(FioFsMountPoint, FioFsFile, fsAllocStrategy, units)
}

// WithTargets add an existing file as fio target
// used when a fio instance is launched to verify
// writes of a previous fio instance
func (e *E2eFioArgsBuilder) WithTargets(targets []string) *E2eFioArgsBuilder {
	for _, tgt := range targets {
		e.targets = append(e.targets, fioTarget{
			target: tgt,
		})
	}
	return e
}

// WithRawBlock add a fio target device
func (e *E2eFioArgsBuilder) WithRawBlock(devicePath string) *E2eFioArgsBuilder {
	e.targets = append(e.targets, fioTarget{
		target: devicePath,
		device: devicePath,
	})
	return e
}

// WithDefaultRawBlock add the "standard" fio target device
func (e *E2eFioArgsBuilder) WithDefaultRawBlock() *E2eFioArgsBuilder {
	return e.WithRawBlock(FioBlockFilename)
}

// WithRuntime set duration for time based run
// duration of 0 => no runtime limit
func (e *E2eFioArgsBuilder) WithRuntime(duration int) *E2eFioArgsBuilder {
	e.duration = duration
	return e
}

// WithLivenessParameters add liveness checks for fio targets
func (e *E2eFioArgsBuilder) WithLivenessParameters(readInterval uint, timeout uint) *E2eFioArgsBuilder {
	e.liveness = true
	e.livenessReadInterval = readInterval
	e.livenessTimeout = timeout
	return e
}

// WithLiveness add liveness checks for fio targets
func (e *E2eFioArgsBuilder) WithLiveness() *E2eFioArgsBuilder {
	e.liveness = true
	return e
}

// WithExitValue add exit value override
func (e *E2eFioArgsBuilder) WithExitValue(exitValue int) *E2eFioArgsBuilder {
	e.exitValue = exitValue
	return e
}

// WithZeroFill add commands to zerofill the targets
func (e *E2eFioArgsBuilder) WithZeroFill(val bool) *E2eFioArgsBuilder {
	e.zerofill = val
	return e
}

// WithBlockSize set fio block size
func (e *E2eFioArgsBuilder) WithBlockSize(val uint) *E2eFioArgsBuilder {
	e.blockSize = val
	return e
}

// WithLoops number of fio iterations over disk
func (e *E2eFioArgsBuilder) WithLoops(val uint) *E2eFioArgsBuilder {
	e.loops = val
	return e
}

// WithDirectIO set args for direct IO
func (e *E2eFioArgsBuilder) WithDirectIO(direct bool) {
	e.indirectIO = !direct
}

func (e *E2eFioArgsBuilder) Build() ([]string, error) {
	var cmdLine []string
	var fioArgs []string

	if e.err != nil {
		return cmdLine, e.err
	}

	bs := fmt.Sprintf("--bs=%d", e.blockSize)
	switch e.argsSet {
	case NoFioArgs:
		break
	case DefaultFioArgs:
		fioArgs = append(fioArgs, bs)
		fioArgs = append(fioArgs, fioDefaultParams...)
	case RandWriteFioArgs:
		fioArgs = append(fioArgs, bs)
		fioArgs = append(fioArgs, fioRandWriteParams...)
	case RandReadFioArgs:
		fioArgs = append(fioArgs, bs)
		fioArgs = append(fioArgs, fioRandReadParams...)
	case RandReadWriteFioArgs:
		fioArgs = append(fioArgs, bs)
		fioArgs = append(fioArgs, fioRandReadWriteParams...)
	case CustomFioArgs:
		fioArgs = append(fioArgs, bs)
		fioArgs = append(fioArgs, e.customArgs...)
	case PerfSeqReadFioArgs:
		fioArgs = append(fioArgs, fioPerfSeqReadParams...)
	case PerfSeqWriteFioArgs:
		fioArgs = append(fioArgs, fioPerfSeqWriteParams...)
	case PerfSeqMixedFioArgs:
		fioArgs = append(fioArgs, fioPerfSeqMixedParams...)
	case PerfRandReadFioArgs:
		fioArgs = append(fioArgs, fioPerfRandReadParams...)
	case PerfRandWriteFioArgs:
		fioArgs = append(fioArgs, fioPerfRandWriteParams...)
	case PerfRandMixedFioArgs:
		fioArgs = append(fioArgs, fioPerfRandMixedParams...)
	}

	{
		// configuration overrides arguments passed in
		// add direct io argument derived from the configuration
		// and discard any --direct=* in the argument list
		tmp := []string{"--direct=1"}
		if e.indirectIO {
			tmp = []string{"--direct=0"}
		}
		for _, arg := range fioArgs {
			if strings.HasPrefix(arg, "--direct=") {
			} else {
				tmp = append(tmp, arg)
			}
		}
		fioArgs = tmp
	}

	// 0. exit value
	if e.exitValue != 0 {
		cmdLine = append(cmdLine, []string{"exitv", fmt.Sprintf("%d", e.exitValue)}...)
	}

	// 1. make files
	for _, tgt := range e.targets {
		// note target may be a path to a file, which is not created by the wrapper
		if tgt.fsFile != "" {
			switch tgt.fsAllocStrategy {
			case FioFsAllocDefault:
				cmdLine = append(cmdLine, []string{
					"makefile", tgt.device, tgt.fsFile, "availblockslessby", "20", ";",
				}...)
			case FioFsAllocLessByBlocks:
				cmdLine = append(cmdLine, []string{
					"makefile", tgt.device, tgt.fsFile, "availblockslessby", fmt.Sprintf("%d", tgt.fsAllocUnit), ";",
				}...)
			case FioFsAllocPercentage:
				cmdLine = append(cmdLine, []string{
					"makefile", tgt.device, tgt.fsFile, "availblockspercent", fmt.Sprintf("%d", tgt.fsAllocUnit), ";",
				}...)
			}
		}
	}

	// 2. zerofill
	if e.zerofill {
		for _, tgt := range e.targets {
			if tgt.fsFile != "" {
				cmdLine = append(cmdLine, []string{
					"zerofill", tgt.device + "/" + tgt.fsFile, ";",
				}...)
			} else {
				if tgt.device != "" {
					cmdLine = append(cmdLine, []string{
						"zerofill", tgt.device, ";",
					}...)
				}
			}
		}
	}

	// 3. liveness
	if e.liveness {
		for _, tgt := range e.targets {
			if tgt.device != "" && tgt.target != "" {
				cmdLine = append(cmdLine, []string{
					"liveness",
					tgt.target,
					fmt.Sprintf("%v", e.livenessReadInterval),
					fmt.Sprintf("%v", e.livenessTimeout),
					";",
				}...)
			}
		}
	}

	// 4. target size
	for _, tgt := range e.targets {
		if tgt.fsFile != "" {
			cmdLine = append(cmdLine, []string{
				"filesize",
				fmt.Sprintf("%s/%s", tgt.device, tgt.fsFile),
				";",
			}...)
		} else {
			cmdLine = append(cmdLine, []string{
				"filesize",
				tgt.device,
				";",
			}...)
		}
	}

	// 5. fio
	if len(fioArgs) != 0 {
		cmdLine = append(cmdLine, []string{"---", "fio", "--verify_dump=1"}...)
		if e.duration != 0 {
			cmdLine = append(cmdLine, "--loops=99999")
		} else if e.loops != 0 {
			cmdLine = append(cmdLine, fmt.Sprintf("--loops=%d", e.loops))
		}
		cmdLine = append(cmdLine, fioArgs...)
		cmdLine = append(cmdLine, e.additionalArgs...)

		for ix, tgt := range e.targets {
			cmdLine = append(cmdLine, fmt.Sprintf("--name=benchtest%d", ix))
			if tgt.targetSize != 0 {
				cmdLine = append(cmdLine, fmt.Sprintf("--size=%v", tgt.targetSize))
			}
			cmdLine = append(cmdLine, fmt.Sprintf("--filename=%s", tgt.target))
		}
		cmdLine = append(cmdLine, ";")
	}

	// 6. duration
	if e.duration != 0 {
		cmdLine = append(cmdLine, []string{
			"sigterm",
			fmt.Sprintf("%d", e.duration),
			";",
		}...)
	}

	return cmdLine, nil
}

func (e *E2eFioArgsBuilder) GetTargets() []string {
	targets := []string{}
	for _, tgt := range e.targets {
		targets = append(targets, tgt.target)
	}
	return targets
}
