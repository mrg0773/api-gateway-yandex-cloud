package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	TargetURL        string        `envconfig:"TARGET_URL" required:"true"`
	AvailableMethods []string      `envconfig:"AVAILABLE_METHODS" default:"GET,POST,PUT"`
	AttemptLimit     int           `envconfig:"ATTEMPT_LIMIT" required:"true"`
	PauseAttempt     time.Duration `envconfig:"PAUSE_ATTEMPT" required:"true"`

	QueueName string `envconfig:"QUEUE_NAME" required:"true"`
}

var cfg *configuration

// HTTP http requests to yandex cloud
func HTTP(w http.ResponseWriter, r *http.Request) {
	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Printf("failed to load configuration: %s", err)
		return
	}

	method := r.Method
	if !inSlice(method, cfg.AvailableMethods) {
		err := fmt.Errorf("method %s is not availabel", method)
		fmt.Printf("method check: available method %s", cfg.AvailableMethods)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	r.Body.Close()

	msg := message{
		Target:    cfg.TargetURL,
		Endpoint:  r.RequestURI,
		Method:    method,
		Headers:   r.Header,
		PathQuery: r.URL.Query(),

		Body:         body,
		AttemptLimit: cfg.AttemptLimit,
		PauseAttempt: cfg.PauseAttempt,
	}

	messageID, err := sendMessage(r.Context(), msg)
	if err != nil {
		fmt.Printf("send request to mq: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	fmt.Printf("send request to mq id: %s ", messageID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(messageID))
}

func MQ(ctx context.Context, msg *message) {
	url := msg.Target + msg.Endpoint

	req, err := http.NewRequestWithContext(ctx, url, msg.Method, bytes.NewReader(msg.Body))
	if err != nil {
		fmt.Printf("create request: %s", err)
		return
	}

	for key, params := range msg.PathQuery {
		for _, val := range params {
			req.URL.Query().Add(key, val)
		}
	}

	for key, param := range msg.Headers {
		req.Header.Add(key, param[0])
	}

	client := &http.Client{}
	for i := 0; i < msg.AttemptLimit; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("%d - send request: %s", i, err)
			time.Sleep(msg.PauseAttempt)
			continue
		}
		code := res.StatusCode
		fmt.Printf("%d - response status code", code)
		if code >= 200 || code <= 299 {
			return
		}
	}
}
