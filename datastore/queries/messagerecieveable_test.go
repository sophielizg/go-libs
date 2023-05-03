package queries_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
	"github.com/sophielizg/go-libs/datastore/queries/queriestest"
	"github.com/sophielizg/go-libs/testutils"
)

type MockMessageRecieveableBackend struct {
	ErrorRval    error
	MessagesRval []mutator.MappedFieldValues
	messageIdx   int
}

func (b *MockMessageRecieveableBackend) HasMessage() (bool, error) {
	return b.messageIdx < len(b.MessagesRval), b.ErrorRval
}

func (b *MockMessageRecieveableBackend) RecieveMessage() (string, mutator.MappedFieldValues, error) {
	id := strconv.Itoa(b.messageIdx)
	message := b.MessagesRval[b.messageIdx]
	return id, message, b.ErrorRval
}

func (b *MockMessageRecieveableBackend) AckSuccess(messageId []string) error {
	return b.ErrorRval
}

func (b *MockMessageRecieveableBackend) AckFailure(messageId []string) error {
	return b.ErrorRval
}

func TestRecieveMessage(t *testing.T) {
	type recieveInputVal struct {
		messagesRval []mutator.MappedFieldValues
		errorRval    error
	}

	type recieveExpectedVal struct {
		messageIdOutput string
		messageOutput   *queriestest.MockNonKeyedEntry
		errOutput       error
	}

	mockErr := errors.New("mock error")

	tests := &testutils.Tests[*recieveInputVal, *recieveExpectedVal]{
		Cases: []testutils.TestCase[*recieveInputVal, *recieveExpectedVal]{
			{
				Name: "successfully recieves message",
				Input: &recieveInputVal{
					messagesRval: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
						},
					},
					errorRval: nil,
				},
				Expected: &recieveExpectedVal{
					messageIdOutput: "0",
					messageOutput: &queriestest.MockNonKeyedEntry{
						Data: &queriestest.MockData{
							Data: "test1",
						},
					},
					errOutput: nil,
				},
			},
			{
				Name: "returns error from backend",
				Input: &recieveInputVal{
					messagesRval: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
						},
					},
					errorRval: mockErr,
				},
				Expected: &recieveExpectedVal{
					messageIdOutput: "",
					messageOutput:   nil,
					errOutput:       mockErr,
				},
			},
		},
		Func: func(t *testing.T, input *recieveInputVal, expected *recieveExpectedVal) {
			backend := &MockMessageRecieveableBackend{
				ErrorRval:    input.errorRval,
				MessagesRval: input.messagesRval,
			}

			recieveable := queries.MessageReceiveable[queriestest.MockNonKeyedEntry, *queriestest.MockNonKeyedEntry]{}
			recieveable.SetBackend(backend)

			actualId, actualMessage, actualErr := recieveable.RecieveMessage()

			testutils.AssertEquals(t, expected.messageIdOutput, actualId)

			if expected.errOutput == nil {
				testutils.AssertOk(t, actualErr)
			} else {
				testutils.AssertErrorEquals(t, expected.errOutput, actualErr)
			}

			if expected.messageOutput == nil {
				testutils.AssertTrue(t, actualMessage == nil)
			} else {
				queriestest.AssertMockNonKeyedEntryEquals(t, expected.messageOutput, actualMessage)
			}
		},
	}

	tests.Run(t)
}
