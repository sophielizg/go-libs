package mutator

type MutatableFactory[M any, PM Mutatable[M]] struct{}

func (f MutatableFactory[M, PM]) Create() PM {
	return PM(new(M))
}

func (f MutatableFactory[M, PM]) CreateFromFields(fields MappedFieldValues) (PM, error) {
	mutatable := f.Create()
	err := mutatable.Mutator().SetFields(fields)
	return mutatable, err
}

func (f MutatableFactory[M, PM]) CreateFieldValues(mutatable PM) MappedFieldValues {
	return mutatable.Mutator().GetFields()
}

func (f MutatableFactory[M, PM]) CreateFromFieldsList(fieldsList []MappedFieldValues) ([]PM, error) {
	mutatables := make([]PM, len(fieldsList))

	for i, fields := range fieldsList {
		if fields == nil {
			mutatables[i] = nil
			continue
		}

		var err error
		mutatables[i], err = f.CreateFromFields(fields)
		if err != nil {
			return nil, err
		}
	}

	return mutatables, nil
}

func (f MutatableFactory[M, PM]) CreateFieldValuesList(mutatables []PM) []MappedFieldValues {
	fieldsList := make([]MappedFieldValues, len(mutatables))

	for i, mutatable := range mutatables {
		fieldsList[i] = f.CreateFieldValues(mutatable)
	}

	return fieldsList
}
