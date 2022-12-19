package yadisk

import (
	"encoding/json"
	"io"
	"net/url"
)

// CreateDirectoryResponse struct is returned by the API for create directory request.
type CreateDirectoryResponse struct {
	HRef      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

// CreateDirectory will put specified data to Yandex.Disk.
func (a *API) CreateDirectory(remotePath string) error {
	_, err := a.CreateDirectoryRequest(remotePath)
	if err != nil {
		return err
	}

	return nil
}

// CreateDirectoryRequest will make a create request and return a URL to created directory.
func (a *API) CreateDirectoryRequest(remotePath string) (*CreateDirectoryResponse, error) {
	values := url.Values{}
	values.Add("path", remotePath)

	req, err := a.scopedRequest("PUT", "/v1/disk/resources?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := CheckAPIError(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	ur, err := ParseCreateDirectoryResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

// ParseCreateDirectoryResponse tries to read and parse CreateDirectoryResponse struct.
func ParseCreateDirectoryResponse(data io.Reader) (*CreateDirectoryResponse, error) {
	dec := json.NewDecoder(data)
	var ur CreateDirectoryResponse

	if err := dec.Decode(&ur); err == io.EOF {
		// ok
	} else if err != nil {
		return nil, err
	}

	// TODO: check if there is any trash data after JSON and crash if there is.

	return &ur, nil
}
