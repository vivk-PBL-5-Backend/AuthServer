package filereader

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func ReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(bytes)
}
