package types

type Dictionary struct {
	Values map[string]string
}

func (r *MagicaReader) GetDictionary() Dictionary {
	items := r.GetInt32()
	values := make(map[string]string)
	for i := 0; i < items; i++ {
		key := r.GetString()
		value := r.GetString()
		values[key] = value
	}

	return Dictionary{Values: values}
}