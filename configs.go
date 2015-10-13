package common_utils

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"log"
)

var re *regexp.Regexp
var pat = "[#].*\\n|\\s+\\n|\\S+[=]|.*\n"

func init() {
	re, _ = regexp.Compile(pat)
}

func Load(configPath string, dest map[string]string) error {

	log.Println("Loading configurations...")

	// Check if config file available.
	checkFile, err := os.Stat(configPath)
	if err != nil {
		log.Println("ERROR: Could not find " + configPath)
		return err
	}

	// Open and Read from config file
	confFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	buff := make([]byte, checkFile.Size())
	confFile.Read(buff)
	confFile.Close()
	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	s2 := re.FindAllString(str, -1)

	for i := 0; i < len(s2); {
		if strings.HasPrefix(s2[i], "#") {
			i++
		} else if strings.HasSuffix(s2[i], "=") {
			key := strings.ToLower(s2[i])[0 : len(s2[i]) - 1]
			i++
			if strings.HasSuffix(s2[i], "\n") {
				val := s2[i][0 : len(s2[i]) - 1]
				if strings.HasSuffix(val, "\r") {
					val = val[0 : len(val) - 1]
				}
				i++
				dest[key] = val
			}
		} else if strings.Index(" \t\r\n", s2[i][0:1]) > -1 {
			i++
		} else {
			return errors.New("Unable to process line in cfg file containing " + s2[i])
		}
	}
	return nil
}