package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
)

func scan[M any, PM mutator.Mutatable[M]](inChan chan mutator.MappedFieldValues, inErrorChan chan error, entryFactory *mutator.MutatableFactory[M, PM]) (chan PM, chan error) {
	outChan := make(chan PM, 1)
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

				entry, err := entryFactory.CreateFromFields(inFields)
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
