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

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira/database/commerce"
	"github.com/ivohutasoit/alira/messaging"
	"github.com/ivohutasoit/alira/util"
)

type Customer struct{}

func (s *Customer) Get(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("not enough parameter")
	}
	token, _ := args[0].(string)
	id, _ := args[1].(string)
	customer := &commerce.Customer{}
	alira.GetConnection().Where("id = ?", id).First(&customer)
	if customer.Model.ID == "" {
		return nil, errors.New("invalid customer")
	}
	custUser := &commerce.CustomerUser{}
	alira.GetConnection().Where("customer_id = ? AND role = ?", id, "OWNER").First(&custUser)
	if custUser.Model.ID == "" {
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
	var userProfile messaging.UserProfile
	parser := &util.Parser{}
	parser.UnmarshalResponse(body, http.StatusOK, &userProfile)

	return map[interface{}]interface{}{
		"customer": customer,
		"profile":  userProfile,
	}, nil
}

func (s *Customer) Create(args ...interface{}) (map[interface{}]interface{}, error) {
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
	alira.GetConnection().Where("code = ?", code).First(customer)
	if customer.Model.ID != "" {
		return nil, errors.New("customer id has been used")
	}
	customer = &commerce.Customer{
		Code:    code,
		Status:  "ACTIVE",
		Payment: payment,
	}

	data := map[string]interface{}{
		"username":   username,
		"email":      email,
		"mobile":     mobile,
		"first_name": firstName,
		"last_name":  lastName,
		"active":     true,
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
	var userProfile messaging.UserProfile

	parser := &util.Parser{}
	parser.UnmarshalResponse(body, http.StatusCreated, &userProfile)

	alira.GetConnection().Create(&customer)
	custUser := &commerce.CustomerUser{
		CustomerID: customer.Model.ID,
		UserID:     userProfile.ID,
		Role:       "OWNER",
	}
	alira.GetConnection().Create(&custUser)

	return map[interface{}]interface{}{
		"status":   "SUCCESS",
		"message":  "Customer has been created",
		"customer": customer,
		"owner":    custUser,
	}, nil
}

func (s *Customer) Search(args ...interface{}) (map[interface{}]interface{}, error) {
	var customers, count []model.CustomerProfile
	token, _ := args[0].(string)
	page, _ := strconv.Atoi(args[1].(string))
	limit, _ := strconv.Atoi("5")
	if len(args) == 2 {
		paginator := util.NewPaginator(&util.PaginatorParameter{
			Database: alira.GetConnection().Table("customers c").Select("c.id, c.code, c.status, c.payment, cu.user_id").Joins("JOIN customer_users cu ON cu.customer_id=c.id AND cu.role='OWNER'"),
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
			fmt.Println(customer.UserID)
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
			var userProfile messaging.UserProfile

			parser := &util.Parser{}
			parser.UnmarshalResponse(body, http.StatusOK, &userProfile)

			customer.UserID = userProfile.ID
			customer.Username = userProfile.Username
			customer.Email = userProfile.Email
			customer.Mobile = userProfile.PrimaryMobile
			customer.Name = fmt.Sprintf("%s %s", userProfile.FirstName, userProfile.LastName)

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
