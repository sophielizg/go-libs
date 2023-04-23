package inmemory

import "github.com/sophielizg/go-libs/datastore"

type Connection struct {
	appendTables map[string]*AppendTable
	hashTables   map[string]HashTable
	queues       map[string]*Queue
}

func (c *Connection) Close() {}

func (c *Connection) GetAppendTable(settings *datastore.TableSettings) *AppendTable {
	return c.appendTables[settings.Name]
}

func (c *Connection) DropAppendTable(settings *datastore.TableSettings) {
	c.appendTables[settings.Name] = nil
}

func (c *Connection) GetHashTable(settings *datastore.TableSettings) HashTable {
	return c.hashTables[settings.Name]
}

func (c *Connection) SetHashTable(settings *datastore.TableSettings, newTable HashTable) {
	c.hashTables[settings.Name] = newTable
}

func (c *Connection) DropHashTable(settings *datastore.TableSettings) {
	c.hashTables[settings.Name] = nil
}

func (c *Connection) GetQueue(settings *datastore.TableSettings) *Queue {
	return c.queues[settings.Name]
}

func (c *Connection) DropQueue(settings *datastore.TableSettings) {
	c.queues[settings.Name] = nil
}

func NewConnection() *Connection {
	return &Connection{
		appendTables: map[string]*AppendTable{},
		hashTables:   map[string]HashTable{},
		queues:       map[string]*Queue{},
	}
}
