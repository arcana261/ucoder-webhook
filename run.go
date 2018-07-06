package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func run(dir string, name string, args ...string) {
	var buffer bytes.Buffer

	for i, value := range args {
		if i > 0 {
			buffer.WriteString(" ")
		}

		buffer.WriteString(value)
	}

	argument := buffer.String()

	log.Printf("'%s %s'\n", name, argument)

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = nil
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Printf("Error calling '%s %s' resulted in error %s\n", name, argument, err)
	}
}
