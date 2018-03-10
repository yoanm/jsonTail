package jsonHelper

import "github.com/tidwall/gjson"

func extractFieldListFromObject(encoded string, fieldToMatchList []string) (resultList map[string]interface{}, err error) {
	gjsonGlobalObject := gjson.Parse(encoded);

	resultList = make(map[string]interface{});
	for counter := 0; counter < len(fieldToMatchList); counter++ {
		gjsonResult := gjsonGlobalObject.Get(fieldToMatchList[counter]);

		if gjsonResult.Exists() {
			resultList[fieldToMatchList[counter]] = gjsonResult.Value();
		}
	}

	return resultList, nil;
}

