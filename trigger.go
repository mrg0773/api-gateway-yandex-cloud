package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

func MQ() (err error) {
	fmt.Println("mq trigger")

	ctx := context.Background()

	msg, err := recieveMessage(ctx)
	if err != nil {
		fmt.Printf("recieve message: %s", err)
		return
	}

	fmt.Printf("received msg from mq: %+v\n", msg)
	url := msg.Target + msg.Endpoint

	req, err := http.NewRequestWithContext(ctx, url, msg.Method, bytes.NewReader(msg.Body))
	if err != nil {
		fmt.Printf("create request: %s", err)
		return
	}

	q := req.URL.Query()
	for key, params := range msg.PathQuery {
		for _, val := range params {
			q.Add(key, val)
		}
	}

	for key, param := range msg.Headers {
		req.Header.Add(key, param[0])
	}

	fmt.Printf("sent request:\n%s", formatRequest(req))

	client := &http.Client{}
	for i := 0; i < msg.AttemptLimit; i++ {
		var res *http.Response
		res, err = client.Do(req)
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

	return
}
