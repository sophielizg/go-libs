package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type ScanableBackend interface {
	Scan(batchSize int) (chan mutator.MappedFieldValues, chan error)
}

type Scanable[E any, PE mutator.Mutatable[E]] struct {
	backend      ScanableBackend
	entryFactory mutator.MutatableFactory[E, PE]
}

func (s *Scanable[E, PE]) SetBackend(tableBackend ScanableBackend) {
	s.backend = tableBackend
}

func (s *Scanable[E, PE]) Scan(batchSize int) (chan PE, chan error) {
	inChan, inErrorChan := s.backend.Scan(batchSize)

	outChan := make(chan PE, 1)
	outErrorChan := make(chan error, 1)
	go func() {
		defer close(outChan)
		defer close(outErrorChan)

		for {
			select {
			case err, more := <-inErrorChan:
				if !more {
					inErrorChan = nil
					break
				}

				outErrorChan <- err

			case inFields, more := <-inChan:
				if !more {
					inChan = nil
					break
				}

				entry, err := s.entryFactory.CreateFromFields(inFields)
				if err != nil {
					outErrorChan <- err
				} else {
					outChan <- entry
				}
			}

			if inErrorChan == nil && inChan == nil {
				break
			}
		}
	}()

	return outChan, outErrorChan
}
