package mutator_test

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/mutator/mocks"
	"github.com/sophielizg/go-libs/testutils"
)

func TestGetField(t *testing.T) {
	mockData := &mocks.MockData{
		MockField1: "test",
		MockField2: 1,
		MockField3: nil,
	}
	m := mockData.Mutator()

	testutils.Case(t, "returns correct types and values", func(t *testing.T) {
		field1, ok := m.GetField("1").(string)
		testutils.AssertTrue(t, ok)
		field2, ok := m.GetField("2").(int)
		testutils.AssertTrue(t, ok)
		field3, ok := m.GetField("3").(*float32)
		testutils.AssertTrue(t, ok)

		testutils.AssertEquals(t, "test", field1)
		testutils.AssertEquals(t, 1, field2)
		testutils.AssertEquals(t, nil, field3)
	})

	testutils.Case(t, "returns correct values after change", func(t *testing.T) {
		var newField3 float32 = 3.1
		mockData.MockField3 = &newField3

		field3, ok := m.GetField("3").(*float32)
		testutils.AssertTrue(t, ok)
		testutils.AssertEquals(t, &newField3, field3)
	})
}

func TestGetFields(t *testing.T) {
	mockData := &mocks.MockData{
		MockField1: "test",
		MockField2: 1,
		MockField3: nil,
	}
	m := mockData.Mutator()

	testutils.Case(t, "returns correct types and values", func(t *testing.T) {
		fields := m.GetFields()
		field1, ok := fields["1"].(string)
		testutils.AssertTrue(t, ok)
		field2, ok := fields["2"].(int)
		testutils.AssertTrue(t, ok)
		field3, ok := fields["3"].(*float32)
		testutils.AssertTrue(t, ok)

		testutils.AssertEquals(t, "test", field1)
		testutils.AssertEquals(t, 1, field2)
		testutils.AssertEquals(t, nil, field3)
	})

	testutils.Case(t, "returns correct values after change", func(t *testing.T) {
		var newField3 float32 = 3.1
		mockData.MockField3 = &newField3

		fields := m.GetFields()
		field3, ok := fields["3"].(*float32)
		testutils.AssertTrue(t, ok)
		testutils.AssertEquals(t, &newField3, field3)
	})
}

func TestSetField(t *testing.T) {
	mockData := &mocks.MockData{}
	m := mockData.Mutator()

	testutils.Case(t, "successfully sets fields", func(t *testing.T) {
		var field3 float32 = 3.0
		testutils.AssertOk(t, m.SetField("1", "test"))
		testutils.AssertOk(t, m.SetField("2", 1))
		testutils.AssertOk(t, m.SetField("3", &field3))

		expected := &mocks.MockData{
			MockField1: "test",
			MockField2: 1,
			MockField3: &field3,
		}
		mocks.AssertMockDataEquals(t, expected, mockData)
	})

	testutils.Case(t, "throws error for type mismatch", func(t *testing.T) {
		err := m.SetField("1", 1)
		testutils.AssertErrorEquals(t, mutator.SetFieldTypeError, err)
	})
}

func TestSetFields(t *testing.T) {
	mockData := &mocks.MockData{}
	m := mockData.Mutator()
	var field3 float32 = 3.0

	testutils.Case(t, "successfully sets all fields", func(t *testing.T) {
		fieldValues := mutator.MappedFieldValues{
			"1": "test",
			"2": 1,
			"3": &field3,
		}

		testutils.AssertOk(t, m.SetFields(fieldValues))

		expected := &mocks.MockData{
			MockField1: "test",
			MockField2: 1,
			MockField3: &field3,
		}
		mocks.AssertMockDataEquals(t, expected, mockData)
	})

	testutils.Case(t, "successfully sets some fields", func(t *testing.T) {
		fieldValues := mutator.MappedFieldValues{
			"1": "test2",
		}

		testutils.AssertOk(t, m.SetFields(fieldValues))

		expected := &mocks.MockData{
			MockField1: "test2",
			MockField2: 1,
			MockField3: &field3,
		}
		mocks.AssertMockDataEquals(t, expected, mockData)
	})

	testutils.Case(t, "throws error for type mismatch", func(t *testing.T) {
		fieldValues := mutator.MappedFieldValues{
			"1": 1,
		}
		err := m.SetFields(fieldValues)
		testutils.AssertErrorEquals(t, mutator.SetFieldTypeError, err)
	})
}
