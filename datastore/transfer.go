package datastore

func transfer[D any](batchSize int, dataChan chan D, errorChan chan error, insertFunc func([]D) error) error {
	for {
		buf := make([]D, 0, batchSize)

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
			err := insertFunc(buf)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
