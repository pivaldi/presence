package presence

import "sync"

// MarshalUnsetBehavior controls how unset values are marshaled to JSON.
type MarshalUnsetBehavior int

const (
	// UnsetSkip omits unset fields from JSON output (requires omitempty tag).
	UnsetSkip MarshalUnsetBehavior = iota
	// UnsetNull marshals unset fields as null.
	UnsetNull
)

// ScanNullBehavior controls how SQL NULL values are scanned.
type ScanNullBehavior int

const (
	// ScanNullAsNull interprets SQL NULL as explicit null (isSet=true, val=nil).
	ScanNullAsNull ScanNullBehavior = iota
	// ScanNullAsUnset interprets SQL NULL as unset (isSet=false, val=nil).
	ScanNullAsUnset
)

var (
	defaultMarshalUnset MarshalUnsetBehavior = UnsetSkip
	defaultScanNull     ScanNullBehavior     = ScanNullAsNull
	configMu            sync.RWMutex
)

// SetDefaultMarshalUnset sets the package-level default for marshal unset behavior.
func SetDefaultMarshalUnset(b MarshalUnsetBehavior) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultMarshalUnset = b
}

// GetDefaultMarshalUnset returns the package-level default for marshal unset behavior.
func GetDefaultMarshalUnset() MarshalUnsetBehavior {
	configMu.RLock()
	defer configMu.RUnlock()

	return defaultMarshalUnset
}

// SetDefaultScanNull sets the package-level default for scan null behavior.
func SetDefaultScanNull(b ScanNullBehavior) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultScanNull = b
}

// GetDefaultScanNull returns the package-level default for scan null behavior.
func GetDefaultScanNull() ScanNullBehavior {
	configMu.RLock()
	defer configMu.RUnlock()

	return defaultScanNull
}
