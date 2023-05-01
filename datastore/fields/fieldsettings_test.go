package fields_test

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/testutils"
)

func TestNewFieldSettings(t *testing.T) {
	settings := fields.NewFieldSettings(
		fields.WithAutoGenerate("test"),
		fields.WithNumBytes("test", 31),
	)

	setting := settings["test"]
	testutils.AssertTrue(t, setting.AutoGenerate)
	testutils.AssertEquals(t, 31, setting.NumBytes)
}
