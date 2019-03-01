package h2memcache

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
)

var ErrUnauthorized = errors.New("cache: Unauthorized")
var ErrNotFound = errors.New("cache: Entry not found")

type Cache struct {
	client *http.Client
	url    string
	apikey string
}

func NewCache(httpClient *http.Client, url string, apikey string) *Cache {
	c := &Cache{
		client: httpClient,
		url:    strings.TrimRight(url, "/"),
		apikey: apikey,
	}

	if c.client == nil {
		c.client = http.DefaultClient
	}

	return c
}

func (c *Cache) Get(key string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/items/"+key, nil)
	if err != nil {
		return nil, err
	}

	if c.apikey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apikey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}

	b.ReadFrom(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		return b.Bytes(), nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	default:
		return nil, errors.New("Unknown error")
	}
}

func (c *Cache) Set(key string, value []byte) error {
	b := bytes.NewBuffer(value)

	req, err := http.NewRequest(http.MethodPut, c.url+"/items/"+key, b)
	if err != nil {
		return err
	}

	if c.apikey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apikey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return errors.New("Unknown error")
	}
}

func (c *Cache) Delete(key string) error {
	req, err := http.NewRequest(http.MethodDelete, c.url+"/items/"+key, nil)
	if err != nil {
		return err
	}

	if c.apikey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apikey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return errors.New("Unknown error")
	}
}

func (c *Cache) Clear() error {
	req, err := http.NewRequest(http.MethodDelete, c.url+"/items", nil)
	if err != nil {
		return err
	}

	if c.apikey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apikey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return errors.New("Unknown error")
	}
}
