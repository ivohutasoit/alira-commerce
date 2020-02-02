package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ivohutasoit/alira/util"
)

type User struct{}

func (s *User) ChangeUserPin(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("not enough parameter")
	}
	token, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not type string")
	}
	pin, ok := args[1].(string)
	if !ok {
		return nil, errors.New("plain text is not type string")
	}

	data := map[string]interface{}{
		"pin": pin,
	}
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://localhost:9000/api/alpha/account/pin", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	parser := &util.Parser{}
	response, err := parser.UnmarshalResponse(body, http.StatusAccepted, nil)
	if err != nil {
		return nil, err
	}

	return map[interface{}]interface{}{
		"message": response.Message,
	}, nil
}
