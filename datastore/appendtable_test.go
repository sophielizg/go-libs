package datastore_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/examples/logtable"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS

func AssertLogDataRowEquals(t *testing.T, expected, actual *logtable.LogDataRow) {
	t.Helper()
	testutils.AssertEquals(t, expected.Message, actual.Message)
	testutils.AssertEquals(t, expected.Source, actual.Source)
	testutils.AssertEquals(t, expected.Level, actual.Level)
	testutils.AssertEquals(t, expected.CreatedTime, actual.CreatedTime)
}

func AssertLogDataRowFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedMessageVal, expectedOk := expected[logtable.MessageKey].(fields.String)
	actualMessageVal, ok := actual[logtable.MessageKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedMessageVal, actualMessageVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualMessageVal)
	}

	expectedSourceVal, expectedOk := expected[logtable.SourceKey].(fields.String)
	actualSourceVal, ok := actual[logtable.SourceKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedSourceVal, actualSourceVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualSourceVal)
	}

	expectedLevelVal, expectedOk := expected[logtable.LevelKey].(fields.String)
	actualLevelVal, ok := actual[logtable.LevelKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedLevelVal, actualLevelVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualLevelVal)
	}

	expectedCreatedTimeVal, expectedOk := expected[logtable.CreatedTimeKey].(fields.Time)
	actualCreatedTimeVal, ok := actual[logtable.CreatedTimeKey].(fields.Time)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedCreatedTimeVal, actualCreatedTimeVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, fields.Time{}, actualCreatedTimeVal)
	}
}

// MOCKS

type MockAppendTableBackendOps struct {
	ErrorRval     error
	DataRowsRval  []mutator.MappedFieldValues
	DataRowsInput []mutator.MappedFieldValues
}

func (b *MockAppendTableBackendOps) Scan(batchSize int) (chan *datastore.ScanFields, chan error) {
	dataChan := make(chan *datastore.ScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errorChan)

		for _, dataRow := range b.DataRowsRval {
			dataChan <- &datastore.ScanFields{
				DataRow: dataRow,
			}
		}

		if b.ErrorRval != nil {
			errorChan <- b.ErrorRval
		}
	}()

	return dataChan, errorChan
}

func (b *MockAppendTableBackendOps) AddMultiple(data []mutator.MappedFieldValues) error {
	b.DataRowsInput = data
	return b.ErrorRval
}

// TESTS

func TestAppendTableSettings(t *testing.T) {
	table := logtable.NewLogTable()
	table.Init()
	actual := table.GetSettings()

	testutils.AssertEquals(t, "Log", actual.Name)
	testutils.AssertEquals(t, &logtable.LogDataRowSettings, actual.DataRowSettings)

	// Check that defaults are generated
	AssertLogDataRowFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.DataRowSettings.EmptyValues)
}

func TestAppendTableScan(t *testing.T) {
	type scanInputVals struct {
		Error    error
		DataRows []mutator.MappedFieldValues
	}

	type scanExpectedVals struct {
		Errors     []error
		ScanFields []*logtable.LogDataRow
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*scanInputVals, *scanExpectedVals]{
		Cases: []testutils.TestCase[*scanInputVals, *scanExpectedVals]{
			{
				Name: "properly formats good values",
				Input: &scanInputVals{
					Error: nil,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: "test1"},
						{logtable.MessageKey: "test2"},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []*logtable.LogDataRow{
						{Message: "test1"},
						{Message: "test2"},
					},
				},
			},
			{
				Name: "returns error for mismatched types",
				Input: &scanInputVals{
					Error: nil,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: 0},
					},
				},
				Expected: &scanExpectedVals{
					Errors: []error{mutator.SetFieldTypeError},
				},
			},
			{
				Name: "handles error and result input",
				Input: &scanInputVals{
					Error: mockError,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: "test"},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []*logtable.LogDataRow{
						{Message: "test"},
					},
					Errors: []error{mockError},
				},
			},
		},
		Func: func(input *scanInputVals, expected *scanExpectedVals) {
			mockBackend := &MockAppendTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRows,
			}
			table := logtable.NewLogTable()
			table.SetBackend(mockBackend)

			actualScanFieldsChan, actualErrorChan := table.Scan(10)

			for _, expectedScanFields := range expected.ScanFields {
				actualScanFields, more := <-actualScanFieldsChan
				if !more {
					t.Errorf("actualScanFieldsChan ended prematurely")
				}

				AssertLogDataRowEquals(t, expectedScanFields, actualScanFields)
			}

			_, more := <-actualScanFieldsChan
			if more {
				t.Errorf("actualScanFieldsChan longer than expected")
			}

			for _, expectedError := range expected.Errors {
				actualError, more := <-actualErrorChan
				if !more {
					t.Errorf("actualErrorChan ended prematurely")
				}

				testutils.AssertErrorEquals(t, expectedError, actualError)
			}

			_, more = <-actualErrorChan
			if more {
				t.Errorf("actualErrorChan longer than expected")
			}
		},
	}

	tests.Run(t)
}

