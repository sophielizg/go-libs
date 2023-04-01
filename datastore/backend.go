package datastore

// TODO: revisit some of this and other naming (and filenames)
type Backend[C Connection] struct {
	Conn C
}

func (b *Backend[C]) RegisterTables(registerFuncs ...func(b *Backend[C]) error) error {
	for _, register := range registerFuncs {
		if err := register(b); err != nil {
			return err
		}
	}

	return nil
}

func NewBackend[C Connection](options ...func(*Backend[C])) *Backend[C] {
	backend := &Backend[C]{}

	for _, option := range options {
		option(backend)
	}

	return backend
}

func WithConnection[C Connection](conn C) func(*Backend[C]) {
	return func(b *Backend[C]) {
		b.Conn = conn
	}
}

func RegisterAppendTable[C Connection, TB AppendTableBackend[C], T Table[AppendTableBackendOps]](table T, tableBackend TB) func(*Backend[C]) error {
	return func(b *Backend[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(b.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterHashTable[C Connection, TB HashTableBackend[C], T Table[HashTableBackendOps]](table T, tableBackend TB) func(*Backend[C]) error {
	return func(b *Backend[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(b.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}

func RegisterSortTable[C Connection, TB SortTableBackend[C], T Table[SortTableBackendOps]](table T, tableBackend TB) func(*Backend[C]) error {
	return func(b *Backend[C]) error {
		table.SetBackend(tableBackend)
		table.Init()
		tableBackend.SetConnection(b.Conn)
		tableBackend.SetSettings(table.GetSettings())
		return tableBackend.Register()
	}
}
