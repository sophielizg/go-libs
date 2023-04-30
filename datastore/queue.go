package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type Queue[M any, PM mutator.Mutatable[M]] struct {
	backend        QueueBackendQueries
	Settings       *TableSettings
	MessageFactory mutator.MutatableFactory[M, PM]
	*queries.Countable
	*queries.MessageReceiveable[M, PM]
}

func (q *Queue[M, PM]) Init() {
	q.Settings.ApplyOption(WithEntry[M, PM]())
	q.Countable = &queries.Countable{}
	q.MessageReceiveable = &queries.MessageReceiveable[M, PM]{}
}

func (q *Queue[M, PM]) GetSettings() *TableSettings {
	return q.Settings
}

func (q *Queue[M, PM]) SetBackend(backend QueueBackendQueries) {
	q.backend = backend
	q.Countable.SetBackend(backend)
	q.MessageReceiveable.SetBackend(backend)
}

func (q *Queue[M, PM]) SendMessage(messages ...PM) error {
	return q.backend.SendMessage(q.MessageFactory.CreateFieldValuesList(messages))
}

func (q *Queue[M, PM]) TransferTo(newQueue *Queue[M, PM], batchSize int) error {
	bufMessages := make([]PM, 0, batchSize)
	bufIds := make([]string, 0, batchSize)
	for {
		size, err := q.Count()
		if err != nil {
			return err
		} else if size == 0 {
			break
		}

		id, message, err := q.RecieveMessage()
		if err != nil {
			return err
		}

		bufMessages = append(bufMessages, message)
		bufIds = append(bufIds, id)
		if len(bufMessages) == batchSize {
			if err = newQueue.SendMessage(bufMessages...); err != nil {
				q.AckFailure(bufIds...)
				return err
			}

			q.AckSuccess(bufIds...)
			bufMessages = make([]PM, 0, batchSize)
			bufIds = make([]string, 0, batchSize)
		}
	}

	if err := newQueue.SendMessage(bufMessages...); err != nil {
		q.AckFailure(bufIds...)
		return err
	}

	q.AckSuccess(bufIds...)
	return nil
}
