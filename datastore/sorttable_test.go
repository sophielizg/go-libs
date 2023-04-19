package datastore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/compare"
	"github.com/sophielizg/go-libs/datastore/examples/purchase"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS

func AssertPurchaseDataRowEquals(t *testing.T, expected, actual *purchase.DataRow) {
	t.Helper()
	testutils.AssertEquals(t, expected.Department, actual.Department)
	testutils.AssertEquals(t, expected.Price, actual.Price)
	testutils.AssertEquals(t, expected.Quantity, actual.Quantity)
	testutils.AssertEquals(t, expected.LastUpdated, actual.LastUpdated)
}

func AssertPurchaseDataRowFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedDepartmentVal, expectedOk := expected[purchase.DepartmentKey].(fields.String)
	actualDepartmentVal, ok := actual[purchase.DepartmentKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedDepartmentVal, actualDepartmentVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualDepartmentVal)
	}

	expectedPriceVal, expectedOk := expected[purchase.PriceKey].(fields.Float)
	actualPriceVal, ok := actual[purchase.PriceKey].(fields.Float)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedPriceVal, actualPriceVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, 0.0, actualPriceVal)
	}

	expectedQuantityVal, expectedOk := expected[purchase.QuantityKey].(fields.Int)
	actualQuantityVal, ok := actual[purchase.QuantityKey].(fields.Int)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedQuantityVal, actualQuantityVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, 0, actualQuantityVal)
	}

	expectedLastUpdatedVal, expectedOk := expected[purchase.LastUpdatedKey].(fields.Time)
	actualLastUpdatedVal, ok := actual[purchase.LastUpdatedKey].(fields.Time)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedLastUpdatedVal, actualLastUpdatedVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, fields.Time{}, actualLastUpdatedVal)
	}
}

func AssertPurchaseHashKeyEquals(t *testing.T, expected, actual *purchase.HashKey) {
	t.Helper()
	testutils.AssertEquals(t, expected.CustomerName, actual.CustomerName)
}

func AssertPurchaseHashKeyFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedBrandVal, expectedOk := expected[purchase.CustomerNameKey].(fields.String)
	actualBrandVal, ok := actual[purchase.CustomerNameKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedBrandVal, actualBrandVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualBrandVal)
	}
}

func AssertPurchaseSortKeyEquals(t *testing.T, expected, actual *purchase.SortKey) {
	t.Helper()
	testutils.AssertEquals(t, expected.PurchaseTime, actual.PurchaseTime)
	testutils.AssertEquals(t, expected.ItemBrand, actual.ItemBrand)
	testutils.AssertEquals(t, expected.ItemName, actual.ItemName)
}

func AssertPurchaseSortKeyFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedPurchaseTimeVal, expectedOk := expected[purchase.PurchaseTimeKey].(fields.Time)
	actualPurchaseTimeVal, ok := actual[purchase.PurchaseTimeKey].(fields.Time)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedPurchaseTimeVal, actualPurchaseTimeVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, fields.Time{}, actualPurchaseTimeVal)
	}

	expectedItemBrandVal, expectedOk := expected[purchase.ItemBrandKey].(fields.String)
	actualItemBrandVal, ok := actual[purchase.ItemBrandKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedItemBrandVal, actualItemBrandVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualItemBrandVal)
	}

	expectedItemNameVal, expectedOk := expected[purchase.ItemNameKey].(fields.String)
	actualItemNameVal, ok := actual[purchase.ItemNameKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedItemNameVal, actualItemNameVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualItemNameVal)
	}
}

func AssertPurchaseSortKeyComparatorEquals(t *testing.T, expected, actual *purchase.SortKeyComparator) {
	t.Helper()
	testutils.AssertTrue(t, expected.PurchaseTime.Equals(actual.PurchaseTime))
	testutils.AssertTrue(t, expected.ItemBrand.Equals(actual.ItemBrand))
	testutils.AssertTrue(t, expected.ItemName.Equals(actual.ItemName))
}

func AssertPurchaseSortKeyComparatorFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedPurchaseTimeVal, expectedOk := expected[purchase.PurchaseTimeKey].(*compare.Comparator[fields.Time])
	actualPurchaseTimeVal, ok := actual[purchase.PurchaseTimeKey].(*compare.Comparator[fields.Time])
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertTrue(t, expectedPurchaseTimeVal.Equals(actualPurchaseTimeVal))
	} else {
		// Check for default value
		testutils.AssertNull(t, actualPurchaseTimeVal)
	}

	expectedItemBrandVal, expectedOk := expected[purchase.ItemBrandKey].(*compare.Comparator[fields.String])
	actualItemBrandVal, ok := actual[purchase.ItemBrandKey].(*compare.Comparator[fields.String])
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertTrue(t, expectedItemBrandVal.Equals(actualItemBrandVal))
	} else {
		// Check for default value
		testutils.AssertNull(t, actualItemBrandVal)
	}

	expectedItemNameVal, expectedOk := expected[purchase.ItemNameKey].(*compare.Comparator[fields.String])
	actualItemNameVal, ok := actual[purchase.ItemNameKey].(*compare.Comparator[fields.String])
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertTrue(t, expectedItemNameVal.Equals(actualItemNameVal))
	} else {
		// Check for default value
		testutils.AssertNull(t, actualItemNameVal)
	}
}

// MOCKS

type MockSortTableBackendOps struct {
	ErrorRval       error
	DataRowsRval    []mutator.MappedFieldValues
	HashKeysRval    []mutator.MappedFieldValues
	SortKeysRval    []mutator.MappedFieldValues
	DataRowsInput   []mutator.MappedFieldValues
	HashKeysInput   []mutator.MappedFieldValues
	SortKeysInput   []mutator.MappedFieldValues
	ComparatorInput mutator.MappedFieldValues
}

func (b *MockSortTableBackendOps) Scan(batchSize int) (chan *datastore.ScanFields, chan error) {
	dataChan := make(chan *datastore.ScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errorChan)

		for i := range b.DataRowsRval {
			dataChan <- &datastore.ScanFields{
				DataRow: b.DataRowsRval[i],
				HashKey: b.HashKeysRval[i],
				SortKey: b.SortKeysRval[i],
			}
		}

		if b.ErrorRval != nil {
			errorChan <- b.ErrorRval
		}
	}()

	return dataChan, errorChan
}

func (b *MockSortTableBackendOps) GetMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	return b.DataRowsRval, b.ErrorRval
}

func (b *MockSortTableBackendOps) AddMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, []mutator.MappedFieldValues, error) {
	b.DataRowsInput = data
	b.HashKeysInput = hashKeys
	b.SortKeysInput = sortKeys
	return b.HashKeysRval, b.SortKeysRval, b.ErrorRval
}

func (b *MockSortTableBackendOps) UpdateMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) error {
	b.DataRowsInput = data
	b.HashKeysInput = hashKeys
	b.SortKeysInput = sortKeys
	return b.ErrorRval
}

func (b *MockSortTableBackendOps) DeleteMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues) error {
	return b.ErrorRval
}

func (b *MockSortTableBackendOps) GetWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues) ([]mutator.MappedFieldValues, []mutator.MappedFieldValues, error) {
	b.ComparatorInput = comparator
	return b.DataRowsRval, b.SortKeysRval, b.ErrorRval
}

func (b *MockSortTableBackendOps) UpdateWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues, data mutator.MappedFieldValues) error {
	b.ComparatorInput = comparator
	b.DataRowsInput = []mutator.MappedFieldValues{data}
	return b.ErrorRval
}

func (b *MockSortTableBackendOps) DeleteWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues) error {
	b.ComparatorInput = comparator
	return b.ErrorRval
}

// TESTS

func TestSortTableSettings(t *testing.T) {
	table := purchase.NewTable()
	table.Init()
	actual := table.GetSettings()

	testutils.AssertEquals(t, "Purchase", actual.Name)
	testutils.AssertEquals(t, &purchase.DataRowSettings, actual.DataRowSettings)
	testutils.AssertEquals(t, &purchase.HashKeySettings, actual.HashKeySettings)
	testutils.AssertEquals(t, &purchase.SortKeySettings, actual.SortKeySettings)

	// Check that defaults are generated
	AssertPurchaseDataRowFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.DataRowSettings.EmptyValues)
	AssertPurchaseHashKeyFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.HashKeySettings.EmptyValues)
	AssertPurchaseSortKeyFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.SortKeySettings.EmptyValues)
}

