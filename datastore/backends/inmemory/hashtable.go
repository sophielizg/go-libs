package inmemory

import (
	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type HashTable = map[string]mutator.MappedFieldValues

type HashTableBackend struct {
	conn     *Connection
	settings *datastore.TableSettings
}

func (b *HashTableBackend) SetSettings(settings *datastore.TableSettings) {
	b.settings = settings
}

func (b *HashTableBackend) SetConnection(conn *Connection) {
	b.conn = conn
}

func (b *HashTableBackend) Register() error {
	if err := validateAutoGenerateSettings(b.settings.DataSettings); err != nil {
		return err
	} else if err := validateAutoGenerateSettings(b.settings.KeySettings); err != nil {
		return err
	}

	if b.conn.GetHashTable(b.settings) == nil {
		b.conn.SetHashTable(b.settings, HashTable{})
	}

	return nil
}

func (b *HashTableBackend) Drop() error {
	b.conn.DropHashTable(b.settings)
	return nil
}

func (b *HashTableBackend) Count() (int, error) {
	return len(b.conn.GetHashTable(b.settings)), nil
}

func (b *HashTableBackend) Scan(batchSize int) (chan mutator.MappedFieldValues, chan error) {
	outChan := make(chan mutator.MappedFieldValues, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		for _, entry := range b.conn.GetHashTable(b.settings) {
			outChan <- entry
		}
	}()

	return outChan, errorChan
}

func (b *HashTableBackend) Get(keys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	table := b.conn.GetHashTable(b.settings)

	data := make([]mutator.MappedFieldValues, len(keys))
	for i, key := range keys {
		keyStr, err := stringifyKey(key)
		if err != nil {
			return nil, err
		}

		data[i] = table[keyStr]
	}

	return data, nil
}

func (b *HashTableBackend) Add(entries []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	table := b.conn.GetHashTable(b.settings)

	for _, entry := range entries {
		key := getKeyFromEntry(b.settings, entry)
		keyStr, err := stringifyKey(key)
		if err != nil {
			return nil, err
		} else if table[keyStr] != nil {
			return nil, KeyExistsError
		}

		table[keyStr] = entry
	}

	b.conn.SetHashTable(b.settings, table)
	return entries, nil
}

func (b *HashTableBackend) Update(entries []mutator.MappedFieldValues) error {
	table := b.conn.GetHashTable(b.settings)

	for _, entry := range entries {
		key := getKeyFromEntry(b.settings, entry)
		keyStr, err := stringifyKey(key)
		if err != nil {
			return err
		} else if table[keyStr] == nil {
			return KeyDoesNotExistError
		}

		table[keyStr] = entry
	}

	b.conn.SetHashTable(b.settings, table)
	return nil
}

func (b *HashTableBackend) Delete(keys []mutator.MappedFieldValues) error {
	table := b.conn.GetHashTable(b.settings)

	for _, key := range keys {
		keyStr, err := stringifyKey(key)
		if err != nil {
			return err
		} else if table[keyStr] == nil {
			return KeyDoesNotExistError
		}

		delete(table, keyStr)
	}

	b.conn.SetHashTable(b.settings, table)
	return nil
}
