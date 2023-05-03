package datastore_test

// TODO adapt old test for TransferTo

// func TestQueueTransferTo(t *testing.T) {
// 	type transferInputVals struct {
// 		Error   error
// 		Message mutator.MappedFieldValues
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
// 					Message: mutator.MappedFieldValues{
// 						shipping.NameKey: "test1",
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
// 					Message: mutator.MappedFieldValues{
// 						shipping.NameKey: 1,
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
// 					Message: mutator.MappedFieldValues{
// 						shipping.NameKey: "test",
// 					},
// 				},
// 				Expected: &transferExpectedVals{
// 					Error: mockError,
// 				},
// 			},
// 		},
// 		Func: func(input *transferInputVals, expected *transferExpectedVals) {
// 			mockBackendSrc := &MockQueueBackendOps{
// 				ErrorRval:   input.Error,
// 				MessageRval: input.Message,
// 				SizeRval:    1,
// 			}
// 			srcQueue := shipping.NewPendingShipmentQueue()
// 			srcQueue.SetBackend(mockBackendSrc)

// 			mockBackendDest := &MockQueueBackendOps{}
// 			destQueue := shipping.NewPendingShipmentQueue()
// 			destQueue.SetBackend(mockBackendDest)

// 			err := srcQueue.TransferTo(destQueue, 10)
// 			testutils.AssertErrorEquals(t, expected.Error, err)

// 			if expected.Error != nil {
// 				return
// 			}

// 			testutils.AssertEquals(t, 1, len(mockBackendDest.MessagesInput))
// 			AssertShippingMessageFieldsEqualOrDefault(t, input.Message, mockBackendDest.MessagesInput[0])
// 		},
// 	}

// 	tests.Run(t)
// }