func TestSortTableScan(t *testing.T) {
	type scanInputVals struct {
		Error    error
		DataRows []mutator.MappedFieldValues
		HashKeys []mutator.MappedFieldValues
		SortKeys []mutator.MappedFieldValues
	}

	type scanExpectedVals struct {
		Errors     []error
		ScanFields []purchase.TableScan
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*scanInputVals, *scanExpectedVals]{
		Cases: []testutils.TestCase[*scanInputVals, *scanExpectedVals]{
			{
				Name: "properly formats good values",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []purchase.TableScan{
						{
							DataRow: &purchase.DataRow{Price: 9.99},
							HashKey: &purchase.HashKey{CustomerName: "test1"},
							SortKey: &purchase.SortKey{ItemName: "test1"},
						},
						{
							DataRow: &purchase.DataRow{Price: 10.99},
							HashKey: &purchase.HashKey{CustomerName: "test2"},
							SortKey: &purchase.SortKey{ItemName: "test2"},
						},
					},
				},
			},
			{
				Name: "returns error for data row mismatched types",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 0},
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
						{purchase.CustomerNameKey: 0},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
					},
				},
				Expected: &scanExpectedVals{
					Errors: []error{mutator.SetFieldTypeError},
				},
			},
			{
				Name: "returns error for sort key mismatched types",
				Input: &scanInputVals{
					Error: nil,
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: 0},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
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
						{purchase.CustomerNameKey: "test1"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
					},
				},
				Expected: &scanExpectedVals{
					ScanFields: []purchase.TableScan{
						{
							DataRow: &purchase.DataRow{Price: 9.99},
							HashKey: &purchase.HashKey{CustomerName: "test1"},
							SortKey: &purchase.SortKey{ItemName: "test1"},
						},
					},
					Errors: []error{mockError},
				},
			},
		},
		Func: func(input *scanInputVals, expected *scanExpectedVals) {
			mockBackend := &MockSortTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRows,
				HashKeysRval: input.HashKeys,
				SortKeysRval: input.SortKeys,
			}
			table := purchase.NewTable()
			table.SetBackend(mockBackend)

			actualScanFieldsChan, actualErrorChan := table.Scan(10)

			for _, expectedScanFields := range expected.ScanFields {
				actualScanFields, more := <-actualScanFieldsChan
				if !more {
					t.Errorf("actualScanFieldsChan ended prematurely")
				}

				AssertPurchaseDataRowEquals(t, expectedScanFields.DataRow, actualScanFields.DataRow)
				AssertPurchaseHashKeyEquals(t, expectedScanFields.HashKey, actualScanFields.HashKey)
				AssertPurchaseSortKeyEquals(t, expectedScanFields.SortKey, actualScanFields.SortKey)
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

func TestSortTableGetMultiple(t *testing.T) {
	type getInputVals struct {
		Error        error
		HashKeys     []*purchase.HashKey
		SortKeys     []*purchase.SortKey
		DataRowsRval []mutator.MappedFieldValues
	}

	type getExpectedVals struct {
		Error    error
		DataRows []*purchase.DataRow
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*getInputVals, *getExpectedVals]{
		Cases: []testutils.TestCase[*getInputVals, *getExpectedVals]{
			{
				Name: "successfully gets",
				Input: &getInputVals{
					Error: nil,
					HashKeys: []*purchase.HashKey{
						{CustomerName: "test1"},
						{CustomerName: "test2"},
					},
					SortKeys: []*purchase.SortKey{
						{ItemName: "test1"},
						{ItemName: "test2"},
					},
					DataRowsRval: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
					},
				},
				Expected: &getExpectedVals{
					Error: nil,
					DataRows: []*purchase.DataRow{
						{Price: 9.99},
						{Price: 10.99},
					},
				},
			},
			{
				Name: "returns error",
				Input: &getInputVals{
					Error:        mockError,
					HashKeys:     []*purchase.HashKey{},
					SortKeys:     []*purchase.SortKey{},
					DataRowsRval: []mutator.MappedFieldValues{},
				},
				Expected: &getExpectedVals{
					Error:    mockError,
					DataRows: nil,
				},
			},
		},
		Func: func(input *getInputVals, expected *getExpectedVals) {
			mockBackend := &MockSortTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRowsRval,
			}
			table := purchase.NewTable()
			table.SetBackend(mockBackend)

			actualDataRows, err := table.GetMultiple(input.HashKeys, input.SortKeys)
			testutils.AssertErrorEquals(t, expected.Error, err)

			testutils.AssertEquals(t, len(expected.DataRows), len(actualDataRows))
			for i := range expected.DataRows {
				AssertPurchaseDataRowEquals(t, expected.DataRows[i], actualDataRows[i])
			}
		},
	}

	tests.Run(t)
}

