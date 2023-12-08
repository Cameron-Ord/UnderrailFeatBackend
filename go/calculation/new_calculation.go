package calculation

import "fmt"

func Prepare_Data(data RequestData) ([]byte, error) {
	var allAllocatedFeats []string
	var Feats []map[string]string
	var err error
	Feats = Unload_Json()
	err = Run_Calculation(Feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	fmt.Println("ALL ALLOCATED : ", allAllocatedFeats)
	jsonData, err := Marshal_Data(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}