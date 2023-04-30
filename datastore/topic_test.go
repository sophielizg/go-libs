package datastore_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore/examples/shipping"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS -- defined in queue_test.go

// MOCKS

type MockTopicBackendOps struct {
	ErrorRval      error
	HasMessageRval bool
	MessageIdRval  string
	DataRowRval    mutator.MappedFieldValues
	DataRowsInput  []mutator.MappedFieldValues
}

func (t *MockTopicBackendOps) Publish(messages []mutator.MappedFieldValues) error {
	t.DataRowsInput = messages
	return t.ErrorRval
}

func (t *MockTopicBackendOps) Subscribe(subscriptionId string) error {
	return t.ErrorRval
}

func (t *MockTopicBackendOps) Unsubscribe(subscriptionId string) error {
	return t.ErrorRval
}

func (t *MockTopicBackendOps) HasMessage(subscriptionId string) (bool, error) {
	return t.HasMessageRval, t.ErrorRval
}

func (t *MockTopicBackendOps) RecieveMessage(subscriptionId string) (string, mutator.MappedFieldValues, error) {
	t.HasMessageRval = false
	return t.MessageIdRval, t.DataRowRval, t.ErrorRval
}

func (t *MockTopicBackendOps) AckSuccess(subscriptionId string, messageId string) error {
	return t.ErrorRval
}

func (t *MockTopicBackendOps) AckFailure(subscriptionId string, messageId string) error {
	return t.ErrorRval
}

// TESTS

func TestTopicSettings(t *testing.T) {
	topic := shipping.NewShippedTopic()
	topic.Init()
	actual := topic.GetSettings()

	testutils.AssertEquals(t, "Shipped", actual.Name)
	testutils.AssertEquals(t, &shipping.DataRowSettings, actual.DataSettings)

	// Check that defaults are generated
	AssertShippingDataRowFieldsEqualOrDefault(t, mutator.MappedFieldValues{}, actual.EmptyValues)
}

func TestTopicSubscribe(t *testing.T) {
	mockBackend := &MockTopicBackendOps{}
	topic := shipping.NewShippedTopic()
	topic.SetBackend(mockBackend)

	subscriptionId := "test"
	testSubscription, err := topic.Subscribe(subscriptionId)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, subscriptionId, testSubscription.Id)
	testutils.AssertEquals(t, topic, testSubscription.ParentTopic)
}

func TestTopicRecieveMessage(t *testing.T) {
	type recieveMessageInputVals struct {
		MessageId string
		DataRow   mutator.MappedFieldValues
		Error     error
	}

	type recieveMessageExpectedVals struct {
		MessageId string
		DataRow   *shipping.DataRow
		Error     error
	}

	mockError := errors.New("test")

	tests := testutils.Tests[*recieveMessageInputVals, *recieveMessageExpectedVals]{
		Cases: []testutils.TestCase[*recieveMessageInputVals, *recieveMessageExpectedVals]{
			{
				Name: "successfully recieves",
				Input: &recieveMessageInputVals{
					MessageId: "testId",
					DataRow: mutator.MappedFieldValues{
						shipping.NameKey:            "test name",
						shipping.QuantityKey:        2,
						shipping.ShippingAddressKey: "test address",
					},
					Error: nil,
				},
				Expected: &recieveMessageExpectedVals{
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
				Input: &recieveMessageInputVals{
					MessageId: "",
					DataRow:   mutator.MappedFieldValues{},
					Error:     mockError,
				},
				Expected: &recieveMessageExpectedVals{
					MessageId: "",
					DataRow:   nil,
					Error:     mockError,
				},
			},
		},
		Func: func(input *recieveMessageInputVals, expected *recieveMessageExpectedVals) {
			mockBackend := &MockTopicBackendOps{
				MessageIdRval: input.MessageId,
				DataRowRval:   input.DataRow,
				ErrorRval:     input.Error,
			}
			topic := shipping.NewShippedTopic()
			topic.SetBackend(mockBackend)

			actualMessageId, actualDataRow, err := topic.RecieveMessage("test")
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
