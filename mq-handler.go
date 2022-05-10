package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

func MQHandler(ctx context.Context, msg *message) {
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
