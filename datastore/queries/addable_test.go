package queries_test

import (
	"errors"
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
	"github.com/sophielizg/go-libs/datastore/queries/queriestest"
	"github.com/sophielizg/go-libs/testutils"
)

type MockAddableBackend struct {
	ErrorRval    error
	EntriesRval  []mutator.MappedFieldValues
	EntriesInput []mutator.MappedFieldValues
}

func (b *MockAddableBackend) Add(entries []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	b.EntriesInput = entries
	return b.EntriesRval, b.ErrorRval
}

func TestAdd(t *testing.T) {
	type addInputVal struct {
		entriesInput []*queriestest.MockNonKeyedEntry
		entriesRval  []mutator.MappedFieldValues
		errorRval    error
	}

	type addExpectedVal struct {
		entriesSaved  []mutator.MappedFieldValues
		entriesOutput []*queriestest.MockNonKeyedEntry
		errOutput     error
	}

	mockErr := errors.New("mock error")

	tests := &testutils.Tests[*addInputVal, *addExpectedVal]{
		Cases: []testutils.TestCase[*addInputVal, *addExpectedVal]{
			{
				Name: "successfully adds",
				Input: &addInputVal{
					entriesInput: []*queriestest.MockNonKeyedEntry{
						{
							Data: &queriestest.MockData{
								Data: "test1",
							},
						},
						{
							Data: &queriestest.MockData{
								Data: "test2",
							},
						},
					},
					entriesRval: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
						},
						{
							queriestest.DataKey: "test2",
						},
					},
					errorRval: nil,
				},
				Expected: &addExpectedVal{
					entriesSaved: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
						},
						{
							queriestest.DataKey: "test2",
						},
					},
					entriesOutput: []*queriestest.MockNonKeyedEntry{
						{
							Data: &queriestest.MockData{
								Data: "test1",
							},
						},
						{
							Data: &queriestest.MockData{
								Data: "test2",
							},
						},
					},
					errOutput: nil,
				},
			},
			{
				Name: "returns error from backend",
				Input: &addInputVal{
					entriesInput: []*queriestest.MockNonKeyedEntry{
						{
							Data: &queriestest.MockData{
								Data: "test1",
							},
						},
						{
							Data: &queriestest.MockData{
								Data: "test2",
							},
						},
					},
					entriesRval: nil,
					errorRval:   mockErr,
				},
				Expected: &addExpectedVal{
					entriesSaved: []mutator.MappedFieldValues{
						{
							queriestest.DataKey: "test1",
						},
						{
							queriestest.DataKey: "test2",
						},
					},
					entriesOutput: nil,
					errOutput:     mockErr,
				},
			},
		},
		Func: func(t *testing.T, input *addInputVal, expected *addExpectedVal) {
			backend := &MockAddableBackend{
				ErrorRval:   input.errorRval,
				EntriesRval: input.entriesRval,
			}

			addable := queries.Addable[queriestest.MockNonKeyedEntry, *queriestest.MockNonKeyedEntry]{}
			addable.SetBackend(backend)

			actualEntries, actualErr := addable.Add(input.entriesInput...)

			testutils.AssertEquals(t, len(expected.entriesSaved), len(backend.EntriesInput))
			for i := range expected.entriesSaved {
				queriestest.AssertMockNonKeyedEntryFieldsEqual(t, expected.entriesSaved[i], backend.EntriesInput[i])
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
					queriestest.AssertMockNonKeyedEntryEquals(t, expected.entriesOutput[i], actualEntries[i])
				}
			}
		},
	}

	tests.Run(t)
}
