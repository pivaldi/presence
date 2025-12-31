package tests

import (
	"testing"

	"github.com/pivaldi/nullable"
	"github.com/stretchr/testify/assert"
)

func TestMarshalUnsetBehavior(t *testing.T) {
	t.Run("UnsetSkip is default", func(t *testing.T) {
		assert.Equal(t, nullable.MarshalUnsetBehavior(0), nullable.UnsetSkip)
	})

	t.Run("UnsetNull is alternative", func(t *testing.T) {
		assert.Equal(t, nullable.MarshalUnsetBehavior(1), nullable.UnsetNull)
	})
}

func TestScanNullBehaviorConstants(t *testing.T) {
	t.Run("ScanNullAsNull is default", func(t *testing.T) {
		assert.Equal(t, nullable.ScanNullBehavior(0), nullable.ScanNullAsNull)
	})

	t.Run("ScanNullAsUnset is alternative", func(t *testing.T) {
		assert.Equal(t, nullable.ScanNullBehavior(1), nullable.ScanNullAsUnset)
	})
}

func TestDefaultConfiguration(t *testing.T) {
	t.Run("default marshal unset is skip", func(t *testing.T) {
		assert.Equal(t, nullable.UnsetSkip, nullable.GetDefaultMarshalUnset())
	})

	t.Run("default scan null is null", func(t *testing.T) {
		assert.Equal(t, nullable.ScanNullAsNull, nullable.GetDefaultScanNull())
	})
}

func TestPerValueConfiguration(t *testing.T) {
	t.Run("SetMarshalUnset configures per-value behavior", func(t *testing.T) {
		n := nullable.Of[string]{}
		n.SetMarshalUnset(nullable.UnsetNull)
		assert.Equal(t, nullable.UnsetNull, n.GetMarshalUnset())
	})

	t.Run("SetScanNull configures per-value behavior", func(t *testing.T) {
		n := nullable.Of[string]{}
		n.SetScanNull(nullable.ScanNullAsUnset)
		assert.Equal(t, nullable.ScanNullAsUnset, n.GetScanNull())
	})

	t.Run("default uses package default for marshal", func(t *testing.T) {
		n := nullable.Of[string]{}
		assert.Equal(t, nullable.GetDefaultMarshalUnset(), n.GetMarshalUnset())
	})

	t.Run("default uses package default for scan", func(t *testing.T) {
		n := nullable.Of[string]{}
		assert.Equal(t, nullable.GetDefaultScanNull(), n.GetScanNull())
	})
}
