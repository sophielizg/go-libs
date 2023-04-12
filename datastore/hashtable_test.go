package datastore_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/examples/product"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

type MockHashTableBackendOps struct {
	ErrorRval     error
	DataRowsRval  []mutator.MappedFieldValues
	HashKeyRval   []mutator.MappedFieldValues
	DataRowsInput []mutator.MappedFieldValues
	HashKeysInput []mutator.MappedFieldValues
}

func (b *MockHashTableBackendOps) Scan(batchSize int) (chan *datastore.ScanFields, chan error) {
	dataChan := make(chan *datastore.ScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errorChan)

		for i := range b.DataRowsRval {
			dataChan <- &datastore.ScanFields{
				DataRow: b.DataRowsRval[i],
				HashKey: b.HashKeyRval[i],
			}
		}

		if b.ErrorRval != nil {
			errorChan <- b.ErrorRval
		}
	}()

	return dataChan, errorChan
}

func (b *MockHashTableBackendOps) GetMultiple(hashKeys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	return b.DataRowsRval, b.ErrorRval
}

func (b *MockHashTableBackendOps) AddMultiple(hashKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	b.DataRowsInput = data
	b.HashKeysInput = hashKeys
	return b.HashKeyRval, b.ErrorRval
}

func (b *MockHashTableBackendOps) UpdateMultiple(hashKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) error {
	b.DataRowsInput = data
	b.HashKeysInput = hashKeys
	return b.ErrorRval
}

func (b *MockHashTableBackendOps) DeleteMultiple(hashKeys []mutator.MappedFieldValues) error {
	return b.ErrorRval
}

func TestHashTableSettings(t *testing.T) {
	table := product.NewTable()
	table.Init()
	actual := table.GetSettings()

	testutils.AssertEquals(t, "Product", actual.Name)
	testutils.AssertEquals(t, &product.DataRowSettings, actual.DataRowSettings)
	testutils.AssertEquals(t, &product.HashKeySettings, actual.HashKeySettings)

	actualDepartment, ok := actual.DataRowSettings.EmptyValues["Department"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualDepartment)

	actualPrice, ok := actual.DataRowSettings.EmptyValues["Price"].(fields.Float)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, 0.0, actualPrice)

	actualQuantity, ok := actual.DataRowSettings.EmptyValues["Quantity"].(fields.Int)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, 0, actualQuantity)

	actualLastUpdated, ok := actual.DataRowSettings.EmptyValues["LastUpdated"].(fields.Time)
	testutils.AssertTrue(t, ok)
	testutils.AssertTrue(t, fields.Time.IsZero(actualLastUpdated))

	actualBrand, ok := actual.HashKeySettings.EmptyValues["Brand"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualBrand)

	actualName, ok := actual.HashKeySettings.EmptyValues["Name"].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, "", actualName)
}

