//go:build openharmony

// filecheck_openharmony.go provides a no-op implementation for OpenHarmony.
// OpenHarmony's kernel makes Unix UID checks and symlink checks
// unreliable (process UID != file owner UID, filesystem overrides chmod).
// Security is enforced by the OS rather than file permission checks.

package filecheck

const noopValidateBeforeRead = true

// ValidateBeforeRead is a no-op on OpenHarmony.
func ValidateBeforeRead(_ string) error {
	return nil
}
