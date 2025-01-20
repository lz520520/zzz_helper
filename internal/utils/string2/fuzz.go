package string2

func GenerateFuzzRangeDict(meta string, start, end int, prefix string) (dict []string) {
	for i := start; i <= end; i++ {
		dict = append(dict, GenerateFuzzDict(meta, i)...)
	}
	for index, value := range dict {
		dict[index] = prefix + value
	}
	return
}

func GenerateFuzzDict(meta string, length int) (dict []string) {
	dict = make([]string, 0)
	if length < 1 {
	} else if length == 1 {
		metaBytes := []byte(meta)
		for _, i := range metaBytes {
			dict = append(dict, string(i))
		}
	} else {
		tmpDict := GenerateFuzzDict(meta, length-1)
		singleDict := GenerateFuzzDict(meta, 1)
		for _, i := range singleDict {
			for _, j := range tmpDict {
				dict = append(dict, i+j)
			}
		}

	}
	return
}
