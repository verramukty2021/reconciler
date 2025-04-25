package util

import (
	"errors"
	"fmt"
)

func DeleteFromSliceMap(mapRecords map[string][]interface{}, key string, startIndexDelete int, endIndexDelete int) (bool, error) {

	slice, isExist := mapRecords[key]

	if isExist {
		if startIndexDelete >= 0 && endIndexDelete <= len(slice) && startIndexDelete < endIndexDelete {
			// delete from slice
			mapRecords[key] = append(slice[:startIndexDelete], slice[endIndexDelete:]...)
			return true, nil
		} else {
			return false, errors.New("Invalid start or end index")
		}
	} else {
		fmt.Println("Key not found")
		return false, nil
	}

}
