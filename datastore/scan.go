package datastore

func scan[I any, O any](inChan chan I, inErrorChan chan error, convertFunc func(I) (O, error)) (chan O, chan error) {
	outChan := make(chan O)
	outErrorChan := make(chan error)
	go func() {
		select {
		case err := <-inErrorChan:
			outErrorChan <- err
		case inFields, more := <-inChan:
			converted, err := convertFunc(inFields)

			if err != nil {
				outErrorChan <- err
			} else {
				outChan <- converted
			}

			if !more {
				close(outChan)
			}
		}
		close(outErrorChan)
	}()

	return outChan, outErrorChan
}
