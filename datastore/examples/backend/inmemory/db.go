package inmemory

import (
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/sophielizg/go-libs/datastore"
)

type Table = map[string][]datastore.DataRowFields

type InsertOrder = []string

var (
	db            = map[string]Table{}
	dbInsertOrder = map[string]InsertOrder{}
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func createOrUpdateDbSchema(tableName string) error {
	table, ok := db[tableName]
	if !ok || table == nil {
		db[tableName] = Table{}
		dbInsertOrder[tableName] = InsertOrder{}
	}

	return nil
}

func scanDb(tableName string, batchSize int) (chan *datastore.DataRowScanFields, chan error) {
	outChan := make(chan *datastore.DataRowScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		table := db[tableName]
		insertOrder := dbInsertOrder[tableName]
		if table == nil || insertOrder == nil {
			errorChan <- errors.New("No table exists with given schema name")
			return
		}

		for _, key := range insertOrder {
			for _, dataRowFields := range table[key] {
				outChan <- &datastore.DataRowScanFields{
					DataRow: dataRowFields,
				}
			}
		}
	}()

	return outChan, errorChan
}

func generateStringKey(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func stringifyHashKey(hashKey datastore.DataRowFields) (string, error) {
	bytes, err := json.Marshal(hashKey)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func insertWithHashKey(tableName string, hashKey datastore.DataRowFields, dataRow datastore.DataRowFields) error {
	key, err := stringifyHashKey(hashKey)
	if err != nil {
		return err
	}

	table := db[tableName]
	insertOrder := dbInsertOrder[tableName]
	if table == nil || insertOrder == nil {
		return errors.New("No table exists with given schema name")
	}

	table[key] = []datastore.DataRowFields{dataRow}
	dbInsertOrder[tableName] = append(insertOrder, key)
	return nil
}
