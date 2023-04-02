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

type MockAppendTableBackendOps struct {
	ErrorRval    error
	DataRowsRval []mutator.MappedFieldValues
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
	return b.ErrorRval
}

func (b *MockAppendTableBackendOps) DeleteAll(batchSize int) error {
	return b.ErrorRval
}

func TestSettings(t *testing.T) {
	table := logtable.NewLogTable()
	table.Init()
	actual := table.GetSettings()

	testutils.AssertEquals(t, "Log", actual.Name)
	testutils.AssertEquals(t, &logtable.LogDataRowSettings, actual.DataRowSettings)

	actualMessage, ok := actual.DataRowSettings.EmptyValues["Message"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualMessage)

	actualSource, ok := actual.DataRowSettings.EmptyValues["Source"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualSource)

	actualLevel, ok := actual.DataRowSettings.EmptyValues["Level"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualLevel)

	actualCreatedTime, ok := actual.DataRowSettings.EmptyValues["CreatedTime"].(fields.Time)
	testutils.AssertTrue(t, ok)
	testutils.AssertTrue(t, fields.Time.IsZero(actualCreatedTime))
}

func TestScan(t *testing.T) {
	type scanExpectedVals struct {
		errors     []error
		scanFields []logtable.LogDataRow
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*MockAppendTableBackendOps, *scanExpectedVals]{
		Cases: []testutils.TestCase[*MockAppendTableBackendOps, *scanExpectedVals]{
			{
				Name: "properly formats good values",
				Input: &MockAppendTableBackendOps{
					ErrorRval: nil,
					DataRowsRval: []mutator.MappedFieldValues{
						{"Message": "test1"},
						{"Message": "test2"},
					},
				},
				Expected: &scanExpectedVals{
					scanFields: []logtable.LogDataRow{
						{Message: "test1"},
						{Message: "test2"},
					},
				},
			},
			{
				Name: "returns error for mismatched types",
				Input: &MockAppendTableBackendOps{
					ErrorRval: nil,
					DataRowsRval: []mutator.MappedFieldValues{
						{"Message": 0},
					},
				},
				Expected: &scanExpectedVals{
					errors: []error{mutator.SetFieldTypeError},
				},
			},
			{
				Name: "handles error and result input",
				Input: &MockAppendTableBackendOps{
					ErrorRval: mockError,
					DataRowsRval: []mutator.MappedFieldValues{
						{"Message": "test"},
					},
				},
				Expected: &scanExpectedVals{
					scanFields: []logtable.LogDataRow{
						{Message: "test"},
					},
					errors: []error{mockError},
				},
			},
		},
		Func: func(mockBackend *MockAppendTableBackendOps, expected *scanExpectedVals) {
			table := logtable.NewLogTable()
			table.SetBackend(mockBackend)

			actualScanFieldsChan, actualErrorChan := table.Scan(10)

			for _, expectedScanFields := range expected.scanFields {
				actualScanFields, more := <-actualScanFieldsChan
				if !more {
					t.Errorf("actualScanFieldsChan ended prematurely")
				}

				testutils.AssertEquals(t, expectedScanFields.Message, actualScanFields.Message)
			}

			_, more := <-actualScanFieldsChan
			if more {
				t.Errorf("actualScanFieldsChan longer than expected")
			}

			for _, expectedError := range expected.errors {
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

func TestAppend(t *testing.T) {

}

func TestDeleteAll(t *testing.T) {

}

func TestTransferTo(t *testing.T) {

}
