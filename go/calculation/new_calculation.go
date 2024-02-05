package calculation

func Prepare_Data(data RequestData) ([]byte, error) {
	var allAllocatedFeats []FeatData
	var Feats []map[string]string
	var err error
	Feats = Unload_Json()
	err = Run_Calculation(Feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}

	jsonData, err := Marshal_Data(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
