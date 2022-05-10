package main

import "time"

type message struct {
	Target    string              `json:"target"`
	Endpoint  string              `json:"endpoint"`
	Method    string              `json:"method"`
	Headers   map[string][]string `json:"headers"`
	PathQuery map[string][]string `json:"path_quey"`
	Body      []byte              `json:"body"`

	AttemptLimit int           `json:"attempt_limit"`
	PauseAttempt time.Duration `json:"pause_attemp"`
}
