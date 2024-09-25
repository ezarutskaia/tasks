package pdf

import (
	"time"
	"bytes"
	"context"
	"io/ioutil"
    "encoding/json"
    "net/http"
)

type Pdf struct {
	URL string
}

func (pdf *Pdf) TaskToPdf(tasks []*TaskDTO) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

	taskJson, err := json.Marshal(tasks)
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