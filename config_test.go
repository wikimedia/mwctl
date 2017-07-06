package main

import (
	"io/ioutil"
	"testing"
)

func TestConfigGetSet(t *testing.T) {
	config := NewConfig()
	err := config.Set("service.restbase.source", "~/dev/src/git/restbase")
	if err != nil {
		t.Fail()
	}
	if val, _ := config.Get("service.restbase.source"); val != "~/dev/src/git/restbase" {
		t.Fail()
	}

	// Use an invalid key (this should fail)
	if config.Set("bogus.key", "correspondingly bogus value") == nil {
		t.Fail()
	}
}

func TestReadWriteConfigFile(t *testing.T) {
	// Create a new config object...
	config := NewConfig()
	config.Set("service.restbase.source", "~/dev/src/git/restbase")

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
	if len(config.Services) != 1 {
		t.Fail()
	}

	if val, _ := config.Get("service.restbase.source"); val != "~/dev/src/git/restbase" {
		t.Fail()
	}
}
