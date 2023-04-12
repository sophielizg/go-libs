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

func TestMergeSettings(t *testing.T) {
	settings1 := fields.NewFieldSettings(
		fields.WithAutoGenerate("test1"),
	)

	settings2 := fields.NewFieldSettings(
		fields.WithNumBytes("test2", 31),
	)

	merged := fields.MergeSettings(settings1, settings2)

	testutils.AssertTrue(t, merged["test1"].AutoGenerate)
	testutils.AssertEquals(t, 31, merged["test2"].NumBytes)
}
