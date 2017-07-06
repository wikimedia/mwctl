package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"os/exec"
)

var (
	version = "0.0.1"

	app = kingpin.New("mwctl", "Install and drive the mediawiki-containers development environment.")

	devel       = app.Command("develop", "Prepare a service for development")
	devServices = devel.Arg("service", "Name of service(s) to develop").Required().Strings()

	apply         = app.Command("apply", "Apply a change in the dev environment")
	applyServices = apply.Arg("service", "Name of services(s) to apply changes").Required().Strings()

	test         = app.Command("test", "Test a service in a fresh container")
	testServices = test.Arg("service", "Name of service(s) whose changes you would like to test").Required().Strings()

	config      = app.Command("config", "Get and set configuration values")
	configKey   = config.Arg("key", "Key").Required().String()
	configValue = config.Arg("value", "Value").String()
)

func applyConfig(config string) ([]byte, error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	stdin, _ := cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, config)
	}()

	return cmd.CombinedOutput()
}

func Croakf(s string, args...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}

func Croak(e interface{}) {
	Croakf("%s\n", e)
}

func main() {
	app.Version(version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case devel.FullCommand():
		fmt.Println("devel")
		for _, service := range *devServices {
			fmt.Println("  " + service)
		}

	case apply.FullCommand():
		out, err := exec.Command("minikube", "ip").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The IP is: %s", out)

	case test.FullCommand():
		fmt.Println("test")

	case config.FullCommand():
		// Read the configuration if disk (if it exists)
		config, err := GetConfig(GetConfigPath())
		if err != nil {
			Croakf("Error reading config: %s\n", err)
		}

		// This is a Get
		if *configValue == "" {
			value, err := config.Get(*configKey)
			if err != nil {
				Croak(err)
			}
			fmt.Println(value)
		// This is a Set
		} else {
			// Update the config object...
			err := config.Set(*configKey, *configValue)
			if err != nil {
				Croak(err)
			}

			// ...and write it to disk
			err = WriteConfigFile(GetConfigPath(), config)
			if err != nil {
				Croak(err)
			}
		}
	}
}
