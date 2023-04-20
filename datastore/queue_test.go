package datastore_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore/examples/shipping"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS

func AssertShippingDataRowEquals(t *testing.T, expected, actual *shipping.DataRow) {
	t.Helper()
	testutils.AssertEquals(t, expected.Department, actual.Department)
	testutils.AssertEquals(t, expected.Brand, actual.Brand)
	testutils.AssertEquals(t, expected.Name, actual.Name)
	testutils.AssertEquals(t, expected.PurchaseTime, actual.PurchaseTime)
	testutils.AssertEquals(t, expected.Quantity, actual.Quantity)
	testutils.AssertEquals(t, expected.ShipmentTime, actual.ShipmentTime)
	testutils.AssertEquals(t, expected.ShippingAddress, actual.ShippingAddress)
}

func AssertShippingDataRowFieldsEqualOrDefault(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedDepartmentVal, expectedOk := expected[shipping.DepartmentKey].(fields.String)
	actualDepartmentVal, ok := actual[shipping.DepartmentKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedDepartmentVal, actualDepartmentVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualDepartmentVal)
	}

	expectedBrandVal, expectedOk := expected[shipping.BrandKey].(fields.String)
	actualBrandVal, ok := actual[shipping.BrandKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedBrandVal, actualBrandVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualBrandVal)
	}

	expectedNameVal, expectedOk := expected[shipping.NameKey].(fields.String)
	actualNameVal, ok := actual[shipping.NameKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedNameVal, actualNameVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualNameVal)
	}

	expectedPurchaseTimeVal, expectedOk := expected[shipping.PurchaseTimeKey].(fields.Time)
	actualPurchaseTimeVal, ok := actual[shipping.PurchaseTimeKey].(fields.Time)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedPurchaseTimeVal, actualPurchaseTimeVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, fields.Time{}, actualPurchaseTimeVal)
	}

	expectedQuantityVal, expectedOk := expected[shipping.QuantityKey].(fields.Int)
	actualQuantityVal, ok := actual[shipping.QuantityKey].(fields.Int)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedQuantityVal, actualQuantityVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, 0, actualQuantityVal)
	}

	expectedShipmentTimeVal, expectedOk := expected[shipping.ShipmentTimeKey].(fields.NullTime)
	actualShipmentTimeVal, ok := actual[shipping.ShipmentTimeKey].(fields.NullTime)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedShipmentTimeVal, actualShipmentTimeVal)
	} else {
		// Check for default value
		testutils.AssertNull(t, actualShipmentTimeVal)
	}

	expectedShippingAddressVal, expectedOk := expected[shipping.ShippingAddressKey].(fields.String)
	actualShippingAddressVal, ok := actual[shipping.ShippingAddressKey].(fields.String)
	testutils.AssertTrue(t, ok)
	if expectedOk {
		testutils.AssertEquals(t, expectedShippingAddressVal, actualShippingAddressVal)
	} else {
		// Check for default value
		testutils.AssertEquals(t, "", actualShippingAddressVal)
	}
}

// MOCKS

type MockQueueBackendOps struct {
	ErrorRval     error
	SizeRval      int
	MessageIdRval string
	DataRowRval   mutator.MappedFieldValues
	DataRowsInput []mutator.MappedFieldValues
}

func (q *MockQueueBackendOps) Size() (int, error) {
	return q.SizeRval, q.ErrorRval
}

func (q *MockQueueBackendOps) Push(messages []mutator.MappedFieldValues) error {
	q.DataRowsInput = messages
	return q.ErrorRval
}

func (q *MockQueueBackendOps) Pop() (string, mutator.MappedFieldValues, error) {
	q.SizeRval -= 1
	return q.MessageIdRval, q.DataRowRval, q.ErrorRval
}

func (q *MockQueueBackendOps) AckSuccess(messageId string) error {
	return q.ErrorRval
}

func (q *MockQueueBackendOps) AckFailure(messageId string) error {
	return q.ErrorRval
}

// TESTS

