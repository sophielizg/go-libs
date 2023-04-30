package inmemory

import (
	"container/list"
	"errors"
	"strconv"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type InFlightMessages = map[string]mutator.MappedFieldValues

type QueueItem struct {
	id      int
	message mutator.MappedFieldValues
}

type Queue struct {
	messageQueue     *list.List
	inFlightMessages InFlightMessages
}

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
	if err := validateAutoGenerateSettings(b.settings.DataSettings); err != nil {
		return err
	}

	if queue := b.conn.GetQueue(b.settings); queue == nil {
		queue = &Queue{
			messageQueue:     &list.List{},
			inFlightMessages: InFlightMessages{},
		}
	}

	return nil
}

func (b *QueueBackend) Drop() error {
	b.conn.DropQueue(b.settings)
	return nil
}

func (b *QueueBackend) Count() (int, error) {
	queue := b.conn.GetQueue(b.settings)
	return queue.messageQueue.Len(), nil
}

func (b *QueueBackend) HasMessage() (bool, error) {
	count, err := b.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *QueueBackend) SendMessage(messages []mutator.MappedFieldValues) error {
	queue := b.conn.GetQueue(b.settings)
	back := queue.messageQueue.Back()

	lastItem, ok := back.Value.(QueueItem)
	if !ok {
		// this should never happen
		return errors.New("send message failed, item with invalid type in queue")
	}

	lastId := lastItem.id
	for _, message := range messages {
		lastId += 1
		queue.messageQueue.PushBack(QueueItem{
			id:      lastId,
			message: message,
		})
	}

	return nil
}

func (b *QueueBackend) RecieveMessage() (string, mutator.MappedFieldValues, error) {
	queue := b.conn.GetQueue(b.settings)
	popped := queue.messageQueue.Front()
	queue.messageQueue.Remove(popped)

	item, ok := popped.Value.(QueueItem)
	if !ok {
		// this should never happen
		return "", nil, errors.New("recieve message failed, item with invalid type in queue")
	}

	idStr := strconv.Itoa(item.id)
	queue.inFlightMessages[idStr] = item.message
	return idStr, item.message, nil
}

func (b *QueueBackend) AckSuccess(messageIds []string) error {
	queue := b.conn.GetQueue(b.settings)

	for _, messageId := range messageIds {
		if queue.inFlightMessages[messageId] == nil {
			return KeyDoesNotExistError
		}

		delete(queue.inFlightMessages, messageId)
	}

	return nil
}

func (b *QueueBackend) AckFailure(messageIds []string) error {
	queue := b.conn.GetQueue(b.settings)

	for _, messageId := range messageIds {
		if queue.inFlightMessages[messageId] == nil {
			return KeyDoesNotExistError
		}

		messageIdInt, err := strconv.Atoi(messageId)
		if err != nil {
			return err
		}

		queue.messageQueue.PushFront(QueueItem{
			id:      messageIdInt,
			message: queue.inFlightMessages[messageId],
		})
	}

	return nil
}
