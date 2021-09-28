package logger

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type elastic struct {
	conn        *elasticsearch.Client
	index       string
	serviceName string
}

func NewElastic() *elastic {
	host := os.Getenv("ELASTIC_URI")
	username := os.Getenv("ELASTIC_USERNAME")
	password := os.Getenv("ELASTIC_PASSWORD")

	config := elasticsearch.Config{
		Addresses: []string{host},
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &elastic{
		conn:  client,
		index: fmt.Sprintf("log_%s", os.Getenv("APP_ENV")),
	}
}

func (e *elastic) SendDataToElastic(data *DataLogger, errChan chan error) {
	var ctx = context.Background()
	e.serviceName = data.Service

	if err := e.findOrCreateIndex(ctx); err != nil {
		errChan <- err
		return
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		errChan <- err
		return
	}

	req := esapi.IndexRequest{
		Index: e.index,
		Body:  bytes.NewBuffer(dataByte),
	}

	res, err := req.Do(ctx, e.conn)
	if err != nil {
		errChan <- err
		return
	}

	defer res.Body.Close()
}

func (e *elastic) findOrCreateIndex(ctx context.Context) error {
	index, err := e.conn.Indices.ExistsAlias([]string{e.index})
	if err != nil {
		return err
	}

	if index.StatusCode == http.StatusNotFound {
		req := esapi.IndicesCreateRequest{
			Index: e.generateIndex(),
			Body:  strings.NewReader(fmt.Sprintf(`{"aliases":{"%s":{"is_write_index":true}}}`, e.index)),
		}

		res, err := req.Do(ctx, e.conn)
		if err != nil {
			return err
		}

		defer res.Body.Close()
	}

	return nil
}

func (e *elastic) generateIndex() string {
	index := fmt.Sprintf("<%s-{now/d{yyyy.MM.dd|+07:00}}-000001>", e.index)
	return url.QueryEscape(index)
}