func TestQueueSettings(t *testing.T) {
	queue := shipping.NewPendingShipmentQueue()
	queue.Init()
	actual := queue.GetSettings()

	testutils.AssertEquals(t, "PendingShipment", actual.Name)
	testutils.AssertEquals(t, &shipping.DataRowSettings, actual.DataRowSettings)

	// Check that defaults are generated
	AssertShippingDataRowFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.DataRowSettings.EmptyValues)
}

func TestQueuePop(t *testing.T) {
	type popInputVals struct {
		MessageId string
		DataRow   mutator.MappedFieldValues
		Error     error
	}

	type popExpectedVals struct {
		MessageId string
		DataRow   *shipping.DataRow
		Error     error
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*popInputVals, *popExpectedVals]{
		Cases: []testutils.TestCase[*popInputVals, *popExpectedVals]{
			{
				Name: "successfully pops",
				Input: &popInputVals{
					MessageId: "testId",
					DataRow: mutator.MappedFieldValues{
						shipping.NameKey:            "test name",
						shipping.QuantityKey:        2,
						shipping.ShippingAddressKey: "test address",
					},
					Error: nil,
				},
				Expected: &popExpectedVals{
					MessageId: "testId",
					DataRow: &shipping.DataRow{
						Name:            "test name",
						Quantity:        2,
						ShippingAddress: "test address",
					},
					Error: nil,
				},
			},
			{
				Name: "returns error",
				Input: &popInputVals{
					MessageId: "",
					DataRow:   mutator.MappedFieldValues{},
					Error:     mockError,
				},
				Expected: &popExpectedVals{
					MessageId: "",
					DataRow:   nil,
					Error:     mockError,
				},
			},
		},
		Func: func(input *popInputVals, expected *popExpectedVals) {
			mockBackend := &MockQueueBackendOps{
				MessageIdRval: input.MessageId,
				DataRowRval:   input.DataRow,
				ErrorRval:     input.Error,
			}
			queue := shipping.NewPendingShipmentQueue()
			queue.SetBackend(mockBackend)

			actualMessageId, actualDataRow, err := queue.Pop()
			testutils.AssertEquals(t, expected.MessageId, actualMessageId)

			if expected.DataRow == nil {
				testutils.AssertNull(t, actualDataRow)
			} else {
				AssertShippingDataRowEquals(t, expected.DataRow, actualDataRow)
			}

			testutils.AssertErrorEquals(t, expected.Error, err)
		},
	}

	tests.Run(t)
}

func TestQueueTransferTo(t *testing.T) {
	type transferInputVals struct {
		Error   error
		DataRow mutator.MappedFieldValues
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
					DataRow: mutator.MappedFieldValues{
						shipping.NameKey: "test1",
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
					DataRow: mutator.MappedFieldValues{
						shipping.NameKey: 1,
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
					DataRow: mutator.MappedFieldValues{
						shipping.NameKey: "test",
					},
				},
				Expected: &transferExpectedVals{
					Error: mockError,
				},
			},
		},
		Func: func(input *transferInputVals, expected *transferExpectedVals) {
			mockBackendSrc := &MockQueueBackendOps{
				ErrorRval:   input.Error,
				DataRowRval: input.DataRow,
				SizeRval:    1,
			}
			srcQueue := shipping.NewPendingShipmentQueue()
			srcQueue.SetBackend(mockBackendSrc)

			mockBackendDest := &MockQueueBackendOps{}
			destQueue := shipping.NewPendingShipmentQueue()
			destQueue.SetBackend(mockBackendDest)

			err := srcQueue.TransferTo(destQueue, 10)
			testutils.AssertErrorEquals(t, expected.Error, err)

			if expected.Error != nil {
				return
			}

			testutils.AssertEquals(t, 1, len(mockBackendDest.DataRowsInput))
			AssertShippingDataRowFieldsEqualOrDefault(t, input.DataRow, mockBackendDest.DataRowsInput[0])
		},
	}

	tests.Run(t)
}
