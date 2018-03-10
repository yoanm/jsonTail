package jsonHelper

import (
	"encoding/json"
	"github.com/tidwall/sjson"
)

func excludeFieldsFromObject(encoded string, fieldToExcludeList []string) (resultList map[string]interface{}, err error) {
	for counter := 0; counter < len(fieldToExcludeList); counter++ {
		if encoded, err = sjson.Delete(encoded, fieldToExcludeList[counter]); err != nil {
			return nil, err;
		}
	}

	if err = json.Unmarshal([]byte(encoded), &resultList); err != nil {
		return nil, err;
	}

	return resultList, nil;
}

