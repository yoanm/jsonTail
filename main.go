package main

import (
	"os"
	"flag"
	"fmt"
	"bytes"
	"errors"
	"encoding/json"
	"github.com/hpcloud/tail"
	"github.com/tidwall/sjson"
	"github.com/tidwall/gjson"
	"github.com/hokaccha/go-prettyjson"
	"github.com/fatih/color"
)

/****************** HELPERS ********************/

/****************** HELPERS - JSON ********************/
func prettifyLine(buffer *bytes.Buffer, encoded string, options JsonTailOptions) (err error) {
	var newEncoded []byte;

	if !gjson.Valid(encoded) {
		return errors.New("Invalid JSON");
	}

	formatter := &prettyjson.Formatter{
		KeyColor:        color.New(color.FgBlue, color.Italic),
		StringColor:     color.New(color.FgGreen),
		BoolColor:       color.New(color.FgYellow),
		NumberColor:     color.New(color.FgCyan),
		NullColor:       color.New(color.FgBlack, color.Bold),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          2,
	};

	if encoded[0] == '{' {
		decoded, err := processJsonObjectString(encoded, options);
		if err != nil {return err;}

		newEncoded, err = formatter.Marshal(decoded);
	} else {
		newEncoded, err = formatter.Format([]byte(encoded));
	}

	if err != nil {return err;}

	buffer.Write(newEncoded);

	return nil;
}

func processJsonObjectString(encoded string, options JsonTailOptions) (decoded interface{}, err error) {
	if (len(options.ExcludedFieldList) > 0) {
		return excludeFieldsFromJsonObjectString(encoded, options.ExcludedFieldList);
	} else if (len(options.OnlyFieldList) > 0) {
		return extractFieldListFromJsonObjectString(encoded, options.OnlyFieldList);
	}

	// Just decode it (nothing to do)
	if err := json.Unmarshal([]byte(encoded), &decoded); err != nil {
		return nil, err;
	}

	return decoded, nil;

}

func excludeFieldsFromJsonObjectString(encoded string, fieldToExcludeList []string) (resultList map[string]interface{}, err error) {
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

func extractFieldListFromJsonObjectString(encoded string, fieldToMatchList []string) (resultList map[string]interface{}, err error) {
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
/****************** END - HELPERS - JSON ********************/

/****************** HELPERS - TYPE ********************/
type JsonTailOptions struct {
	FilePath string;
	FollowFile bool;
	ReopenFile bool;
	ShowDate bool;
	ExcludedFieldList ArrayStringFlag;
	OnlyFieldList ArrayStringFlag;
	Location *tail.SeekInfo
}

func (this *JsonTailOptions) CheckIntegrity() (errorMessage string, code int) {
	// Check if given file exist
	_, err := os.Stat((*this).FilePath);
	if err != nil || os.IsNotExist(err) {
		return fmt.Sprintf("File %s does not exist", (*this).FilePath), 1;
	}

	// Check if file is readable
	file, err := os.Open((*this).FilePath);
	if err != nil || file.Close() != nil {
		return fmt.Sprintf("File %s is not readable", (*this).FilePath), 2;
	}

	if (len((*this).OnlyFieldList) > 0 && len((*this).ExcludedFieldList) > 0) {
		return fmt.Sprintf("You cannot target some fields and exclude some others in same time !"), 3;
	}

	return "", 0;
}

type ArrayStringFlag []string;

func (this *ArrayStringFlag) String() (string) {
	var buffer bytes.Buffer;
	for counter := 0; counter < len(*this); counter++ {
		if (counter > 0) {
			buffer.WriteString(", ");
		}
		buffer.WriteString((*this)[counter]);
	}

	return buffer.String();
}

func (this *ArrayStringFlag) Set(value string) (error) {
	*this = append(*this, value);

	return nil;
}
/****************** END - HELPERS - TYPE ********************/

/****************** END - HELPERS ********************/

func displayLine(line *tail.Line, options JsonTailOptions) (err error) {
	if line.Err != nil { return line.Err;}

	var lineBuffer bytes.Buffer;

	err = prettifyLine(&lineBuffer, line.Text, options);
	if err != nil { return err; }

	if options.ShowDate == true {
		fmt.Printf("[%s]", line.Time)
	}

	lineBuffer.WriteTo(os.Stdout);
	fmt.Printf("\n");

	return nil;
}

func main() {

	var options JsonTailOptions;

	flag.Var(&(options.OnlyFieldList), "only", "Only output this field (multiple flag allow). Take predescence over --exclude")
	flag.Var(&(options.ExcludedFieldList), "exclude", "Field to exclude (multiple flag allow)")
	flag.BoolVar(&(options.FollowFile), "f", false, "follow")
	flag.BoolVar(&(options.ReopenFile), "F", false, "follow & re-open")
	flag.BoolVar(&(options.ShowDate), "date", false, "Display date for each lines")

	flag.Parse();

	options.FilePath = flag.Arg(0);

	message, code := options.CheckIntegrity();

	options.Location = &tail.SeekInfo{
		Offset: 0,
		Whence: os.SEEK_SET,
	}

	if (true == options.FollowFile) {
		options.Location.Whence = os.SEEK_END;
	}

	if (code != 0) {
		fmt.Printf(message)
		os.Exit(code)
	}

	fileLineStream, err := tail.TailFile(
		options.FilePath,
		tail.Config{
			Follow: (options.FollowFile || options.ReopenFile),
			ReOpen: (options.FollowFile && options.ReopenFile),
			Location: options.Location,
			Logger: tail.DiscardingLogger,
		});

	if err != nil {
		fmt.Printf("Error during file opening: %s\n", err)
		os.Exit(3)
	}

	for line := range fileLineStream.Lines {
		err := displayLine(line, options)
		if err != nil {
			fmt.Printf("Error during line display: %s =>%s\n", line.Text, err)
		}
	}
}