func TestHashTableScan(t *testing.T) {
	type scanInputVals struct {
		Error    error
		DataRows []mutator.MappedFieldValues
		HashKeys []mutator.MappedFieldValues
	}

	type scanExpectedVals struct {
		Errors     []error
		ScanFields []product.TableScan
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*scanInputVals, *scanExpectedVals]{
		Cases: []testutils.TestCase[*scanInputVals, *scanExpectedVals]{
			{
				Name: "properly formats good values",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{"Name": "test1"},
						{"Name": "test2"},
					},
					DataRows: []mutator.MappedFieldValues{
						{"Price": 9.99},
						{"Price": 10.99},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []product.TableScan{
						{
							DataRow: &product.DataRow{Price: 9.99},
							HashKey: &product.HashKey{Name: "test1"},
						},
						{
							DataRow: &product.DataRow{Price: 10.99},
							HashKey: &product.HashKey{Name: "test2"},
						},
					},
				},
			},
			{
				Name: "returns error for data row mismatched types",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{"Name": "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{"Price": 0},
					},
				},
				Expected: &scanExpectedVals{
					Errors: []error{mutator.SetFieldTypeError},
				},
			},
			{
				Name: "returns error for hash key mismatched types",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{"Name": 0},
					},
					DataRows: []mutator.MappedFieldValues{
						{"Price": 9.99},
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
					HashKeys: []mutator.MappedFieldValues{
						{"Name": "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{"Price": 9.99},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []product.TableScan{
						{
							DataRow: &product.DataRow{Price: 9.99},
							HashKey: &product.HashKey{Name: "test1"},
						},
					},
					Errors: []error{mockError},
				},
			},
		},
		Func: func(input *scanInputVals, expected *scanExpectedVals) {
			mockBackend := &MockHashTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRows,
				HashKeyRval:  input.HashKeys,
			}
			table := product.NewTable()
			table.SetBackend(mockBackend)

			actualScanFieldsChan, actualErrorChan := table.Scan(10)

			for _, expectedScanFields := range expected.ScanFields {
				actualScanFields, more := <-actualScanFieldsChan
				if !more {
					t.Errorf("actualScanFieldsChan ended prematurely")
				}

				testutils.AssertEquals(t, expectedScanFields.DataRow.Price, actualScanFields.DataRow.Price)
				testutils.AssertEquals(t, expectedScanFields.HashKey.Name, actualScanFields.HashKey.Name)
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

func TestHashTableAddMultiple(t *testing.T) {
	type addInputVals struct {
		Error        error
		HashKeys     []*product.HashKey
		DataRows     []*product.DataRow
		HashKeysRval []mutator.MappedFieldValues
	}

	type addExpectedVals struct {
		Error          error
		HashKeysRval   []*product.HashKey
		HashKeysStored []mutator.MappedFieldValues
		DataRowsStored []mutator.MappedFieldValues
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*addInputVals, *addExpectedVals]{
		Cases: []testutils.TestCase[*addInputVals, *addExpectedVals]{
			{
				Name: "successfully adds",
				Input: &addInputVals{
					Error: nil,
					HashKeys: []*product.HashKey{
						{Name: "test1"},
						{Name: "test2"},
					},
					DataRows: []*product.DataRow{
						{Price: 9.99},
						{Price: 10.99},
					},
					HashKeysRval: []mutator.MappedFieldValues{
						{"Name": "test1"},
						{"Name": "test2"},
					},
				},
				Expected: &addExpectedVals{
					Error: nil,
					HashKeysRval: []*product.HashKey{
						{Name: "test1"},
						{Name: "test2"},
					},
					HashKeysStored: []mutator.MappedFieldValues{
						{"Name": "test1"},
						{"Name": "test2"},
					},
					DataRowsStored: []mutator.MappedFieldValues{
						{"Price": 9.99},
						{"Price": 10.99},
					},
				},
			},
			{
				Name: "returns error",
				Input: &addInputVals{
					Error:        mockError,
					DataRows:     []*product.DataRow{},
					HashKeys:     []*product.HashKey{},
					HashKeysRval: []mutator.MappedFieldValues{},
				},
				Expected: &addExpectedVals{
					Error:          mockError,
					HashKeysRval:   nil,
					HashKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
		},
		Func: func(input *addInputVals, expected *addExpectedVals) {
			mockBackend := &MockHashTableBackendOps{
				ErrorRval:   input.Error,
				HashKeyRval: input.HashKeysRval,
			}
			table := product.NewTable()
			table.SetBackend(mockBackend)

			actualHashKeysRval, err := table.AddMultiple(input.HashKeys, input.DataRows)
			testutils.AssertErrorEquals(t, expected.Error, err)

			testutils.AssertEquals(t, len(expected.DataRowsStored), len(mockBackend.DataRowsInput))
			for i := range expected.DataRowsStored {
				expectedVal, ok := expected.DataRowsStored[i]["Price"].(fields.Float)
				testutils.AssertTrue(t, ok)
				actualVal, ok := mockBackend.DataRowsInput[i]["Price"].(fields.Float)
				testutils.AssertTrue(t, ok)
				testutils.AssertEquals(t, expectedVal, actualVal)
			}

			testutils.AssertEquals(t, len(expected.HashKeysStored), len(mockBackend.HashKeysInput))
			for i := range expected.HashKeysStored {
				expectedVal, ok := expected.HashKeysStored[i]["Name"].(fields.String)
				testutils.AssertTrue(t, ok)
				actualVal, ok := mockBackend.HashKeysInput[i]["Name"].(fields.String)
				testutils.AssertTrue(t, ok)
				testutils.AssertEquals(t, expectedVal, actualVal)
			}

			testutils.AssertEquals(t, len(expected.HashKeysRval), len(actualHashKeysRval))
			for i := range expected.HashKeysStored {
				testutils.AssertEquals(t, expected.HashKeysRval[i].Name, expected.HashKeysRval[i].Name)
			}
		},
	}

	tests.Run(t)
}

// func TestHashTableTransferTo(t *testing.T) {
// 	type transferInputVals struct {
// 		Error    error
// 		DataRows []mutator.MappedFieldValues
// 	}

// 	type transferExpectedVals struct {
// 		Error error
// 	}

// 	mockError := errors.New("test")

// 	tests := testutils.Tests[*transferInputVals, *transferExpectedVals]{
// 		Cases: []testutils.TestCase[*transferInputVals, *transferExpectedVals]{
// 			{
// 				Name: "transfers good values",
// 				Input: &transferInputVals{
// 					Error: nil,
// 					DataRows: []mutator.MappedFieldValues{
// 						{"Message": "test1"},
// 						{"Message": "test2"},
// 					},
// 				},
// 				Expected: &transferExpectedVals{
// 					Error: nil,
// 				},
// 			},
// 			{
// 				Name: "returns error for mismatched types",
// 				Input: &transferInputVals{
// 					Error: nil,
// 					DataRows: []mutator.MappedFieldValues{
// 						{"Message": 1},
// 					},
// 				},
// 				Expected: &transferExpectedVals{
// 					Error: mutator.SetFieldTypeError,
// 				},
// 			},
// 			{
// 				Name: "handles error from backend",
// 				Input: &transferInputVals{
// 					Error: mockError,
// 					DataRows: []mutator.MappedFieldValues{
// 						{"Message": "test"},
// 					},
// 				},
// 				Expected: &transferExpectedVals{
// 					Error: mockError,
// 				},
// 			},
// 		},
// 		Func: func(input *transferInputVals, expected *transferExpectedVals) {
// 			mockBackendSrc := &MockAppendTableBackendOps{
// 				ErrorRval:    input.Error,
// 				DataRowsRval: input.DataRows,
// 			}
// 			srcTable := product.Newproduct()
// 			srcTable.SetBackend(mockBackendSrc)

// 			mockBackendDest := &MockAppendTableBackendOps{}
// 			destTable := product.Newproduct()
// 			destTable.SetBackend(mockBackendDest)

// 			err := srcTable.TransferTo(destTable, 10)
// 			testutils.AssertErrorEquals(t, expected.Error, err)

// 			if expected.Error != nil {
// 				return
// 			}

// 			testutils.AssertEquals(t, len(input.DataRows), len(mockBackendDest.DataRowsInput))
// 			for i := range input.DataRows {
// 				expectedVal, ok := input.DataRows[i]["Message"].(fields.String)
// 				testutils.AssertTrue(t, ok)
// 				actualVal, ok := mockBackendDest.DataRowsInput[i]["Message"].(fields.String)
// 				testutils.AssertTrue(t, ok)
// 				testutils.AssertEquals(t, expectedVal, actualVal)
// 			}
// 		},
// 	}

// 	tests.Run(t)
// }
