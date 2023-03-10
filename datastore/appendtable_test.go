package datastore_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/testutils"
)

type mockAppendTableBackend struct {
	errorRval          error
	errorChanRval      chan error
	scanFieldsChanRval chan *datastore.AppendTableScanFields
}

func (b *mockAppendTableBackend) ValidateSchema(schema *datastore.AppendTableSchema) error {
	return b.errorRval
}

func (b *mockAppendTableBackend) CreateOrUpdateSchema(schema *datastore.AppendTableSchema) error {
	return b.errorRval
}

func (b *mockAppendTableBackend) Scan(schema *datastore.AppendTableSchema, batchSize int) (chan *datastore.AppendTableScanFields, chan error) {
	return b.scanFieldsChanRval, b.errorChanRval
}

func (b *mockAppendTableBackend) AppendMultiple(schema *datastore.AppendTableSchema, data []datastore.DataRow) error {
	return b.errorRval
}

func (b *mockAppendTableBackend) DeleteAll(schema *datastore.AppendTableSchema, batchSize int) error {
	return b.errorRval
}

type testDataRow struct {
	val string
}

func (d *testDataRow) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"val": d.val,
	}
}

type testDataRowFactory struct{}

func (f *testDataRowFactory) CreateDefault() *testDataRow {
	return nil
}

func (f *testDataRowFactory) CreateFromFields(fields datastore.DataRowFields) (*testDataRow, error) {
	return &testDataRow{
		val: fields["val"].(string),
	}, nil
}

func (f *testDataRowFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"val": &datastore.StringField{NumChars: 64},
	}
}

func testTable(backend datastore.AppendTableBackend) datastore.AppendTable[*testDataRow] {
	return datastore.AppendTable[*testDataRow]{
		Name:           "Test",
		DataRowFactory: &testDataRowFactory{},
		Backend:        backend,
	}
}

func TestValidateSchema(t *testing.T) {
	tests := testutils.Tests[error, error]{
		Cases: []testutils.TestCase[error, error]{
			{
				Name:     "returns ok if no error from backend",
				Input:    nil,
				Expected: nil,
			},
			{
				Name:     "returns error if one comes from backend",
				Input:    errors.New("error"),
				Expected: errors.New("error"),
			},
		},
		Func: func(errorRval error, expected error) {
			table := testTable(&mockAppendTableBackend{
				errorRval: errorRval,
			})
			testutils.ErrorEquals(t, expected, table.ValidateSchema())
		},
	}

	tests.Run(t)
}

func TestCreateOrUpdateSchema(t *testing.T) {
	tests := testutils.Tests[error, error]{
		Cases: []testutils.TestCase[error, error]{
			{
				Name:     "returns ok if no error from backend",
				Input:    nil,
				Expected: nil,
			},
			{
				Name:     "returns error if one comes from backend",
				Input:    errors.New("error"),
				Expected: errors.New("error"),
			},
		},
		Func: func(errorRval error, expected error) {
			table := testTable(&mockAppendTableBackend{
				errorRval: errorRval,
			})
			testutils.ErrorEquals(t, expected, table.CreateOrUpdateSchema())
		},
	}

	tests.Run(t)
}

func TestScan(t *testing.T) {
	type scanInVals struct {
		errors     []error
		scanFields []datastore.AppendTableScanFields
	}

	type scanExpectedVals struct {
		errors     []error
		scanFields []datastore.DataRowScan[*testDataRow]
	}

	tests := testutils.Tests[*scanInVals, *scanExpectedVals]{
		Cases: []testutils.TestCase[*scanInVals, *scanExpectedVals]{
			{
				Name: "properly formats good values",
				Input: &scanInVals{
					scanFields: []datastore.AppendTableScanFields{
						{
							DataRow: datastore.DataRowFields{
								"val": "test1",
							},
						},
						{
							DataRow: datastore.DataRowFields{
								"val": "test2",
							},
						},
					},
				},
				Expected: &scanExpectedVals{
					scanFields: []datastore.DataRowScan[*testDataRow]{
						{
							DataRow: &testDataRow{
								val: "test1",
							},
						},
						{
							DataRow: &testDataRow{
								val: "test2",
							},
						},
					},
				},
			},
			{
				Name: "returns error for mismatched types",
				Input: &scanInVals{
					scanFields: []datastore.AppendTableScanFields{
						{
							DataRow: datastore.DataRowFields{
								"val": 0,
							},
						},
					},
				},
				Expected: &scanExpectedVals{
					errors: []error{
						errors.New("error"),
					},
				},
			},
			{
				Name: "handles error and result input",
				Input: &scanInVals{
					scanFields: []datastore.AppendTableScanFields{
						{
							DataRow: datastore.DataRowFields{
								"val": "test",
							},
						},
					},
					errors: []error{
						errors.New("error"),
					},
				},
				Expected: &scanExpectedVals{
					scanFields: []datastore.DataRowScan[*testDataRow]{
						{
							DataRow: &testDataRow{
								val: "test",
							},
						},
					},
					errors: []error{
						errors.New("error"),
					},
				},
			},
		},
		Func: func(inVals *scanInVals, expected *scanExpectedVals) {
			batchSize := 10

			scanFieldsChan := make(chan *datastore.AppendTableScanFields, batchSize)
			errorChan := make(chan error, batchSize)

			for _, scanFields := range inVals.scanFields {
				closure := scanFields
				scanFieldsChan <- &closure
			}

			for _, err := range inVals.errors {
				errorChan <- err
			}

			close(scanFieldsChan)
			close(errorChan)

			table := testTable(&mockAppendTableBackend{
				scanFieldsChanRval: scanFieldsChan,
				errorChanRval:      errorChan,
			})

			actualScanFieldsChan, actualErrorChan := table.Scan(batchSize)

			for _, expectedScanFields := range expected.scanFields {
				actualScanFields, more := <-actualScanFieldsChan
				if !more {
					t.Errorf("actualScanFieldsChan ended prematurely")
				}

				testutils.Equals(t, expectedScanFields.DataRow.val, actualScanFields.DataRow.val)
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

				testutils.ErrorEquals(t, expectedError, actualError)
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
