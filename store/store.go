package store

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vaiktorg/grimoire/helpers"

	"github.com/vaiktorg/grimoire/errs"
)

// TODO: Make Serializer concurrency safe.
type Serializer struct {
	name            string
	filepath        string
	typp            reflect.Type
	dataToSerialize interface{}
}

func Init(pathToFile string, dataToSerialize interface{}) (s *Serializer) {
	s = new(Serializer)
	s.filepath = pathToFile
	s.name = fmt.Sprintf("%T\n", dataToSerialize) + "_" + helpers.MakeTimestampNum()
	if reflect.ValueOf(dataToSerialize).Kind() == reflect.Ptr {
		s.dataToSerialize = dataToSerialize
	} else {
		s.dataToSerialize = &dataToSerialize
	}
	return
}

func (s *Serializer) Load(filePath string) {
	file := helpers.OpenFile(filePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(s.dataToSerialize)
	errs.Must(err)
}

func (s *Serializer) Persist() {
	file := helpers.OpenFile(s.filepath)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err := encoder.Encode(s.dataToSerialize)
	errs.Must(err)
}
