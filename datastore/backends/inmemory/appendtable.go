package inmemory

import (
	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/fields"
)

type AppendTable = []fields.MappedFieldValues

type Connection struct {
	table map[string]*AppendTable
}

func (c *Connection) Close() {}

func (c *Connection) GetTable(settings *datastore.TableSettings) *AppendTable {
	return c.table[settings.Name]
}

func NewConnection() *Connection {
	return &Connection{
		table: map[string]*AppendTable{},
	}
}

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
	if table := b.conn.GetTable(b.settings); table == nil {
		table = &AppendTable{}
	}

	return nil
}

func (b *AppendTableBackend) Drop() error {
	b.conn.table[b.settings.Name] = nil
	return nil
}

func (b *AppendTableBackend) Scan(batchSize int) (chan *datastore.ScanFields, chan error) {
	outChan := make(chan *datastore.ScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		for _, fields := range *b.conn.GetTable(b.settings) {
			outChan <- &datastore.ScanFields{DataRow: fields}
		}
	}()

	return outChan, errorChan
}

func (b *AppendTableBackend) AddMultiple(data []fields.MappedFieldValues) error {
	table := b.conn.GetTable(b.settings)
	*table = append(*table, data...)
	return nil
}
