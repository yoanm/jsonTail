package types

import (
	"fmt"
	"github.com/yoanm/jsonTail/helpers/file"
)

type JsonTailOptions struct {
	FilePath string;
	FollowFile bool;
	ReopenFile bool;
	ShowDate bool;
	ExcludedFieldList ArrayStringFlag;
	OnlyFieldList ArrayStringFlag;
}

func (this *JsonTailOptions) CheckIntegrity() (errorMessage string, code int) {
	if !file.FileExists((*this).FilePath) {
		return fmt.Sprintf("File %s does not exist", (*this).FilePath), 1;
	} else if !file.IsFileReadable((*this).FilePath) {
		return fmt.Sprintf("File %s is not readable", (*this).FilePath), 2;
	}

	if (len((*this).OnlyFieldList) > 0 && len((*this).ExcludedFieldList) > 0) {
		return fmt.Sprintf("You cannot target some fields and exclude some others in same time !"), 3;
	}

	return "", 0;
}

