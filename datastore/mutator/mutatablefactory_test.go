package mutator_test

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/mutator/mocks"
	"github.com/sophielizg/go-libs/testutils"
)

func TestCreateFromFields(t *testing.T) {
	factory := mutator.MutatableFactory[mocks.MockData, *mocks.MockData]{}

	testutils.Case(t, "succeccfully creates with all fields", func(t *testing.T) {
		var field3 float32
		fields := mutator.MappedFieldValues{
			"1": "test",
			"2": 1,
			"3": &field3,
		}

		actual, err := factory.CreateFromFields(fields)
		testutils.AssertOk(t, err)

		expected := &mocks.MockData{
			MockField1: "test",
			MockField2: 1,
			MockField3: &field3,
		}
		mocks.AssertMockDataEquals(t, expected, actual)
	})

	testutils.Case(t, "succeccfully creates with some fields", func(t *testing.T) {
		fields := mutator.MappedFieldValues{
			"1": "test",
		}

		actual, err := factory.CreateFromFields(fields)
		testutils.AssertOk(t, err)

		expected := &mocks.MockData{
			MockField1: "test",
		}
		mocks.AssertMockDataEquals(t, expected, actual)
	})

	testutils.Case(t, "returns error for type mismatch", func(t *testing.T) {
		fields := mutator.MappedFieldValues{
			"1": 1,
		}

		_, err := factory.CreateFromFields(fields)
		testutils.AssertErrorEquals(t, mutator.SetFieldTypeError, err)
	})

}

func TestCreateFromFieldsList(t *testing.T) {
	factory := mutator.MutatableFactory[mocks.MockData, *mocks.MockData]{}

	fieldsList := []mutator.MappedFieldValues{
		{"2": 1},
		{"2": 2},
		{"2": 3},
	}

	actual, err := factory.CreateFromFieldsList(fieldsList)
	testutils.AssertOk(t, err)

	expected := []*mocks.MockData{
		{MockField2: 1},
		{MockField2: 2},
		{MockField2: 3},
	}
	testutils.AssertEquals(t, len(expected), len(actual))

	for i := range expected {
		mocks.AssertMockDataEquals(t, expected[i], actual[i])
	}
}

func TestCreateFieldValuesList(t *testing.T) {
	factory := mutator.MutatableFactory[mocks.MockData, *mocks.MockData]{}

	dataList := []*mocks.MockData{
		{MockField2: 1},
		{MockField2: 2},
		{MockField2: 3},
	}

	actual := factory.CreateFieldValuesList(dataList)
	testutils.AssertEquals(t, len(dataList), len(actual))

	for i := range actual {
		actualField1, ok := actual[i]["1"].(string)
		testutils.AssertTrue(t, ok)
		testutils.AssertEquals(t, "", actualField1)

		actualField2, ok := actual[i]["2"].(int)
		testutils.AssertTrue(t, ok)
		testutils.AssertEquals(t, i+1, actualField2)
	}
}
