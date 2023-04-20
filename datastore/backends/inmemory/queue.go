package inmemory

import (
	"container/list"
	"errors"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type Queue = list.List

type QueueBackend struct {
	conn     *Connection
	settings *datastore.TableSettings
}

func (b *QueueBackend) SetSettings(settings *datastore.TableSettings) {
	b.settings = settings
}

func (b *QueueBackend) SetConnection(conn *Connection) {
	b.conn = conn
}

func (b *QueueBackend) Register() error {
	if err := validateAutoGenerateSettings(b.settings.DataRowSettings); err != nil {
		return err
	}

	if queue := b.conn.GetQueue(b.settings); queue == nil {
		queue = &Queue{}
	}

	return nil
}

func (b *QueueBackend) Drop() error {
	b.conn.DropQueue(b.settings)
	return nil
}

func (b *QueueBackend) Size() (int, error) {
	queue := b.conn.GetQueue(b.settings)
	return queue.Len(), nil
}

func (b *QueueBackend) Push(messages []mutator.MappedFieldValues) error {
	queue := b.conn.GetQueue(b.settings)

	for _, message := range messages {
		queue.PushBack(message)
	}

	return nil
}

func (b *QueueBackend) Pop() (string, mutator.MappedFieldValues, error) {
	queue := b.conn.GetQueue(b.settings)
	popped := queue.Front()
	queue.Remove(popped)

	message, ok := popped.Value.(mutator.MappedFieldValues)
	if !ok {
		// this should never happen
		return "", nil, errors.New("pop failed, item with invalid type in queue")
	}

	return "", message, nil
}

func (b *QueueBackend) AckSuccess(messageId string) error {
	// do nothing for ack
	return nil
}

func (b *QueueBackend) AckFailure(messageId string) error {
	// do nothing for ack
	return nil
}
