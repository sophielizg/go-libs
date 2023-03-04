package datastore

type UpdateOption interface {
	Option
}

type IncludeExcludeOption struct {
	Include bool
	Exclude bool
}

func (o *IncludeExcludeOption) Name() string {
	return "IncludeExcludeOption"
}

func (o *IncludeExcludeOption) OverrideSupported() bool {
	return true
}

func ApplyIncludeExcludeOptions(dataRowFields DataRowFields, options Options[*IncludeExcludeOption]) error {
	if len(options) == 0 {
		return nil
	}

	// Get any option to test whether including or excluding
	var option *IncludeExcludeOption
	for _, option = range options {
		break
	}

	includeFieldNames := map[string]bool{}
	for fieldName := range options {
		err := validateDataRowFieldName(fieldName, dataRowFields)
		if err != nil {
			return err
		}

		if option.Include {
			includeFieldNames[fieldName] = true
		} else {
			delete(dataRowFields, fieldName)
		}
	}

	if option.Include {
		for fieldName := range dataRowFields {
			if !includeFieldNames[fieldName] {
				delete(dataRowFields, fieldName)
			}
		}
	}

	return nil
}
