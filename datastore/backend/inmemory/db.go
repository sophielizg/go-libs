package inmemory

import (
	"errors"

	"github.com/sophielizg/go-libs/datastore"
)

type Table = map[string][]datastore.DataRowFields

type InsertOrder = []string

type InMemoryDatastoreConnection struct {
	db            map[string]Table
	dbInsertOrder map[string]InsertOrder
}

func (c *InMemoryDatastoreConnection) CreateOrUpdateSchema(tableName string) error {
	if c.db == nil {
		c.db = map[string]Table{}
	}

	if c.dbInsertOrder == nil {
		c.dbInsertOrder = map[string]InsertOrder{}
	}

	table, ok := c.db[tableName]
	if !ok || table == nil {
		c.db[tableName] = Table{}
		c.dbInsertOrder[tableName] = InsertOrder{}
	}

	return nil
}

func (c *InMemoryDatastoreConnection) Scan(tableName string, batchSize int) (chan *datastore.HashTableScanFields, chan error) {
	outChan := make(chan *datastore.HashTableScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		table := c.db[tableName]
		insertOrder := c.dbInsertOrder[tableName]
		if table == nil || insertOrder == nil {
			errorChan <- errors.New("No table exists with given schema name")
			return
		}

		for _, key := range insertOrder {
			hashKey, err := unstringifyHashKey(key)
			if err != nil {
				errorChan <- err
				continue
			}

			for _, dataRowFields := range table[key] {
				outChan <- &datastore.HashTableScanFields{
					AppendTableScanFields: datastore.AppendTableScanFields{
						DataRow: dataRowFields,
					},
					HashKey: hashKey,
				}
			}
		}
	}()

	return outChan, errorChan
}

func (c *InMemoryDatastoreConnection) Get(tableName string, hashKey datastore.DataRowFields) ([]datastore.DataRowFields, error) {
	key, err := stringifyHashKey(hashKey)
	if err != nil {
		return nil, err
	}

	table := c.db[tableName]
	if table == nil {
		return nil, errors.New("No table exists with given schema name")
	}

	return table[key], nil
}

func (c *InMemoryDatastoreConnection) Add(tableName string, hashKey datastore.DataRowFields, dataRow datastore.DataRowFields) error {
	key, err := stringifyHashKey(hashKey)
	if err != nil {
		return err
	}

	table := c.db[tableName]
	insertOrder := c.dbInsertOrder[tableName]
	if table == nil || insertOrder == nil {
		return errors.New("No table exists with given schema name")
	}

	table[key] = []datastore.DataRowFields{dataRow}
	c.dbInsertOrder[tableName] = append(insertOrder, key)
	return nil
}

func (c *InMemoryDatastoreConnection) Update(tableName string, hashKey datastore.DataRowFields, dataRow datastore.DataRowFields) error {
	key, err := stringifyHashKey(hashKey)
	if err != nil {
		return err
	}

	table := c.db[tableName]
	if table == nil {
		return errors.New("No table exists with given schema name")
	}

	if _, ok := table[key]; ok {
		table[key] = []datastore.DataRowFields{dataRow}
	}

	return nil
}

func (c *InMemoryDatastoreConnection) Delete(tableName string, hashKey datastore.DataRowFields) error {
	key, err := stringifyHashKey(hashKey)
	if err != nil {
		return err
	}

	table := c.db[tableName]
	if table == nil {
		return errors.New("No table exists with given schema name")
	}

	if _, ok := table[key]; ok {
		delete(table, key)
	}

	return nil
}
