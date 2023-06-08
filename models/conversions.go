package models

func DbIdArrayToInt64Array(arr []DbId) []int64 {
	int64Array := make([]int64, len(arr))
	for i, value := range arr {
		int64Array[i] = int64(value)
	}
	return int64Array
}

func Int64ArrayToDbIdArray(arr []int64) []DbId {
	dbIdArray := make([]DbId, len(arr))
	for i, value := range arr {
		dbIdArray[i] = DbId(value)
	}
	return dbIdArray
}
