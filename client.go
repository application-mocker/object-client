package object_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ObjectPath     = "json"
	ScopeSeparator = "/"
)

type ObjectClient struct {
	host       string
	baseScope  string
	httpClient *http.Client

	objUrl string
}

// NewObjectClient return ObjectClient with target url and base scope
func NewObjectClient(host, baseScope string) (*ObjectClient, error) {
	return NewObjectClientWithHttpClient(host, baseScope, nil)
}

func NewObjectClientWithHttpClient(host, baseScope string, client *http.Client) (*ObjectClient, error) {
	// empty host
	if len(host) == 0 {
		return nil, fmt.Errorf("Empty host of ObjectClient is not allow ")
	}
	address := host

	if client == nil {
		// use default http client
		client = &http.Client{}
	}

	oc := &ObjectClient{
		host:       address,
		baseScope:  baseScope,
		httpClient: client,
	}

	if err := oc.ping(); err != nil {
		return nil, err
	}

	oc.objUrl = fmt.Sprintf("%s/%s/%s", oc.host, ObjectPath, oc.baseScope)

	return oc, nil
}

// check host is success
func (oc *ObjectClient) ping() error {
	// todo: add auth method
	api := "/mock/code/special-http-code/997"

	resp, err := oc.httpClient.Get(oc.host + api)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 997 {
		return fmt.Errorf("Error to connect the server with host: %s, response is: %v", oc.host, body)

	}

	return nil
}

func (oc *ObjectClient) SubClient(scope string) (*ObjectClient, error) {
	return NewObjectClientWithHttpClient(
		oc.host,
		strings.Join([]string{oc.baseScope, scope}, ScopeSeparator),
		oc.httpClient)
}

// InsertOne will insert obj to ObjectMocker, and return the id. The object without extend from BaseNode.
func (oc *ObjectClient) InsertOne(obj interface{}) (string, error) {

	// to json
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(jsonObj)

	// gen request
	req, err := oc.newRequest(http.MethodPost, nil, body)
	if err != nil {
		return "", err
	}

	// invoke
	response, err := oc.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	resBody, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK || err != nil {
		return "", fmt.Errorf("Result code: %v, body: %v, err: %v ", response.StatusCode, string(resBody), err)
	}

	// parse return
	baseStruct := &BaseNode{}
	if err := json.Unmarshal(resBody, baseStruct); err != nil {
		return "", err
	}
	return baseStruct.Id, nil
}

// ListAllValue will return all value by type obj.
func (oc *ObjectClient) ListAllValue() ([]*SimpleNode, error) {
	res := make([]*SimpleNode, 0, 0)

	req, err := oc.newRequest(http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := oc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Result code: %v, body: %v, err: %v ", resp.StatusCode, string(resBody), err)
	}

	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByIdWithoutBaseStruct return object by special id, and the obj not extends from BaseNode
func (oc *ObjectClient) GetByIdWithoutBaseStruct(id string, obj interface{}) (bool, error) {
	newObj := &CommonNode{
		Value: obj,
	}

	return oc.GetById(id, newObj)
}

func (oc *ObjectClient) UpdateByIdWithoutBaseStruct(id string, obj interface{}) (bool, error) {

	queryParams := map[string]string{
		"id": id,
	}

	// to json
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	body := bytes.NewReader(jsonObj)

	// gen request
	req, err := oc.newRequest(http.MethodPut, queryParams, body)
	if err != nil {
		return false, err
	}

	// invoke
	response, err := oc.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	resBody, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK || err != nil {
		return false, fmt.Errorf("Result code: %v, body: %v, err: %v ", response.StatusCode, string(resBody), err)
	}

	return true, nil
}

// DeleteById delete special data by id, if not found, return nil of error and data.
func (oc *ObjectClient) DeleteById(id string) (*SimpleNode, error) {
	res := &SimpleNode{}
	queryParams := map[string]string{
		"id": id,
	}

	req, err := oc.newRequest(http.MethodDelete, queryParams, nil)
	if err != nil {
		return nil, err
	}

	resp, err := oc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Result code: %v, body: %v, err: %v ", resp.StatusCode, string(resBody), err)
	}
	err = json.Unmarshal(resBody, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetById will get special object by id, if not found or error occur, return false and error, when only not found,
//return false and nil. The obj must extend from.
func (oc *ObjectClient) GetById(id string, obj interface{}) (bool, error) {
	queryParams := map[string]string{
		"id": id,
	}

	req, err := oc.newRequest(http.MethodGet, queryParams, nil)
	if err != nil {
		return false, err
	}

	resp, err := oc.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Result code: %v, body: %v, err: %v ", resp.StatusCode, string(resBody), err)
	}
	err = json.Unmarshal(resBody, obj)
	if err != nil {
		return false, err
	}

	return true, nil
}

// newRequest generator new request by special input
func (oc *ObjectClient) newRequest(method string, query map[string]string, body io.Reader) (*http.Request, error) {

	request, err := http.NewRequest(method, oc.objUrl, body)
	if err != nil {
		return nil, err
	}

	if query != nil && len(query) != 0 {
		q := request.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	// todo add auth

	return request, nil
}
