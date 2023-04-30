package queries

type CountableBackend interface {
	Count() (int, error)
}

type Countable struct {
	backend CountableBackend
}

func (a *Countable) SetBackend(tableBackend CountableBackend) {
	a.backend = tableBackend
}

func (a *Countable) Count() (int, error) {
	return a.backend.Count()
}
