package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type Topic[V any, PV mutator.Mutatable[V]] struct {
	Backend        TopicBackendOps
	Settings       *TableSettings
	DataRowFactory mutator.MutatableFactory[V, PV]
}

func (t *Topic[V, PV]) Init() {
	t.Settings.ApplyOption(WithDataRow[V, PV]())
}

func (t *Topic[V, PV]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *Topic[V, PV]) SetBackend(backend TopicBackendOps) {
	t.Backend = backend
}

func (t *Topic[V, PV]) Publish(events ...PV) error {
	return t.Backend.Publish(t.DataRowFactory.CreateFieldValuesList(events))
}

func (t *Topic[V, PV]) Subscribe(subscriptionId string) (*Subscription[V, PV], error) {
	err := t.Backend.Subscribe(subscriptionId)
	if err != nil {
		return nil, err
	}

	return &Subscription[V, PV]{
		Id:          subscriptionId,
		ParentTopic: t,
	}, nil
}

func (t *Topic[V, PV]) Unsubscribe(subscriptionId string) error {
	return t.Backend.Unsubscribe(subscriptionId)
}

func (t *Topic[V, PV]) HasMessage(subscriptionId string) (bool, error) {
	return t.Backend.HasMessage(subscriptionId)
}

func (t *Topic[V, PV]) RecieveMessage(subscriptionId string) (string, PV, error) {
	messageId, messageFields, err := t.Backend.RecieveMessage(subscriptionId)
	if err != nil {
		return "", nil, err
	}

	message, err := t.DataRowFactory.CreateFromFields(messageFields)
	return messageId, message, err
}

func (t *Topic[V, PV]) AckSuccess(subscriptionId string, messageId string) error {
	return t.Backend.AckSuccess(subscriptionId, messageId)
}

func (t *Topic[V, PV]) AckFailure(subscriptionId string, messageId string) error {
	return t.Backend.AckFailure(subscriptionId, messageId)
}

type Subscription[V any, PV mutator.Mutatable[V]] struct {
	Id          string
	ParentTopic *Topic[V, PV]
}

func (s *Subscription[V, PV]) Unsubscribe() error {
	return s.ParentTopic.Unsubscribe(s.Id)
}

func (s *Subscription[V, PV]) HasMessage() (bool, error) {
	return s.ParentTopic.HasMessage(s.Id)
}

func (s *Subscription[V, PV]) RecieveMessage() (string, PV, error) {
	return s.ParentTopic.RecieveMessage(s.Id)
}

func (s *Subscription[V, PV]) AckSuccess(messageId string) error {
	return s.ParentTopic.AckSuccess(s.Id, messageId)
}

func (s *Subscription[V, PV]) AckFailure(messageId string) error {
	return s.ParentTopic.AckFailure(s.Id, messageId)
}