func TestAppendTableAdd(t *testing.T) {
	type addInputVals struct {
		Error    error
		DataRows []*logtable.LogDataRow
	}

	type addExpectedVals struct {
		Error    error
		DataRows []mutator.MappedFieldValues
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*addInputVals, *addExpectedVals]{
		Cases: []testutils.TestCase[*addInputVals, *addExpectedVals]{
			{
				Name: "successfully adds",
				Input: &addInputVals{
					Error: nil,
					DataRows: []*logtable.LogDataRow{
						{Message: "test1"},
						{Message: "test2"},
					},
				},
				Expected: &addExpectedVals{
					Error: nil,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: "test1"},
						{logtable.MessageKey: "test2"},
					},
				},
			},
			{
				Name: "returns error",
				Input: &addInputVals{
					Error:    mockError,
					DataRows: []*logtable.LogDataRow{},
				},
				Expected: &addExpectedVals{
					Error:    mockError,
					DataRows: []mutator.MappedFieldValues{},
				},
			},
		},
		Func: func(input *addInputVals, expected *addExpectedVals) {
			mockBackend := &MockAppendTableBackendOps{
				ErrorRval: input.Error,
			}
			table := logtable.NewLogTable()
			table.SetBackend(mockBackend)

			err := table.Add(input.DataRows...)
			testutils.AssertErrorEquals(t, expected.Error, err)

			testutils.AssertEquals(t, len(expected.DataRows), len(mockBackend.DataRowsInput))
			for i := range expected.DataRows {
				AssertLogDataRowFieldsEqualOrDefault(t, expected.DataRows[i], mockBackend.DataRowsInput[i])
			}
		},
	}

	tests.Run(t)
}

func TestAppendTableTransferTo(t *testing.T) {
	type transferInputVals struct {
		Error    error
		DataRows []mutator.MappedFieldValues
	}

	type transferExpectedVals struct {
		Error error
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*transferInputVals, *transferExpectedVals]{
		Cases: []testutils.TestCase[*transferInputVals, *transferExpectedVals]{
			{
				Name: "transfers good values",
				Input: &transferInputVals{
					Error: nil,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: "test1"},
						{logtable.MessageKey: "test2"},
					},
				},
				Expected: &transferExpectedVals{
					Error: nil,
				},
			},
			{
				Name: "returns error for mismatched types",
				Input: &transferInputVals{
					Error: nil,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: 1},
					},
				},
				Expected: &transferExpectedVals{
					Error: mutator.SetFieldTypeError,
				},
			},
			{
				Name: "handles error from backend",
				Input: &transferInputVals{
					Error: mockError,
					DataRows: []mutator.MappedFieldValues{
						{logtable.MessageKey: "test"},
					},
				},
				Expected: &transferExpectedVals{
					Error: mockError,
				},
			},
		},
		Func: func(input *transferInputVals, expected *transferExpectedVals) {
			mockBackendSrc := &MockAppendTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRows,
			}
			srcTable := logtable.NewLogTable()
			srcTable.SetBackend(mockBackendSrc)

			mockBackendDest := &MockAppendTableBackendOps{}
			destTable := logtable.NewLogTable()
			destTable.SetBackend(mockBackendDest)

			err := srcTable.TransferTo(destTable, 10)
			testutils.AssertErrorEquals(t, expected.Error, err)

			if expected.Error != nil {
				return
			}

			testutils.AssertEquals(t, len(input.DataRows), len(mockBackendDest.DataRowsInput))
			for i := range input.DataRows {
				AssertLogDataRowFieldsEqualOrDefault(t, input.DataRows[i], mockBackendDest.DataRowsInput[i])
			}
		},
	}

	tests.Run(t)
}
