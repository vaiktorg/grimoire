package errs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// errorString is a trivial implementation of error.
type Error string

func (e Error) Error() string {
	return string(e)
}

//--------------------------------------

type ErrorRepository struct {
	repo []error
	c    chan error
	w    sync.WaitGroup
}

var (
	Repo ErrorRepository
)

func init() {
	Repo = ErrorRepository{
		c: make(chan error),
	}
}

func (e *ErrorRepository) Start() {
	go func() {
		for c := range e.c {
			e.repo = append(e.repo, c)
		}
	}()
}
func (e *ErrorRepository) Stop() {
	close(e.c)

	f, err := os.Create(time.Now().Format("20060102150405"))
	if err != nil {
		panic(err)
	}
	err = json.NewEncoder(f).Encode(e.repo)
	if err != nil {
		panic(err)
	}
}

func (e *ErrorRepository) ErrIsNil(err error) bool {
	if err == nil {
		return true
	}

	if e != nil && e.c != nil {
		e.c <- err
	}
	return false
}
func (e *ErrorRepository) Error(err error) string {
	if e != nil && e.c != nil && err != nil {
		e.c <- err
		return err.Error()
	}
	return ""
}

func (e *ErrorRepository) LastError() error {
	if len(e.repo) > 0 {
		return e.repo[len(e.repo)-1]
	}
	return Error("")
}

func ServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(500)
	_, _ = fmt.Fprintln(w, err)
}

func Must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
