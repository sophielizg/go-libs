package inmemory

import (
	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type AppendTable = []mutator.MappedFieldValues

type AppendTableBackend struct {
	conn     *Connection
	settings *datastore.TableSettings
}

func (b *AppendTableBackend) SetSettings(settings *datastore.TableSettings) {
	b.settings = settings
}

func (b *AppendTableBackend) SetConnection(conn *Connection) {
	b.conn = conn
}

func (b *AppendTableBackend) Register() error {
	if err := validateAutoGenerateSettings(b.settings.DataSettings); err != nil {
		return err
	}

	if table := b.conn.GetAppendTable(b.settings); table == nil {
		b.conn.SetAppendTable(b.settings, AppendTable{})
	}

	return nil
}

func (b *AppendTableBackend) Drop() error {
	b.conn.DropAppendTable(b.settings)
	return nil
}

func (b *AppendTableBackend) Scan(batchSize int) (chan mutator.MappedFieldValues, chan error) {
	outChan := make(chan mutator.MappedFieldValues, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		for _, fields := range b.conn.GetAppendTable(b.settings) {
			outChan <- fields
		}
	}()

	return outChan, errorChan
}

func (b *AppendTableBackend) Add(entries []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error) {
	table := b.conn.GetAppendTable(b.settings)
	table = append(table, entries...)
	b.conn.SetAppendTable(b.settings, table)
	return entries, nil
}
