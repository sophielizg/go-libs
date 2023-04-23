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
	if err := validateAutoGenerateSettings(b.settings.DataRowSettings); err != nil {
		return err
	} else if err := validateAutoGenerateSettings(b.settings.HashKeySettings); err != nil {
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

func (b *HashTableBackend) Scan(batchSize int) (chan *datastore.ScanFields, chan error) {
	outChan := make(chan *datastore.ScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		for haskKeyStr, dataRow := range b.conn.GetHashTable(b.settings) {
			hashKey, err := unstringifyKey(haskKeyStr)
			if err != nil {
				errorChan <- err
			} else {
				outChan <- &datastore.ScanFields{
					DataRow: dataRow,
					HashKey: hashKey,
				}
			}
		}
	}()

	return outChan, errorChan
}

func (b *HashTableBackend) GetMultiple(hashKeys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	table := b.conn.GetHashTable(b.settings)

	dataRows := make([]mutator.MappedFieldValues, len(hashKeys))
	for i, hashKey := range hashKeys {
		hashKeyStr, err := stringifyKey(hashKey)
		if err != nil {
			return nil, err
		}

		dataRows[i] = table[hashKeyStr]
	}

	return dataRows, nil
}

func (b *HashTableBackend) AddMultiple(hashKeys []mutator.MappedFieldValues, dataRows []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	table := b.conn.GetHashTable(b.settings)

	for i := range hashKeys {
		hashKeyStr, err := stringifyKey(hashKeys[i])
		if err != nil {
			return nil, err
		} else if table[hashKeyStr] != nil {
			return nil, HashKeyExistsError
		}

		table[hashKeyStr] = dataRows[i]
	}

	b.conn.SetHashTable(b.settings, table)
	return hashKeys, nil
}

func (b *HashTableBackend) UpdateMultiple(hashKeys []mutator.MappedFieldValues, dataRows []mutator.MappedFieldValues) error {
	table := b.conn.GetHashTable(b.settings)

	for i := range hashKeys {
		hashKeyStr, err := stringifyKey(hashKeys[i])
		if err != nil {
			return err
		} else if table[hashKeyStr] == nil {
			return HashKeyDoesNotExistError
		}

		table[hashKeyStr] = dataRows[i]
	}

	b.conn.SetHashTable(b.settings, table)
	return nil
}

func (b *HashTableBackend) DeleteMultiple(hashKeys []mutator.MappedFieldValues) error {
	table := b.conn.GetHashTable(b.settings)

	for i := range hashKeys {
		hashKeyStr, err := stringifyKey(hashKeys[i])
		if err != nil {
			return err
		} else if table[hashKeyStr] == nil {
			return HashKeyDoesNotExistError
		}

		delete(table, hashKeyStr)
	}

	b.conn.SetHashTable(b.settings, table)
	return nil
}
