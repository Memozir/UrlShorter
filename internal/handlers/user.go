package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mail_system/internal/config"
	"mail_system/internal/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserJSON struct {
	Id         uint64 `json:"id"`
	Phone      string `json:"phone"`
	Login      string `json:"login"`
	Pass       string `json:"pass"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
}

func (handler *MailHandlers) RegisterUserHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("User registration handler")

	var userJSON UserJSON
	err := json.NewDecoder(r.Body).Decode(&userJSON)

	if err != nil {
		log.Printf("User decode error: %s", err)
	}

	log.Println(userJSON)
	contextCreateUser, cancelCreateUser := context.WithTimeout(r.Context(), time.Second*2)
	defer cancelCreateUser()

	userId := handler.Db.CreateUser(
		contextCreateUser,
		cancelCreateUser,
		userJSON.FirstName,
		userJSON.SecondName,
		userJSON.MiddleName,
		userJSON.Login,
		userJSON.Pass,
		userJSON.BirthDate)
	rw.Header().Set("Content-type", "application/json")

	fmt.Printf("User id: %d", userId.Val.(uint8))
}

func (handler *MailHandlers) GetUserHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := handler.Db.GetUserById(vars["id"])

	if user.Err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(user.Val.(model.User))
	rw.WriteHeader(http.StatusOK)
}

type UserAuthRequest struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type UserAuthResponse struct {
	Role int8 `json:"role"`
}

func (handler *MailHandlers) AuthUserHandler(rw http.ResponseWriter, r *http.Request) {
	var userAuth UserAuthRequest
	err := json.NewDecoder(r.Body).Decode(&userAuth)

	if err != nil {
		log.Println(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
	}

	//contextAuth, cancelAuth := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancelAuth()
	res := handler.Db.AuthUser(context.Background(), userAuth.Login, userAuth.Pass)

	if res.Err != nil {
		log.Println(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
	} else if res.Val.(model.UserAuth).ClientId > 0 {
		log.Println("SUCCESS CLIENT AUTH")

		response := UserAuthResponse{Role: config.UserRole}
		err = json.NewEncoder(rw).Encode(&response)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}

		rw.WriteHeader(http.StatusOK)
	} else if res.Val.(model.UserAuth).RoleCode != 0 {
		log.Println("SUCCESS CLIENT AUTH")

		response := UserAuthResponse{Role: res.Val.(model.UserAuth).RoleCode}
		err = json.NewEncoder(rw).Encode(response)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}
	} else {
		log.Println("ERROR AUTH")
		rw.WriteHeader(http.StatusBadRequest)
	}
}
