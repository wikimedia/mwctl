package main

import (
	"io/ioutil"
	"testing"
)

func TestReadWriteConfigFile(t *testing.T) {
	// Create a new config object...
	config := NewConfig()
	config.Sources["restbase"] = "git@github.com:wikimedia/restbase"

	// ...create a temporary file
	tmpfile, err := ioutil.TempFile("", "mwctl-cfg")
	if err != nil {
		t.Error("Error creating temp file")
	}

	// ...write the config to temporary file
	err = WriteConfigFile(tmpfile.Name(), config)
	if err != nil {
		t.Error("Error writing config file")
	}

	// ...read it back
	config, err = ReadConfigFile(tmpfile.Name())
	if err != nil {
		t.Error("Error reading config file")
	}

	// ...and validate the result
	if len(config.Sources) != 1 {
		t.Fail()
	}

	if config.Sources["restbase"] != "git@github.com:wikimedia/restbase" {
		t.Fail()
	}
}
