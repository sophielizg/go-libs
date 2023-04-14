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
	subscriptionBackend, subscriptionName, err := t.Backend.Subscribe(subscriptionId)
	if err != nil {
		return nil, err
	}

	subscription := &Subscription[V, PV]{
		Id:          subscriptionId,
		ParentTopic: t,
		Queue: Queue[V, PV]{
			Backend: subscriptionBackend,
			Settings: NewTableSettings(
				WithTableName(subscriptionName),
				WithDataRowSettings(t.Settings.DataRowSettings),
			),
		},
	}
	subscription.Init()
	return subscription, nil
}

func (t *Topic[V, PV]) Unsubscribe(subscriptionId string) error {
	return t.Backend.Unsubscribe(subscriptionId)
}

type Subscription[V any, PV mutator.Mutatable[V]] struct {
	Queue[V, PV]
	Id          string
	ParentTopic *Topic[V, PV]
}

func (s *Subscription[V, PV]) Unsubscribe() error {
	return s.ParentTopic.Unsubscribe(s.Id)
}
