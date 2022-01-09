package helpers

import (
	"os"

	"github.com/vaiktorg/grimoire/errs"
)

func OpenFile(filepath string) (file *os.File) {
	if _, err := os.Stat(filepath); err == nil {
		file, err = os.OpenFile(filepath, os.O_RDWR, os.ModePerm)
		errs.Must(err)
		return file
	} else {
		file, err := os.Create(filepath)
		errs.Must(err)
		return file
	}
}
