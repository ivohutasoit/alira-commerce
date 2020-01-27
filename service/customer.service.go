package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ivohutasoit/alira-commerce/model"

	"github.com/ivohutasoit/alira/model/domain"

	"github.com/ivohutasoit/alira/util"

	alira "github.com/ivohutasoit/alira/model"
	"github.com/ivohutasoit/alira/model/commerce"
)

type CustomerService struct{}

func (s *CustomerService) Create(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough parameter")
	}
	var code, username, email, mobile, firstName, lastName, token string
	payment := false
	for i, v := range args {
		switch i {
		case 6:
			param, ok := v.(int)
			if !ok {
				return nil, errors.New("plain text parameter not type integer")
			}
			if param == 1 {
				payment = true
			}
		default:
			param, ok := v.(string)
			if !ok {
				return nil, errors.New("plain text parameter not type string")
			}
			switch i {
			case 1:
				username = strings.ToLower(strings.TrimSpace(param))
			case 2:
				email = strings.ToLower(strings.TrimSpace(param))
			case 3:
				mobile = strings.TrimSpace(param)
			case 4:
				firstName = strings.Title(strings.TrimSpace(param))
			case 5:
				lastName = strings.Title(strings.TrimSpace(param))
			case 7:
				token = param
			default:
				code = strings.ToUpper(strings.TrimSpace(param))
			}

		}
	}

	customer := &commerce.Customer{}
	alira.GetDatabase().Where("code = ?", code).First(customer)
	if customer.BaseModel.ID != "" {
		return nil, errors.New("customer id has been used")
	}

	data := map[string]string{
		"username":   username,
		"email":      email,
		"mobile":     mobile,
		"first_name": firstName,
		"last_name":  lastName,
	}
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://localhost:9000/api/alpha/account", bytes.NewBuffer(payload))
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
	var response domain.Response
	json.Unmarshal(body, &response)
	if response.Code != http.StatusCreated {
		return nil, errors.New(response.Error)
	}

	customer = &commerce.Customer{
		UserID:  response.Data["user_id"],
		Code:    code,
		Status:  "ACTIVE",
		Payment: payment,
	}
	alira.GetDatabase().Create(&customer)

	return map[interface{}]interface{}{
		"status":   "SUCCESS",
		"message":  "Customer has been created",
		"customer": customer,
	}, nil
}

func (s *CustomerService) Search(args ...interface{}) (map[interface{}]interface{}, error) {
	var customers []commerce.Customer
	token, _ := args[0].(string)
	page, _ := strconv.Atoi(args[1].(string))
	limit, _ := strconv.Atoi("3")
	if len(args) == 2 {
		paginator := util.Initialize(&util.Parameter{
			Database: alira.GetDatabase().Where("status = ?", "ACTIVE"),
			Page:     page,
			Limit:    limit,
			Result:   &customers,
			OrderBy:  []string{"code"},
			ShowSQL:  true,
		})
		records := make([]model.CustomerRecord, 0)
		for _, v := range customers {
			req, err := http.NewRequest("GET",
				fmt.Sprintf("http://localhost:9000/api/alpha/account/%s", v.UserID), nil)
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
			var response domain.Response
			json.Unmarshal(body, &response)
			if response.Code != http.StatusOK {
				return nil, errors.New(response.Error)
			}
			fmt.Println(response.Data)
			records = append(records, model.CustomerRecord{
				ID:       v.BaseModel.ID,
				Code:     v.Code,
				Status:   v.Status,
				Payment:  v.Payment,
				Name:     response.Data["name"],
				Email:    response.Data["email"],
				Mobile:   response.Data["mobile"],
				Username: response.Data["username"],
			})
		}
		paginator.Records = records
		return map[interface{}]interface{}{
			"status":    "SUCCESS",
			"paginator": paginator,
		}, nil
	}
	return map[interface{}]interface{}{}, nil
}