func TestSortTableAddMultiple(t *testing.T) {
	type addInputVals struct {
		Error        error
		HashKeys     []*purchase.HashKey
		SortKeys     []*purchase.SortKey
		DataRows     []*purchase.DataRow
		HashKeysRval []mutator.MappedFieldValues
		SortKeysRval []mutator.MappedFieldValues
	}

	type addExpectedVals struct {
		Error          error
		HashKeysRval   []*purchase.HashKey
		SortKeysRval   []*purchase.SortKey
		HashKeysStored []mutator.MappedFieldValues
		SortKeysStored []mutator.MappedFieldValues
		DataRowsStored []mutator.MappedFieldValues
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*addInputVals, *addExpectedVals]{
		Cases: []testutils.TestCase[*addInputVals, *addExpectedVals]{
			{
				Name: "successfully adds",
				Input: &addInputVals{
					Error: nil,
					HashKeys: []*purchase.HashKey{
						{CustomerName: "test1"},
						{CustomerName: "test2"},
					},
					SortKeys: []*purchase.SortKey{
						{ItemName: "test1"},
						{ItemName: "test2"},
					},
					DataRows: []*purchase.DataRow{
						{Price: 9.99},
						{Price: 10.99},
					},
					HashKeysRval: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeysRval: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
				},
				Expected: &addExpectedVals{
					Error: nil,
					HashKeysRval: []*purchase.HashKey{
						{CustomerName: "test1"},
						{CustomerName: "test2"},
					},
					SortKeysRval: []*purchase.SortKey{
						{ItemName: "test1"},
						{ItemName: "test2"},
					},
					HashKeysStored: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeysStored: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
					DataRowsStored: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
					},
				},
			},
			{
				Name: "returns error",
				Input: &addInputVals{
					Error:        mockError,
					DataRows:     []*purchase.DataRow{},
					SortKeys:     []*purchase.SortKey{},
					HashKeys:     []*purchase.HashKey{},
					HashKeysRval: []mutator.MappedFieldValues{},
					SortKeysRval: []mutator.MappedFieldValues{},
				},
				Expected: &addExpectedVals{
					Error:          mockError,
					HashKeysRval:   nil,
					SortKeysRval:   nil,
					HashKeysStored: []mutator.MappedFieldValues{},
					SortKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
			{
				Name: "returns mismatch length error for hash keys",
				Input: &addInputVals{
					Error:    nil,
					DataRows: []*purchase.DataRow{},
					SortKeys: []*purchase.SortKey{},
					HashKeys: []*purchase.HashKey{},
					HashKeysRval: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
					},
					SortKeysRval: []mutator.MappedFieldValues{},
				},
				Expected: &addExpectedVals{
					Error:          datastore.OutputLengthMismatchError,
					HashKeysRval:   nil,
					SortKeysRval:   nil,
					HashKeysStored: []mutator.MappedFieldValues{},
					SortKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
			{
				Name: "returns mismatch length error for sort keys",
				Input: &addInputVals{
					Error:        nil,
					DataRows:     []*purchase.DataRow{},
					SortKeys:     []*purchase.SortKey{},
					HashKeys:     []*purchase.HashKey{},
					HashKeysRval: []mutator.MappedFieldValues{},
					SortKeysRval: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
					},
				},
				Expected: &addExpectedVals{
					Error:          datastore.OutputLengthMismatchError,
					HashKeysRval:   nil,
					SortKeysRval:   nil,
					HashKeysStored: []mutator.MappedFieldValues{},
					SortKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
		},
		Func: func(input *addInputVals, expected *addExpectedVals) {
			mockBackend := &MockSortTableBackendOps{
				ErrorRval:    input.Error,
				HashKeysRval: input.HashKeysRval,
				SortKeysRval: input.SortKeysRval,
			}
			table := purchase.NewTable()
			table.SetBackend(mockBackend)

			actualHashKeysRval, actualSortKeysRval, err := table.AddMultiple(input.HashKeys, input.SortKeys, input.DataRows)
			testutils.AssertErrorEquals(t, expected.Error, err)

			testutils.AssertEquals(t, len(expected.DataRowsStored), len(mockBackend.DataRowsInput))
			for i := range expected.DataRowsStored {
				AssertPurchaseDataRowFieldsEqualOrDefault(t, expected.DataRowsStored[i], mockBackend.DataRowsInput[i])
			}

			testutils.AssertEquals(t, len(expected.HashKeysStored), len(mockBackend.HashKeysInput))
			for i := range expected.HashKeysStored {
				AssertPurchaseHashKeyFieldsEqualOrDefault(t, expected.HashKeysStored[i], mockBackend.HashKeysInput[i])
			}

			testutils.AssertEquals(t, len(expected.SortKeysStored), len(mockBackend.SortKeysInput))
			for i := range expected.SortKeysStored {
				AssertPurchaseSortKeyFieldsEqualOrDefault(t, expected.SortKeysStored[i], mockBackend.SortKeysInput[i])
			}

			testutils.AssertEquals(t, len(expected.HashKeysRval), len(actualHashKeysRval))
			for i := range expected.HashKeysRval {
				AssertPurchaseHashKeyEquals(t, expected.HashKeysRval[i], actualHashKeysRval[i])
			}

			testutils.AssertEquals(t, len(expected.SortKeysRval), len(actualSortKeysRval))
			for i := range expected.SortKeysRval {
				AssertPurchaseSortKeyEquals(t, expected.SortKeysRval[i], actualSortKeysRval[i])
			}
		},
	}

	tests.Run(t)
}

func TestSortTableUpdateMultiple(t *testing.T) {
	type updateInputVals struct {
		Error    error
		HashKeys []*purchase.HashKey
		SortKeys []*purchase.SortKey
		DataRows []*purchase.DataRow
	}

	type updateExpectedVals struct {
		Error          error
		HashKeysStored []mutator.MappedFieldValues
		SortKeysStored []mutator.MappedFieldValues
		DataRowsStored []mutator.MappedFieldValues
	}

	tests := testutils.Tests[*updateInputVals, *updateExpectedVals]{
		Cases: []testutils.TestCase[*updateInputVals, *updateExpectedVals]{
			{
				Name: "successfully updates",
				Input: &updateInputVals{
					Error: nil,
					HashKeys: []*purchase.HashKey{
						{CustomerName: "test1"},
						{CustomerName: "test2"},
					},
					SortKeys: []*purchase.SortKey{
						{ItemName: "test1"},
						{ItemName: "test2"},
					},
					DataRows: []*purchase.DataRow{
						{Price: 9.99},
						{Price: 10.99},
					},
				},
				Expected: &updateExpectedVals{
					Error: nil,
					HashKeysStored: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeysStored: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
					DataRowsStored: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
					},
				},
			},
			{
				Name: "returns mismatch length error for hash keys",
				Input: &updateInputVals{
					Error:    nil,
					DataRows: []*purchase.DataRow{},
					HashKeys: []*purchase.HashKey{
						{CustomerName: "test"},
					},
					SortKeys: []*purchase.SortKey{},
				},
				Expected: &updateExpectedVals{
					Error:          datastore.InputLengthMismatchError,
					HashKeysStored: []mutator.MappedFieldValues{},
					SortKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
			{
				Name: "returns mismatch length error for sort keys",
				Input: &updateInputVals{
					Error:    nil,
					DataRows: []*purchase.DataRow{},
					HashKeys: []*purchase.HashKey{},
					SortKeys: []*purchase.SortKey{
						{ItemName: "test1"},
					},
				},
				Expected: &updateExpectedVals{
					Error:          datastore.InputLengthMismatchError,
					HashKeysStored: []mutator.MappedFieldValues{},
					SortKeysStored: []mutator.MappedFieldValues{},
					DataRowsStored: []mutator.MappedFieldValues{},
				},
			},
		},
		Func: func(input *updateInputVals, expected *updateExpectedVals) {
			mockBackend := &MockSortTableBackendOps{
				ErrorRval: input.Error,
			}
			table := purchase.NewTable()
			table.SetBackend(mockBackend)

			err := table.UpdateMultiple(input.HashKeys, input.SortKeys, input.DataRows)
			testutils.AssertErrorEquals(t, expected.Error, err)

			testutils.AssertEquals(t, len(expected.DataRowsStored), len(mockBackend.DataRowsInput))
			for i := range expected.DataRowsStored {
				AssertPurchaseDataRowFieldsEqualOrDefault(t, expected.DataRowsStored[i], mockBackend.DataRowsInput[i])
			}

			testutils.AssertEquals(t, len(expected.HashKeysStored), len(mockBackend.HashKeysInput))
			for i := range expected.HashKeysStored {
				AssertPurchaseHashKeyFieldsEqualOrDefault(t, expected.HashKeysStored[i], mockBackend.HashKeysInput[i])
			}

			testutils.AssertEquals(t, len(expected.SortKeysStored), len(mockBackend.SortKeysInput))
			for i := range expected.SortKeysStored {
				AssertPurchaseSortKeyFieldsEqualOrDefault(t, expected.SortKeysStored[i], mockBackend.SortKeysInput[i])
			}
		},
	}

	tests.Run(t)
}

func TestSortTableGetWithSortComparator(t *testing.T) {
	type getInputVals struct {
		Error        error
		HashKey      *purchase.HashKey
		Comparator   *purchase.SortKeyComparator
		DataRowsRval []mutator.MappedFieldValues
		SortKeysRval []mutator.MappedFieldValues
	}

	type getExpectedVals struct {
		Error          error
		ComparatorUsed mutator.MappedFieldValues
		DataRows       []*purchase.DataRow
		SortKeys       []*purchase.SortKey
	}

	mockError := errors.New("test")
	currentTime := time.Now()

	tests := testutils.Tests[*getInputVals, *getExpectedVals]{
		Cases: []testutils.TestCase[*getInputVals, *getExpectedVals]{
			{
				Name: "successfully gets",
				Input: &getInputVals{
					Error:   nil,
					HashKey: &purchase.HashKey{CustomerName: "test1"},
					Comparator: &purchase.SortKeyComparator{
						PurchaseTime: compare.Eq(currentTime),
					},
					DataRowsRval: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
					},
					SortKeysRval: []mutator.MappedFieldValues{
						{purchase.PurchaseTimeKey: currentTime},
						{purchase.PurchaseTimeKey: currentTime},
					},
				},
				Expected: &getExpectedVals{
					Error: nil,
					ComparatorUsed: mutator.MappedFieldValues{
						purchase.PurchaseTimeKey: &compare.Comparator[fields.Time]{
							Op:     compare.EQ,
							Values: []fields.Time{currentTime},
						},
					},
					DataRows: []*purchase.DataRow{
						{Price: 9.99},
						{Price: 10.99},
					},
					SortKeys: []*purchase.SortKey{
						{PurchaseTime: currentTime},
						{PurchaseTime: currentTime},
					},
				},
			},
			{
				Name: "enforces error for wrong sort key order",
				Input: &getInputVals{
					Error:   nil,
					HashKey: &purchase.HashKey{CustomerName: "test1"},
					Comparator: &purchase.SortKeyComparator{
						ItemName: compare.Eq("test"),
					},
					DataRowsRval: []mutator.MappedFieldValues{},
					SortKeysRval: []mutator.MappedFieldValues{},
				},
				Expected: &getExpectedVals{
					Error:          datastore.ComparatorMissingFieldsError,
					ComparatorUsed: nil,
					DataRows:       nil,
					SortKeys:       nil,
				},
			},
			{
				Name: "returns error for data row and sort key mismatch",
				Input: &getInputVals{
					Error:   nil,
					HashKey: &purchase.HashKey{CustomerName: "test1"},
					Comparator: &purchase.SortKeyComparator{
						PurchaseTime: compare.Eq(currentTime),
					},
					DataRowsRval: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
					},
					SortKeysRval: []mutator.MappedFieldValues{
						{purchase.PurchaseTimeKey: currentTime},
						{purchase.PurchaseTimeKey: currentTime},
					},
				},
				Expected: &getExpectedVals{
					Error:          datastore.OutputLengthMismatchError,
					ComparatorUsed: nil,
					DataRows:       nil,
					SortKeys:       nil,
				},
			},
			{
				Name: "returns error",
				Input: &getInputVals{
					Error:        mockError,
					HashKey:      &purchase.HashKey{},
					Comparator:   &purchase.SortKeyComparator{},
					DataRowsRval: []mutator.MappedFieldValues{},
					SortKeysRval: []mutator.MappedFieldValues{},
				},
				Expected: &getExpectedVals{
					Error:          mockError,
					ComparatorUsed: nil,
					DataRows:       nil,
					SortKeys:       nil,
				},
			},
		},
		Func: func(input *getInputVals, expected *getExpectedVals) {
			mockBackend := &MockSortTableBackendOps{
				ErrorRval:    input.Error,
				DataRowsRval: input.DataRowsRval,
				SortKeysRval: input.SortKeysRval,
			}
			table := purchase.NewTable()
			table.SetBackend(mockBackend)

			actualDataRows, actualSortKeys, err := table.GetWithSortComparator(input.HashKey, input.Comparator)
			testutils.AssertErrorEquals(t, expected.Error, err)

			if expected.ComparatorUsed != nil {
				AssertPurchaseSortKeyComparatorFieldsEqualOrDefault(t, expected.ComparatorUsed, mockBackend.ComparatorInput)
			}

			testutils.AssertEquals(t, len(expected.DataRows), len(actualDataRows))
			for i := range expected.DataRows {
				AssertPurchaseDataRowEquals(t, expected.DataRows[i], actualDataRows[i])
			}

			testutils.AssertEquals(t, len(expected.SortKeys), len(actualSortKeys))
			for i := range expected.SortKeys {
				AssertPurchaseSortKeyEquals(t, expected.SortKeys[i], actualSortKeys[i])
			}
		},
	}

	tests.Run(t)
}

func TestSortTableTransferTo(t *testing.T) {
	type transferInputVals struct {
		Error    error
		HashKeys []mutator.MappedFieldValues
		SortKeys []mutator.MappedFieldValues
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
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
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
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: 0},
						{purchase.CustomerNameKey: "test2"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
						{purchase.ItemNameKey: "test2"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
						{purchase.PriceKey: 10.99},
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
					HashKeys: []mutator.MappedFieldValues{
						{purchase.CustomerNameKey: "test1"},
					},
					SortKeys: []mutator.MappedFieldValues{
						{purchase.ItemNameKey: "test1"},
					},
					DataRows: []mutator.MappedFieldValues{
						{purchase.PriceKey: 9.99},
					},
				},
				Expected: &transferExpectedVals{
					Error: mockError,
				},
			},
		},
		Func: func(input *transferInputVals, expected *transferExpectedVals) {
			mockBackendSrc := &MockSortTableBackendOps{
				ErrorRval:    input.Error,
				HashKeysRval: input.HashKeys,
				SortKeysRval: input.SortKeys,
				DataRowsRval: input.DataRows,
			}
			srcTable := purchase.NewTable()
			srcTable.SetBackend(mockBackendSrc)

			mockBackendDest := &MockSortTableBackendOps{
				HashKeysRval: input.HashKeys,
				SortKeysRval: input.SortKeys,
			}
			destTable := purchase.NewTable()
			destTable.SetBackend(mockBackendDest)

			err := srcTable.TransferTo(destTable, 10)
			testutils.AssertErrorEquals(t, expected.Error, err)

			if expected.Error != nil {
				return
			}

			testutils.AssertEquals(t, len(input.HashKeys), len(mockBackendDest.HashKeysInput))
			for i := range input.HashKeys {
				AssertPurchaseHashKeyFieldsEqualOrDefault(t, input.HashKeys[i], mockBackendDest.HashKeysInput[i])
			}

			testutils.AssertEquals(t, len(input.SortKeys), len(mockBackendDest.SortKeysInput))
			for i := range input.SortKeys {
				AssertPurchaseSortKeyFieldsEqualOrDefault(t, input.SortKeys[i], mockBackendDest.SortKeysInput[i])
			}

			testutils.AssertEquals(t, len(input.DataRows), len(mockBackendDest.DataRowsInput))
			for i := range input.DataRows {
				AssertPurchaseDataRowFieldsEqualOrDefault(t, input.DataRows[i], mockBackendDest.DataRowsInput[i])
			}
		},
	}

	tests.Run(t)
}
