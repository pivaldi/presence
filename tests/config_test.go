package tests

import (
	"testing"

	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
)

func TestMarshalUnsetBehavior(t *testing.T) {
	t.Run("UnsetSkip is default", func(t *testing.T) {
		assert.Equal(t, presence.MarshalUnsetBehavior(0), presence.UnsetSkip)
	})

	t.Run("UnsetNull is alternative", func(t *testing.T) {
		assert.Equal(t, presence.MarshalUnsetBehavior(1), presence.UnsetNull)
	})
}

func TestScanNullBehaviorConstants(t *testing.T) {
	t.Run("ScanNullAsNull is default", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullBehavior(0), presence.ScanNullAsNull)
	})

	t.Run("ScanNullAsUnset is alternative", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullBehavior(1), presence.ScanNullAsUnset)
	})
}

func TestDefaultConfiguration(t *testing.T) {
	t.Run("default marshal unset is skip", func(t *testing.T) {
		assert.Equal(t, presence.UnsetSkip, presence.GetDefaultMarshalUnset())
	})

	t.Run("default scan null is null", func(t *testing.T) {
		assert.Equal(t, presence.ScanNullAsNull, presence.GetDefaultScanNull())
	})
}

func TestPerValueConfiguration(t *testing.T) {
	t.Run("SetMarshalUnset configures per-value behavior", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetMarshalUnset(presence.UnsetNull)
		assert.Equal(t, presence.UnsetNull, n.GetMarshalUnset())
	})

	t.Run("SetScanNull configures per-value behavior", func(t *testing.T) {
		n := presence.Of[string]{}
		n.SetScanNull(presence.ScanNullAsUnset)
		assert.Equal(t, presence.ScanNullAsUnset, n.GetScanNull())
	})

	t.Run("default uses package default for marshal", func(t *testing.T) {
		n := presence.Of[string]{}
		assert.Equal(t, presence.GetDefaultMarshalUnset(), n.GetMarshalUnset())
	})

	t.Run("default uses package default for scan", func(t *testing.T) {
		n := presence.Of[string]{}
		assert.Equal(t, presence.GetDefaultScanNull(), n.GetScanNull())
	})
}
