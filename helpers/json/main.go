package jsonHelper

import (
	"bytes"
	"errors"
	"encoding/json"
	"github.com/yoanm/jsonTail/helpers/types"
	"github.com/tidwall/gjson"
	//"github.com/hokaccha/go-prettyjson"
)

const indentationString string = "  ";

func Prettify(buffer *bytes.Buffer, encoded string, options types.JsonTailOptions) (err error) {
	if !gjson.Valid(encoded) {
		return errors.New("Invalid JSON");
	}

	if encoded[0] == '{' {
		decoded, err := prettifyJsonObject(encoded, options);
		if err != nil {return err;}

		return indentJson(buffer, decoded);
	}

	return indentString(buffer, encoded);
}

func prettifyJsonObject(encoded string, options types.JsonTailOptions) (decoded interface{}, err error) {
	if (len(options.ExcludedFieldList) > 0) {
		return excludeFieldsFromObject(encoded, options.ExcludedFieldList);
	} else if (len(options.OnlyFieldList) > 0) {
		return extractFieldListFromObject(encoded, options.OnlyFieldList);
	}

	// Just decode it (nothing to do)
	var resultList map[string]interface{};
	if err := json.Unmarshal([]byte(encoded), &resultList); err != nil {
		return nil, err;
	}

	return resultList, nil;

}

/*func excludeFieldsFromObject(decoded map[string]interface{}, fieldToRemoveList []string) (map[string]interface{}) {
	for counter := 0; counter < len(fieldToRemoveList); counter++ {
		delete(decoded, fieldToRemoveList[counter]);
	}

	return decoded;
}*/

func indentString(buffer *bytes.Buffer, encoded string) (error) {
	if err := json.Indent(buffer, []byte(encoded), "", indentationString); err != nil {
		return err;
	}

	return nil;
}

func indentJson(buffer *bytes.Buffer, decoded interface{})  (error) {
	var encoded []byte;
	var err error;
	//if encoded, err = prettyjson.Marshal(decoded); err != nil {
	if encoded, err = json.MarshalIndent(decoded, "", indentationString); err != nil {
		return err;
	}

	buffer.Write(encoded);

	return nil;
}
