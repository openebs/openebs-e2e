package main

import (
	"fmt"
	"os"
	"regexp"
)

// multipathConfPath and multipathConfBackupPath are under /host because
// this agent's container only sees the host root there (see the host-root
// volume mount) -- the same convention SYSRQ_TRIGGER_FILE uses.
const (
	multipathConfPath       = "/host/etc/multipath.conf"
	multipathConfBackupPath = "/host/etc/multipath.conf.e2e-agent-bak"
)

var findMultipathsRe = regexp.MustCompile(`(?m)^(\s*find_multipaths\s+)\S+(\s*)$`)
var defaultsBlockRe = regexp.MustCompile(`(?m)^(\s*defaults\s*\{)`)

// SetFindMultipaths sets multipath.conf's find_multipaths default to value
// ("yes" or "no"), backing up the file first (once -- a second call before
// RestoreMultipathConf reuses the original backup rather than overwriting
// it with an already-modified file). The caller still has to restart
// multipathd (e.g. via RestartService("multipathd")) for this to take
// effect, and must call RestoreMultipathConf when done: this agent has no
// way to do that on its own if the test process never returns.
func SetFindMultipaths(value string) error {
	if value != "yes" && value != "no" {
		return fmt.Errorf("find_multipaths must be \"yes\" or \"no\", got %q", value)
	}

	if _, err := os.Stat(multipathConfBackupPath); os.IsNotExist(err) {
		original, err := os.ReadFile(multipathConfPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", multipathConfPath, err)
		}
		if err := os.WriteFile(multipathConfBackupPath, original, 0o644); err != nil {
			return fmt.Errorf("backing up %s: %w", multipathConfPath, err)
		}
	} else if err != nil {
		return fmt.Errorf("checking for existing backup %s: %w", multipathConfBackupPath, err)
	}

	content, err := os.ReadFile(multipathConfBackupPath)
	if err != nil {
		return fmt.Errorf("reading backed-up %s: %w", multipathConfBackupPath, err)
	}

	updated := setFindMultipathsIn(string(content), value)

	if err := os.WriteFile(multipathConfPath, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", multipathConfPath, err)
	}
	return nil
}

// setFindMultipathsIn returns content with find_multipaths set to value: an
// existing find_multipaths line (wherever it is) is rewritten in place;
// otherwise it is inserted just inside the first "defaults {" block; failing
// that (no defaults block at all), a new one is appended.
func setFindMultipathsIn(content, value string) string {
	if findMultipathsRe.MatchString(content) {
		return findMultipathsRe.ReplaceAllString(content, fmt.Sprintf("${1}%s${2}", value))
	}
	line := fmt.Sprintf("find_multipaths %s", value)
	if loc := defaultsBlockRe.FindStringIndex(content); loc != nil {
		insertAt := loc[1]
		return content[:insertAt] + "\n\t" + line + content[insertAt:]
	}
	return content + fmt.Sprintf("\ndefaults {\n\t%s\n}\n", line)
}

// RestoreMultipathConf puts back the multipath.conf SetFindMultipaths
// backed up, if any -- a no-op if SetFindMultipaths was never called (or
// was already restored), so repeated cleanup calls are safe. The caller
// still has to restart multipathd for the restored file to take effect.
func RestoreMultipathConf() error {
	original, err := os.ReadFile(multipathConfBackupPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading backup %s: %w", multipathConfBackupPath, err)
	}
	if err := os.WriteFile(multipathConfPath, original, 0o644); err != nil {
		return fmt.Errorf("restoring %s: %w", multipathConfPath, err)
	}
	return os.Remove(multipathConfBackupPath)
}
