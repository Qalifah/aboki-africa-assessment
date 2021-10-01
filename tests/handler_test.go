package tests

import (
	"context"
	"encoding/json"
	"fmt"
	core "github.com/Qalifah/aboki-africa-assessment"
	"github.com/Qalifah/aboki-africa-assessment/handler"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	err := deleteAllFromTable("user_points")
	if !assert.NoError(t, err) {
		return
	}

	err = deleteAllFromTable("users")
	if !assert.NoError(t, err) {
		return
	}

	user1, err := seedOneUser("Quadri", "quadri@gmail.com")
	if !assert.NoError(t, err) {
		return
	}

	_, err = seedPointBalanceForUser(user1.ID, 0)
	if !assert.NoError(t, err) {
		return
	}

	tests := []struct {
		requestBody *handler.UserRequest
		wantCode    int
		checkData   bool
	}{
		{
			requestBody: &handler.UserRequest{
				Name:         "Qalifah",
				Email:        "qalifah@gmail.com",
				ReferralCode: nil,
			},
			checkData: true,
			wantCode:  http.StatusOK,
		},
		{
			requestBody: &handler.UserRequest{
				Name:         "Qali",
				Email:        "qali@gmail.com",
				ReferralCode: nil,
			},
			wantCode: http.StatusOK,
		},
		{
			requestBody: &handler.UserRequest{
				Name:         "Gee",
				Email:        "gee@gmail.com",
				ReferralCode: nil,
			},
			checkData: true,
			wantCode:  http.StatusOK,
		},
		{
			requestBody: &handler.UserRequest{
				Name:         "",
				Email:        "me@gmail.com",
				ReferralCode: nil,
			},
			checkData: false,
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		resp, err := registerUser(test.requestBody)
		if !assert.NoError(t, err) {
			return
		}

		if assert.Equal(t, test.wantCode, resp.StatusCode) {
			return
		}

		if test.checkData {
			body := &core.User{}
			err = getResponseBody(resp.Body, body)
			if !assert.NoError(t, err) {
				return
			}

			assert.NotEmpty(t, body.ID)
			assert.Equal(t, test.requestBody.Name, body.Name)
			assert.Equal(t, test.requestBody.Email, body.Email)
			assert.NotEmpty(t, body.CreatedAt)
			assert.NotEmpty(t, body.UpdatedAt)
			assert.Nil(t, body.DeletedAt)
		}
	}

	pp, err := testHandler.userPointRepository.GetPointsBalance(context.Background(), user1.ID)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 50, pp)
}

func TestTransaction(t *testing.T) {
	err := deleteAllFromTable("user_points")
	if !assert.NoError(t, err) {
		return
	}

	err = deleteAllFromTable("transactions")
	if !assert.NoError(t, err) {
		return
	}

	err = deleteAllFromTable("user_referrals")
	if !assert.NoError(t, err) {
		return
	}

	err = deleteAllFromTable("users")
	if !assert.NoError(t, err) {
		return
	}

	user1, err := seedOneUser("Qalifah", "qalifah@gmail.com")
	if !assert.NoError(t, err) {
		return
	}

	_, err = seedPointBalanceForUser(user1.ID, 4000)
	if !assert.NoError(t, err) {
		return
	}

	tests := []struct {
		user      *core.User
		userPoint *core.Point
		wantCode  int
	}{
		{
			user: &core.User{
				Name:         "qal",
				Email:        "qal@gmail.com",
			},
			userPoint: &core.Point{
				Points: 40000,
			},
			wantCode: http.StatusOK,
		},
		{
			user: &core.User{
				Name:         "quadri",
				Email:        "qdot@gmail.com",
			},
			userPoint: &core.Point{
				Points: 40000,
			},
			wantCode: http.StatusOK,
		},
		{
			user: &core.User{
				Name:         "Andy",
				Email:        "badboy@gmail.com",
			},
			userPoint: &core.Point{
				Points: 40000,
			},
			wantCode: http.StatusOK,
		},
	}

	ctx := context.Background()
	for _, test := range tests {
		err := testHandler.userRepository.CreateUser(ctx, test.user)
		if !assert.NoError(t, err) {
			return
		}

		test.userPoint.UserID = test.user.ID
		err = testHandler.userPointRepository.CreatePoint(ctx, test.userPoint)
		if !assert.NoError(t, err) {
			return
		}

		resp, err := transaction(&handler.TransferPointsRequest{
			SenderID:          test.user.ID,
			RecipientID: 	   user1.ID,
			Points:            210,
		})

		if !assert.NoError(t, err) {
			return
		}

		if assert.Equal(t, test.wantCode, resp.StatusCode) {
			return
		}
	}

	pp, err := testHandler.userPointRepository.GetPointsBalance(context.Background(), user1.ID)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 50, pp)
}

func getResponseBody(respBody io.ReadCloser, data interface{}) error {
	buf, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	fmt.Printf("response body: %+v", string(buf))
	if err = json.Unmarshal(buf, data); err != nil {
		return err
	}

	return nil
}

func registerUser(req *handler.UserRequest) (*http.Response, error) {
	return http.Post(url+"/register", "corelication/json", serialize(req))
}

func transaction(req *handler.TransferPointsRequest) (*http.Response, error) {
	return http.Post(url+"/transaction", "corelication/json", serialize(req))
}

func seedOneUser(name string, email string) (*core.User, error) {
	user := &core.User{
		Name:         name,
		Email:        email,
	}

	err := testHandler.userRepository.CreateUser(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedPointBalanceForUser(userID string, points int) (*core.Point, error) {
	p := &core.Point{
		UserID: userID,
		Points: points,
	}
	err := testHandler.userPointRepository.CreatePoint(context.Background(), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}