package main

import (
	"os"
	"flag"
	"fmt"
	"bytes"
	"github.com/hpcloud/tail"
	"github.com/yoanm/jsonTail/helpers/types"
	"github.com/yoanm/jsonTail/helpers/json"
)

func displayLine(line *tail.Line, options types.JsonTailOptions) (error) {
	if line.Err != nil { return line.Err;}

	var lineBuffer bytes.Buffer;
	var err error;

	err = jsonHelper.Prettify(&lineBuffer, line.Text, options);
	if err != nil { return err; }

	if options.ShowDate == true {
		fmt.Printf("[%s]", line.Time)
	}

	lineBuffer.WriteTo(os.Stdout);
	fmt.Printf("\n");

	return nil;
}

func main() {

	var options types.JsonTailOptions;

	flag.Var(&(options.OnlyFieldList), "only", "Only output this field (multiple flag allow). Take predescence over --exclude")
	flag.Var(&(options.ExcludedFieldList), "exclude", "Field to exclude (multiple flag allow)")
	flag.BoolVar(&(options.FollowFile), "f", false, "follow")
	flag.BoolVar(&(options.ReopenFile), "F", false, "follow & re-open")
	flag.BoolVar(&(options.ShowDate), "date", false, "Display date for each lines")

	flag.Parse();

	options.FilePath = flag.Arg(0);

	fmt.Printf(
		"command: %s, file: %s, follow: %t, reopen: %t, showDate %t\n",
		os.Args[0],
		options.FilePath,
		options.FollowFile,
		options.ReopenFile,
		options.ShowDate,
	)

	message, code := options.CheckIntegrity();

	if (code != 0) {
		fmt.Printf(message)
		os.Exit(code)
	}

	fileLineStream, err := tail.TailFile(
		options.FilePath,
		tail.Config{
			Follow: (options.FollowFile || options.ReopenFile),
			ReOpen: (options.FollowFile && options.ReopenFile),
		})

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
