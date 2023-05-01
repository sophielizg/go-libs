package queries_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
	"github.com/sophielizg/go-libs/datastore/queries/queriestest"
	"github.com/sophielizg/go-libs/testutils"
)

type MockGetableBackend struct {
	ErrorRval   error
	EntriesRval []mutator.MappedFieldValues
	KeysInput   []mutator.MappedFieldValues
}

func (b *MockGetableBackend) Get(keys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	b.KeysInput = keys
	return b.EntriesRval, b.ErrorRval
}

func TestGet(t *testing.T) {
	type getInputVal struct {
		keysInput   []*queriestest.MockKey
		entriesRval []mutator.MappedFieldValues
		errorRval   error
	}

	type getExpectedVal struct {
		keysQueried   []mutator.MappedFieldValues
		entriesOutput []*queriestest.MockKeyedEntry
		errOutput     error
	}

	mockErr := errors.New("mock error")

	tests := &testutils.Tests[*getInputVal, *getExpectedVal]{
		Cases: []testutils.TestCase[*getInputVal, *getExpectedVal]{
			{
				Name: "successfully gets",
				Input: &getInputVal{
					keysInput: []*queriestest.MockKey{
						{
							Id: "test1",
						},
						{
							Id: "test2",
						},
					},
					entriesRval: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
							queriestest.IdKey:   "test1",
						},
						{
							queriestest.DataKey: "test2",
							queriestest.IdKey:   "test2",
						},
					},
					errorRval: nil,
				},
				Expected: &getExpectedVal{
					keysQueried: []mutator.MappedFieldValues{
						{
							queriestest.IdKey: "test1",
						},
						{
							queriestest.IdKey: "test2",
						},
					},
					entriesOutput: []*queriestest.MockKeyedEntry{
						{
							Data: &queriestest.MockData{
								Data: "test1",
							},
							Key: &queriestest.MockKey{
								Id: "test1",
							},
						},
						{
							Data: &queriestest.MockData{
								Data: "test2",
							},
							Key: &queriestest.MockKey{
								Id: "test2",
							},
						},
					},
					errOutput: nil,
				},
			},
			{
				Name: "returns error from backend",
				Input: &getInputVal{
					keysInput: []*queriestest.MockKey{
						{
							Id: "test1",
						},
						{
							Id: "test2",
						},
					},
					entriesRval: nil,
					errorRval:   mockErr,
				},
				Expected: &getExpectedVal{
					keysQueried: []mutator.MappedFieldValues{
						{
							queriestest.IdKey: "test1",
						},
						{
							queriestest.IdKey: "test2",
						},
					},
					entriesOutput: nil,
					errOutput:     mockErr,
				},
			},
		},
		Func: func(t *testing.T, input *getInputVal, expected *getExpectedVal) {
			backend := &MockGetableBackend{
				ErrorRval:   input.errorRval,
				EntriesRval: input.entriesRval,
			}

			getable := queries.Getable[queriestest.MockKey, *queriestest.MockKey, queriestest.MockKeyedEntry, *queriestest.MockKeyedEntry]{}
			getable.SetBackend(backend)

			actualEntries, actualErr := getable.Get(input.keysInput...)

			testutils.AssertEquals(t, len(expected.keysQueried), len(backend.KeysInput))
			for i := range expected.keysQueried {
				queriestest.AssertMockKeyFieldsEqual(t, expected.keysQueried[i], backend.KeysInput[i])
			}

			if expected.errOutput == nil {
				testutils.AssertOk(t, actualErr)
			} else {
				testutils.AssertErrorEquals(t, expected.errOutput, actualErr)
			}

			if expected.entriesOutput == nil {
				testutils.AssertTrue(t, actualEntries == nil)
			} else {
				for i := range expected.entriesOutput {
					queriestest.AssertMockKeyedEntryEquals(t, expected.entriesOutput[i], actualEntries[i])
				}
			}
		},
	}

	tests.Run(t)
}
