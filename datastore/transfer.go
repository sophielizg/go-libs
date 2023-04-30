package datastore

type Transferable[E any] interface {
	Add(entries ...E) ([]E, error)
}

func transfer[E any](batchSize int, dataChan chan E, errorChan chan error, destTable Transferable[E]) error {
	for {
		buf := make([]E, 0, batchSize)

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
			_, err := destTable.Add(buf...)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
