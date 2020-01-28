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

func (s *CustomerService) Get(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("not enough parameter")
	}
	token, _ := args[0].(string)
	id, _ := args[1].(string)
	customer := &commerce.Customer{}
	alira.GetDatabase().Where("id = ?", id).First(&customer)
	if customer.BaseModel.ID == "" {
		return nil, errors.New("invalid customer")
	}
	custUser := &commerce.CustomerUser{}
	alira.GetDatabase().Where("customer_id = ? AND role = ?", id, "OWNER").First(&custUser)
	if custUser.BaseModel.ID == "" {
		return nil, errors.New("invalid customer")
	}

	req, err := http.NewRequest("GET",
		fmt.Sprintf("http://localhost:9000/api/alpha/account/%s", custUser.UserID), nil)
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

	return map[interface{}]interface{}{
		"customer": customer,
		"profile":  response.Data,
	}, nil
}

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
	customer = &commerce.Customer{
		Code:    code,
		Status:  "ACTIVE",
		Payment: payment,
	}
	alira.GetDatabase().Create(&customer)

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
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Code != http.StatusCreated {
		return nil, errors.New(response.Error)
	}

	var user domain.User
	if err := json.Unmarshal([]byte(response.Data), &user); err != nil {
		return nil, err
	}

	custUser := &commerce.CustomerUser{
		CustomerID: customer.BaseModel.ID,
		UserID:     user.BaseModel.ID,
		Role:       "OWNER",
	}
	alira.GetDatabase().Create(&custUser)

	return map[interface{}]interface{}{
		"status":   "SUCCESS",
		"message":  "Customer has been created",
		"customer": customer,
		"owner":    custUser,
	}, nil
}

func (s *CustomerService) Search(args ...interface{}) (map[interface{}]interface{}, error) {
	var customers, count []model.CustomerProfile
	token, _ := args[0].(string)
	page, _ := strconv.Atoi(args[1].(string))
	limit, _ := strconv.Atoi("5")
	if len(args) == 2 {
		paginator := util.NewPaginator(&util.PaginatorParameter{
			Database: alira.GetDatabase().Table("customers c").Select("c.id, c.code, c.status, c.payment, cu.user_id").Joins("JOIN customer_users cu ON cu.customer_id=c.id AND cu.role='OWNER'"),
			Count:    &count,
			Page:     page,
			Limit:    limit,
			Result:   &customers,
			OrderBy:  []string{"c.code"},
			ShowSQL:  true,
		})
		//fmt.Println(paginator)
		temp := make([]model.CustomerProfile, 0)
		for _, customer := range customers {
			req, err := http.NewRequest("GET",
				fmt.Sprintf("http://localhost:9000/api/alpha/account/%s", customer.UserID), nil)
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
			if err := json.Unmarshal(body, &response); err != nil {
				return nil, err
			}

			if response.Code != http.StatusOK {
				return nil, errors.New(response.Error)
			}
			if err := json.Unmarshal([]byte(response.Data), &customer); err != nil {
				return nil, errors.New(response.Error)
			}
			temp = append(temp, customer)
		}
		paginator.Records = temp
		return map[interface{}]interface{}{
			"status":    "SUCCESS",
			"paginator": paginator,
		}, nil
	}
	return map[interface{}]interface{}{}, nil
}
