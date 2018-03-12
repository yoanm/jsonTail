package types

import (
	"bytes"
)

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
