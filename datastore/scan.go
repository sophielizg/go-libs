package datastore

import "github.com/sophielizg/go-libs/datastore/fields"

type ScanFields struct {
	DataRow fields.MappedFieldValues
	HashKey fields.MappedFieldValues
	SortKey fields.MappedFieldValues
}

func scan[O any](inChan chan *ScanFields, inErrorChan chan error, convertFieldsToOutput func(*ScanFields) (O, error)) (chan O, chan error) {
	outChan := make(chan O, 1)
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

				converted, err := convertFieldsToOutput(inFields)

				if err != nil {
					outErrorChan <- err
				} else {
					outChan <- converted
				}
			}

			if inErrorChan == nil && inChan == nil {
				break
			}
		}
	}()

	return outChan, outErrorChan
}
