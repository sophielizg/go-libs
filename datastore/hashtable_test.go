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

// HELPERS

func AssertProductDataRowEquals(t *testing.T, expected, actual *product.DataRow) {
	t.Helper()
	testutils.AssertEquals(t, expected.Department, actual.Department)
	testutils.AssertEquals(t, expected.Price, actual.Price)
	testutils.AssertEquals(t, expected.Quantity, actual.Quantity)
	testutils.AssertEquals(t, expected.LastUpdated, actual.LastUpdated)
}

func AssertProductDataRowFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedDepartmentVal, expectedOk := expected[product.DepartmentKey].(fields.String)
	actualDepartmentVal, ok := actual[product.DepartmentKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedDepartmentVal, actualDepartmentVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualDepartmentVal)
	}

	expectedPriceVal, expectedOk := expected[product.PriceKey].(fields.Float)
	actualPriceVal, ok := actual[product.PriceKey].(fields.Float)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedPriceVal, actualPriceVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, 0.0, actualPriceVal)
	}

	expectedQuantityVal, expectedOk := expected[product.QuantityKey].(fields.Int)
	actualQuantityVal, ok := actual[product.QuantityKey].(fields.Int)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedQuantityVal, actualQuantityVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, 0, actualQuantityVal)
	}

	expectedLastUpdatedVal, expectedOk := expected[product.LastUpdatedKey].(fields.Time)
	actualLastUpdatedVal, ok := actual[product.LastUpdatedKey].(fields.Time)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedLastUpdatedVal, actualLastUpdatedVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, fields.Time{}, actualLastUpdatedVal)
	}
}

func AssertProductHashKeyEquals(t *testing.T, expected, actual *product.HashKey) {
	t.Helper()
	testutils.AssertEquals(t, expected.Brand, actual.Brand)
	testutils.AssertEquals(t, expected.Name, actual.Name)
}

func AssertProductHashKeyFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedBrandVal, expectedOk := expected[product.BrandKey].(fields.String)
	actualBrandVal, ok := actual[product.BrandKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedBrandVal, actualBrandVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualBrandVal)
	}

	expectedNameVal, expectedOk := expected[product.NameKey].(fields.String)
	actualNameVal, ok := actual[product.NameKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedNameVal, actualNameVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualNameVal)
	}
}

// MOCKS

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

// TESTS

func TestHashTableSettings(t *testing.T) {
	table := product.NewTable()
	table.Init()
	actual := table.GetSettings()

	testutils.AssertEquals(t, "Product", actual.Name)
	testutils.AssertEquals(t, &product.DataRowSettings, actual.DataRowSettings)
	testutils.AssertEquals(t, &product.HashKeySettings, actual.HashKeySettings)

	// Check that defaults are generated
	AssertProductDataRowFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.DataRowSettings.EmptyValues)
	AssertProductHashKeyFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.HashKeySettings.EmptyValues)
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
						{product.NameKey: "test1"},
						{product.NameKey: "test2"},
					},
					DataRows: []mutator.MappedFieldValues{
						{product.PriceKey: 9.99},
						{product.PriceKey: 10.99},
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
						{product.NameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{product.PriceKey: 0},
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
						{product.NameKey: 0},
					},
					DataRows: []mutator.MappedFieldValues{
						{product.PriceKey: 9.99},
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
						{product.NameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{product.PriceKey: 9.99},
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

				AssertProductDataRowEquals(t, expectedScanFields.DataRow, actualScanFields.DataRow)
				AssertProductHashKeyEquals(t, expectedScanFields.HashKey, actualScanFields.HashKey)
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
						{product.NameKey: "test1"},
						{product.NameKey: "test2"},
					},
				},
				Expected: &addExpectedVals{
					Error: nil,
					HashKeysRval: []*product.HashKey{
						{Name: "test1"},
						{Name: "test2"},
					},
					HashKeysStored: []mutator.MappedFieldValues{
						{product.NameKey: "test1"},
						{product.NameKey: "test2"},
					},
					DataRowsStored: []mutator.MappedFieldValues{
						{product.PriceKey: 9.99},
						{product.PriceKey: 10.99},
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
				AssertProductDataRowFieldsEqualOrDefault(t, expected.DataRowsStored[i], mockBackend.DataRowsInput[i])
			}

			testutils.AssertEquals(t, len(expected.HashKeysStored), len(mockBackend.HashKeysInput))
			for i := range expected.HashKeysStored {
				AssertProductHashKeyFieldsEqualOrDefault(t, expected.HashKeysStored[i], mockBackend.HashKeysInput[i])
			}

			testutils.AssertEquals(t, len(expected.HashKeysRval), len(actualHashKeysRval))
			for i := range expected.HashKeysStored {
				AssertProductHashKeyEquals(t, expected.HashKeysRval[i], actualHashKeysRval[i])
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
