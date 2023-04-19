package datastore

type Connection interface {
	Close()
}

type ConnectionGroup[C Connection] struct {
	Conn C
}

func (g *ConnectionGroup[C]) RegisterTables(registerFuncs ...func(g *ConnectionGroup[C]) error) error {
	for _, register := range registerFuncs {
		if err := register(g); err != nil {
			return err
		}
	}

	return nil
}

func NewConnectionGroup[C Connection](options ...func(*ConnectionGroup[C])) *ConnectionGroup[C] {
	group := &ConnectionGroup[C]{}

	for _, option := range options {
		option(group)
	}

	return group
}

func WithConnection[C Connection](conn C) func(*ConnectionGroup[C]) {
	return func(g *ConnectionGroup[C]) {
		g.Conn = conn
	}
}

type Table[B any] interface {
	Init()
	GetSettings() *TableSettings
	SetBackend(tableBackend B)
}

func RegisterAppendTable[C Connection, TB AppendTableBackend[C], T Table[AppendTableBackendOps]](table T, tableBackend TB) func(*ConnectionGroup[C]) error {
	return func(g *ConnectionGroup[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(g.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterHashTable[C Connection, TB HashTableBackend[C], T Table[HashTableBackendOps]](table T, tableBackend TB) func(*ConnectionGroup[C]) error {
	return func(g *ConnectionGroup[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(g.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterSortTable[C Connection, TB SortTableBackend[C], T Table[SortTableBackendOps]](table T, tableBackend TB) func(*ConnectionGroup[C]) error {
	return func(g *ConnectionGroup[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(g.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterQueue[C Connection, TB QueueBackend[C], T Table[QueueBackendOps]](table T, tableBackend TB) func(*ConnectionGroup[C]) error {
	return func(g *ConnectionGroup[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(g.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterTopic[C Connection, TB TopicBackend[C], T Table[TopicBackendOps]](table T, tableBackend TB) func(*ConnectionGroup[C]) error {
	return func(g *ConnectionGroup[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(g.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}
