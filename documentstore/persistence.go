package documentstore

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/vaiktorg/grimoire/helpers"

	"github.com/vaiktorg/grimoire/errs"
)

// Creates a directory in the documentstore
func (d *Dir) DirToFS() error {
	if !d.Meta.IsProtected {
		err := os.MkdirAll(d.Meta.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) FileToFS() error {
	if !f.Meta.IsProtected {
		file, err := os.Create(f.Meta.Path)
		if err != nil {
			return err
		}
		defer file.Close()

		if f.Data != nil && len(f.Data) > 0 {
			_, err = file.Write(f.Data)
			if err != nil {
				return err
			}
			f.Meta.size = int64(cap(f.Data))
			f.Meta.Size = strconv.Itoa(int(f.Meta.size))
		}
	}
	return nil
}

// CreateStruct creates structure set on documentstore
func (d *Dir) CreateStructOnFS() error {

	err := d.DirToFS()
	if err != nil {
		return err
	}

	for _, file := range d.Files {
		err := file.FileToFS()
		if err != nil {
			return err
		}
	}

	for _, item := range d.Dirs {
		err := item.CreateStructOnFS()
		if err != nil {
			return err
		}

	}

	return nil
}

type Format string

const (
	JSON Format = "json"
	YAML        = "yaml"
	GOB         = "gob"
	XML         = "xml"
)

// Loads json file to populate directory
func (d *Dir) Load(fmt Format) {
	file := helpers.OpenFile("FSState_" + helpers.MakeTimestampNum() + "." + string(fmt))
	defer file.Close()

	switch fmt {
	case XML:
		err := xml.NewDecoder(file).Decode(d)
		errs.Must(err)
	case GOB:
		err := gob.NewDecoder(file).Decode(d)
		errs.Must(err)
	case YAML:
		err := yaml.NewDecoder(file).Decode(d)
		errs.Must(err)
	case JSON:
		err := json.NewDecoder(file).Decode(d)
		errs.Must(err)
	}
}

// Presists directory for later reference and persistence.
func (d *Dir) Persist(fmt Format) {
	file := helpers.OpenFile("FSState_" + helpers.MakeTimestampNum() + "." + string(fmt))
	defer file.Close()

	switch fmt {
	case XML:
		err := xml.NewEncoder(file).Encode(file)
		errs.Must(err)
	case GOB:
		err := gob.NewEncoder(file).Encode(file)
		errs.Must(err)
	case YAML:
		err := yaml.NewEncoder(file).Encode(file)
		errs.Must(err)
	case JSON:
		err := json.NewEncoder(file).Encode(d)
		errs.Must(err)
	}
}
