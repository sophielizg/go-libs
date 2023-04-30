package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type MessageReceiveableBackend interface {
	HasMessage() (bool, error)
	RecieveMessage() (string, mutator.MappedFieldValues, error)
	AckSuccess(messageId []string) error
	AckFailure(messageId []string) error
}

type MessageReceiveable[M any, PM mutator.Mutatable[M]] struct {
	backend        MessageReceiveableBackend
	messageFactory mutator.MutatableFactory[M, PM]
}

func (s *MessageReceiveable[M, PM]) SetBackend(backend MessageReceiveableBackend) {
	s.backend = backend
}

func (m *MessageReceiveable[M, PM]) HasMessage() (bool, error) {
	return m.backend.HasMessage()
}

func (m *MessageReceiveable[M, PM]) RecieveMessage() (string, PM, error) {
	messageId, messageFields, err := m.backend.RecieveMessage()
	if err != nil {
		return "", nil, err
	}

	message, err := m.messageFactory.CreateFromFields(messageFields)
	return messageId, message, err
}

func (m *MessageReceiveable[M, PM]) AckSuccess(messageId ...string) error {
	return m.backend.AckSuccess(messageId)
}

func (m *MessageReceiveable[M, PM]) AckFailure(messageId ...string) error {
	return m.backend.AckFailure(messageId)
}
