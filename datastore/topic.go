package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type Topic[M any, PM mutator.Mutatable[M]] struct {
	backend        TopicBackendQueries
	Settings       *TableSettings
	MessageFactory mutator.MutatableFactory[M, PM]
}

func (t *Topic[M, PM]) Init() {
	t.Settings.ApplyOption(WithEntry[M, PM]())
}

func (t *Topic[M, PM]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *Topic[M, PM]) SetBackend(backend TopicBackendQueries) {
	t.backend = backend
}

func (t *Topic[M, PM]) Publish(messages ...PM) error {
	return t.backend.Publish(t.MessageFactory.CreateFieldValuesList(messages))
}

func (t *Topic[M, PM]) Subscribe(subscriptionId string) (*Subscription[M, PM], error) {
	subscriptionBackend, err := t.backend.Subscribe(subscriptionId)
	if err != nil {
		return nil, err
	}

	subscription := &Subscription[M, PM]{
		Id:          subscriptionId,
		ParentTopic: t,
	}
	subscription.Init()
	subscription.SetBackend(subscriptionBackend)
	return subscription, nil
}

type Subscription[M any, PM mutator.Mutatable[M]] struct {
	backend     SubscriptionBackendQueries
	Id          string
	ParentTopic *Topic[M, PM]
	*queries.MessageReceiveable[M, PM]
}

func (s *Subscription[M, PM]) Init() {
	s.MessageReceiveable = &queries.MessageReceiveable[M, PM]{}
}

func (s *Subscription[M, PM]) SetBackend(backend SubscriptionBackendQueries) {
	s.backend = backend
	s.MessageReceiveable.SetBackend(backend)
}

func (s *Subscription[M, PM]) Unsubscribe() error {
	return s.backend.Unsubscribe()
}
