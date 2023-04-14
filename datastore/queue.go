package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type Queue[V any, PV mutator.Mutatable[V]] struct {
	Backend        QueueBackendOps
	Settings       *TableSettings
	DataRowFactory mutator.MutatableFactory[V, PV]
}

func (q *Queue[V, PV]) Init() {
	q.Settings.ApplyOption(WithDataRow[V, PV]())
}

func (q *Queue[V, PV]) GetSettings() *TableSettings {
	return q.Settings
}

func (q *Queue[V, PV]) SetBackend(backend QueueBackendOps) {
	q.Backend = backend
}

func (q *Queue[V, PV]) Size() (int, error) {
	return q.Backend.Size()
}

func (q *Queue[V, PV]) Push(messages ...PV) error {
	return q.Backend.Push(q.DataRowFactory.CreateFieldValuesList(messages))
}

func (q *Queue[V, PV]) Pop() (string, PV, error) {
	messageId, messageFields, err := q.Backend.Pop()
	if err != nil {
		return "", nil, err
	}

	message, err := q.DataRowFactory.CreateFromFields(messageFields)
	return messageId, message, err
}

func (q *Queue[V, PV]) AckSuccess(messageId string) error {
	return q.Backend.AckSuccess(messageId)
}

func (q *Queue[V, PV]) AckFailure(messageId string) error {
	return q.Backend.AckFailure(messageId)
}

func (q *Queue[V, PV]) TransferTo(newQueue *Queue[V, PV], batchSize int) error {
	buf := make([]PV, 0, batchSize)
	for {
		size, err := q.Size()
		if err != nil {
			return err
		} else if size == 0 {
			break
		}

		_, message, err := q.Pop()
		if err != nil {
			return err
		}

		buf = append(buf, message)
		if len(buf) == batchSize {
			if err = newQueue.Push(buf...); err != nil {
				return err
			}
			buf = make([]PV, 0, batchSize)
		}
	}

	return newQueue.Push(buf...)
}
