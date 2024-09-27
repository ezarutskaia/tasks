package pdf

import (
	"time"
	"bytes"
	"math/rand"
	"context"
	"io/ioutil"
    "encoding/json"
    "net/http"
)

type Pdf struct {
	URL string
}

func (pdf *Pdf) TaskToPdf(task *TaskDTO) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rand.Intn(300))*time.Second)
    defer cancel()

	taskJson, err := json.Marshal(task)
    if err != nil {
        return nil, err
    }

	req, err := http.NewRequestWithContext(ctx, "POST", pdf.URL, bytes.NewBuffer(taskJson))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return body, nil
		}
	}
}