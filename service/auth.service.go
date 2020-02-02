package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ivohutasoit/alira-commerce/messaging"

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira/database/commerce"
	"github.com/ivohutasoit/alira/util"
)

type Auth struct{}

func (s *Auth) AuthenticateUser(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough parameter")
	}
	userid, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not string type")
	}
	userid = strings.ToLower(userid)

	custUser := &commerce.CustomerUser{}
	alira.GetConnection().Where("active = ? AND (username = ? OR email = ?)",
		true, userid, userid).First(&custUser)
	if custUser.Model.ID == "" {
		return nil, errors.New("invalid login")
	}

	data := map[string]interface{}{
		"user_id": custUser.UserID,
	}
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_LOGIN")),
		bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
	var loggedUser messaging.LoggedUser
	parser := &util.Parser{}
	response, err := parser.UnmarshalResponse(body, http.StatusOK, &loggedUser)
	if err != nil {
		return nil, err
	}

	return map[interface{}]interface{}{
		"message": response.Message,
		"user":    &loggedUser,
	}, nil
}

func (s *Auth) VerifyToken(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("not enough parameter")
	}
	userid, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not string type")
	}
	token, ok := args[1].(string)
	if !ok {
		return nil, errors.New("plain text is not string type")
	}
	data := map[string]interface{}{
		"referer": userid,
		"token":   token,
		"purpose": "LOGIN",
	}
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_TOKENVERiFY")),
		bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
	var loggedUser messaging.LoggedUser
	parser := &util.Parser{}
	response, err := parser.UnmarshalResponse(body, http.StatusOK, &loggedUser)
	if err != nil {
		return nil, err
	}
	loggedUser.UserID = userid
	return map[interface{}]interface{}{
		"message": response.Message,
		"user":    &loggedUser,
	}, nil
}

func (s *Auth) RemoveSessionToken(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough parameters")
	}
	token, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not type string")
	}
	if token == " " {
		return nil, errors.New("invalid token")
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_LOGOUT")),
		nil)
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
	response, err := parser.UnmarshalResponse(body, http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return map[interface{}]interface{}{
		"message": response.Message,
	}, nil
}
