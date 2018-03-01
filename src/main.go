package main

import (
	"os"
	"flag"
	"fmt"
	"bytes"
	"reflect"
	"encoding/json"
	"github.com/hpcloud/tail"
)

func exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

func isReadable(path string) (bool) {
	file, err := os.Open(path);

	if err != nil {
		return false
	} else {
		err = file.Close()
		if err != nil {
			return false
		}
	}

	return true
}

func displayLine(line *tail.Line, showDate bool, fieldToRemovedList []string) (error) {
	if line.Err != nil { return line.Err;}

	// only in case decoded is an object
	var lineBuffer *bytes.Buffer;
	var err error;
	if len(fieldToRemovedList) > 0 && line.Text[0] == '{' {
		lineBuffer, err = prettifyJsonObject(line.Text);
	} else {
		lineBuffer, err = prettifyJsonString(line.Text);
	}
	if err != nil { return err; }

	if showDate == true {
		fmt.Printf("[%s]", line.Time)
	}
	lineBuffer.WriteTo(os.Stdout);
	fmt.Printf("\n");

	return nil;
}

func prettifyJsonObject(encoded string) (*bytes.Buffer, error) {
	//decoded, err := decodeJsonObject(encoded);
	//if err != nil { return nil, err; }

	//printFields(decoded);

	return prettifyJsonString(encoded);

}

func printFields(b *interface{}) {
	val := reflect.ValueOf(b)
	for i := 0; i < val.Type().NumField(); i++ {
		fmt.Println(val.Type().Field(i).Name)
	}
}

func decodeJsonObject(encoded string) (*interface{}, error) {
	var decoded interface{};
	err := json.Unmarshal([]byte(encoded), &decoded);
	if err != nil {
		return nil, err;
	}

	return &decoded, nil;
}

func prettifyJsonString (encoded string) (*bytes.Buffer, error) {
	var decoded bytes.Buffer;

	if err := json.Indent(&decoded, []byte(encoded), "", "  "); err != nil {
		return nil, err;
	}

	return &decoded, nil;
}

func main() {

	var followFile bool;
	var reopenFile bool;
	var showDate bool;

	flag.BoolVar(&followFile, "f", false, "follow")
	flag.BoolVar(&reopenFile, "F", false, "follow & re-open")
	flag.BoolVar(&showDate, "date", false, "Display date for each lines")

	flag.Parse();

	logFilePath := flag.Arg(0);

	fmt.Printf(
		"command: %s, file: %s, follow: %t, reopen: %t, showDate %t\n",
		os.Args[0],
		logFilePath,
		followFile,
		reopenFile,
		showDate,
	)

	if !exists(logFilePath) {
		fmt.Printf("File %s does not exist\n", logFilePath)
		os.Exit(1)
	} else if !isReadable(logFilePath) {
		fmt.Printf("File %s is not readable\n", logFilePath)
		os.Exit(2)
	}

	fileLineStream, err := tail.TailFile(
		logFilePath,
		tail.Config{
			Follow: (followFile || reopenFile),
			ReOpen: (followFile && reopenFile),
		})

	if err != nil {
		fmt.Printf("Error during file opening: %s\n", err)
		os.Exit(3)
	}

	var fieldToRemoveList []string;
	fieldToRemoveList[0] = "plop";

	for line := range fileLineStream.Lines {
		err := displayLine(line, showDate, fieldToRemoveList)
		if err != nil {
			fmt.Printf("Error during line display: %s =>%s\n", line.Text, err)
		}
	}
}
