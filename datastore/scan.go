package datastore

func scan[I any, O any](batchSize int, inChan chan I, inErrorChan chan error, convertFunc func(I) (O, error)) (chan O, chan error) {
	outChan := make(chan O, 1)
	outErrorChan := make(chan error, 1)
	go func() {
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
					close(outChan)
					break
				}

				converted, err := convertFunc(inFields)

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
		close(outErrorChan)
	}()

	return outChan, outErrorChan
}
