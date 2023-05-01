package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type Transferable[E any, PE mutator.Mutatable[E]] struct {
	Scanable *Scanable[E, PE]
	Addable  *Addable[E, PE]
}

func (t *Transferable[E, PE]) TransferTo(newTable *Transferable[E, PE], batchSize int) error {
	dataChan, errorChan := t.Scanable.Scan(batchSize)

	for {
		buf := make([]PE, 0, batchSize)

		for {
			select {
			case err, more := <-errorChan:
				if !more {
					errorChan = nil
					break
				}

				return err
			case data, more := <-dataChan:
				if !more {
					dataChan = nil
					break
				}

				buf = append(buf, data)
			}

			if len(buf) == batchSize || (dataChan == nil && errorChan == nil) {
				break
			}
		}

		if len(buf) > 0 {
			_, err := newTable.Addable.Add(buf...)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
