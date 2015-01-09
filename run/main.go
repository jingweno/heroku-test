package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

type Result struct {
	Output   string
	ExitCode int
	Error    error
}

func execCmd(cmd []string) (*Result, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	out, err := c.CombinedOutput()

	var exitCode int
	if err != nil {
		exitCode = 1
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	return &Result{
		Output:   string(out),
		Error:    err,
		ExitCode: exitCode,
	}, nil
}

func send(r *Result) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	url := "http://requestb.in/1boam5g1"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	r, err := execCmd(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = send(r)
	if err != nil {
		log.Fatal(err)
	}
}
