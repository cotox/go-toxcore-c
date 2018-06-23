package tox

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"unsafe"
)

func safeptr(b []byte) unsafe.Pointer {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(h.Data)
}

func toxerr(errno interface{}) error {
	return fmt.Errorf("toxcore error: %v", errno)
}

func toxerrf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func assert(condition bool, message string) {
	if !condition {
		log.Panic(message)
	}
}

var toxdebug = false

func SetDebug(debug bool) {
	toxdebug = debug
}

var loglevel = 0

func SetLogLevel(level int) {
	loglevel = level
}

func FileExist(fname string) bool {
	_, err := os.Stat(fname)
	if err != nil {
		return false
	}
	return true
}

func (this *Tox) WriteSavedata(fname string) error {
	if !FileExist(fname) {
		// TODO: choose a better file mode
		err := ioutil.WriteFile(fname, this.GetSavedata(), 0755)
		if err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		liveData := this.GetSavedata()
		if bytes.Compare(data, liveData) != 0 {
			// TODO: choose a better file mode
			err := ioutil.WriteFile(fname, this.GetSavedata(), 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Tox) LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

func LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

func ConnStatusString(status ConnectionType) (s string) {
	switch status {
	case ConnectionNone:
		s = "CONNECTION_NONE"
	case ConnectionTCP:
		s = "CONNECTION_TCP"
	case ConnectionUDP:
		s = "CONNECTION_UDP"
	}
	return
}
