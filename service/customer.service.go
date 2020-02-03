package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ivohutasoit/alira/service"

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
		fmt.Sprintf("%s%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_ACCOUNT"), "/"+custUser.UserID),
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
	var userProfile messaging.UserProfile
	parser := &util.Parser{}
	_, err = parser.UnmarshalResponse(body, http.StatusOK, &userProfile)
	if err != nil {
		fmt.Println(err.Error())
	}

	ss := &Store{}

	data, err := ss.Search("customer", customer.ID)
	if err != nil {
		return nil, err
	}

	customerProfile := &model.CustomerProfile{
		ID:        customer.ID,
		UserID:    custUser.UserID,
		Code:      customer.Code,
		Status:    customer.Status,
		Payment:   customer.Payment,
		Username:  userProfile.Username,
		Email:     userProfile.Email,
		Mobile:    userProfile.PrimaryMobile,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
	}

	fmt.Println(customerProfile)

	return map[interface{}]interface{}{
		"customer": customerProfile,
		"stores":   data["stores"],
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
		"username":      username,
		"email":         email,
		"mobile":        mobile,
		"first_name":    firstName,
		"last_name":     lastName,
		"active":        false,
		"customer_user": true,
	}
	payload, _ := json.Marshal(data)
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_ACCOUNT")),
		bytes.NewBuffer(payload))
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
	_, err = parser.UnmarshalResponse(body, http.StatusCreated, &userProfile)
	if err != nil {
		return nil, err
	}

	alira.GetConnection().Create(&customer)
	custUser := &commerce.CustomerUser{
		CustomerID:    customer.Model.ID,
		UserID:        userProfile.ID,
		Username:      username,
		Email:         email,
		PrimaryMobile: mobile,
		Role:          "OWNER",
		Active:        true,
	}
	alira.GetConnection().Create(&custUser)

	mail := &service.Mail{
		From:     os.Getenv("SMTP_SENDER"),
		To:       []string{custUser.Email},
		Subject:  "[Alira] Welcome",
		Template: "views/mail/welcome.html",
		Data: map[interface{}]interface{}{
			"name":     fmt.Sprintf("%s %s", userProfile.FirstName, userProfile.LastName),
			"username": username,
		},
	}
	ms := &service.MailService{}
	_, err = ms.Send(mail)
	if err != nil {
		return nil, err
	}

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

		temp := make([]model.CustomerProfile, 0)
		for _, customer := range customers {
			req, err := http.NewRequest("GET",
				fmt.Sprintf("%s%s/%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_ACCOUNT"), customer.UserID),
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
			var userProfile messaging.UserProfile

			parser := &util.Parser{}
			_, err = parser.UnmarshalResponse(body, http.StatusOK, &userProfile)
			if err != nil {
				fmt.Println(err.Error())
			}

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
