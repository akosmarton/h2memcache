package h2memcache

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var ErrUnauthorized = errors.New("cache: Unauthorized")
var ErrNotFound = errors.New("cache: Key not found")

// Cache
type Cache struct {
	client *http.Client
	url    string
	apikey string
}

// NewCache
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

// Get item
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
		return nil, fmt.Errorf("Unknown status code: %d", resp.StatusCode)
	}
}

// Set item
func (c *Cache) Set(key string, value []byte, exp int) error {
	b := bytes.NewBuffer(value)

	req, err := http.NewRequest(http.MethodPut, c.url+"/items/"+key, b)
	if err != nil {
		return err
	}

	if c.apikey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apikey)
	}

	if exp > 0 {
		req.Header.Set("X-Cache-Expire", strconv.Itoa(exp))
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
		return fmt.Errorf("Unknown status code: %d", resp.StatusCode)
	}
}

// Delete item
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
		return fmt.Errorf("Unknown status code: %d", resp.StatusCode)
	}
}

// Clear entire cache
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
		return fmt.Errorf("Unknown status code: %d", resp.StatusCode)
	}
}
