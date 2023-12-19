package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	db "mail_system/internal/db/postgres"
	"net/http"
	"time"
)

type EmployeeJSON struct {
	User         UserJSON `json:"user"`
	Role         RoleJSON `json:"role"`
	CreatorLogin string   `json:"login"`
}

func (emp EmployeeJSON) String() string {
	return fmt.Sprintf("UserId: %d, Role: %d", emp.User.Id, emp.Role.Code)
}

func (handler *MailHandlers) RegisterEmployeeHandler(rw http.ResponseWriter, r *http.Request) {
	var emp EmployeeJSON
	err := json.NewDecoder(r.Body).Decode(&emp)

	if err != nil {
		log.Printf("Registration employee error: %s", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		r.Context().Done()
	}

	contextCreateUser, cancelCreateUser := context.WithTimeout(r.Context(), time.Second*2)
	defer cancelCreateUser()

	contextGetRole, cancelGetRole := context.WithTimeout(r.Context(), time.Second*2)
	defer cancelGetRole()

	userCh := make(chan db.ResultDB)
	roleCh := make(chan db.ResultDB)
	creatorDepartmentCh := make(chan db.ResultDB)

	go func() {
		userCh <- handler.Db.CreateUser(
			contextCreateUser,
			cancelCreateUser,
			emp.User.FirstName,
			emp.User.SecondName,
			emp.User.Login,
			emp.User.Pass,
			emp.User.MiddleName,
			emp.User.BirthDate)
	}()

	go func() {
		roleCh <- handler.Db.GetRoleByName(contextGetRole, cancelGetRole, emp.Role.Name)
	}()

	go func() {
		creatorDepartmentCh <- handler.Db.GetEmployeeDepartment(r.Context(), emp.CreatorLogin)
	}()

	var user db.ResultDB
	var role db.ResultDB
	var departmentId db.ResultDB

	for i := 0; i < 3; i++ {
		select {
		case user = <-userCh:
			continue
		case role = <-roleCh:
			continue
		case departmentId = <-creatorDepartmentCh:
			continue

		}
	}

	contextCreateEmployee, cancelCreateEmployee := context.WithTimeout(r.Context(), time.Second*2)
	defer cancelCreateEmployee()

	employeeCreateResult := handler.Db.CreateEmployee(
		contextCreateEmployee,
		user.Val.(uint8),
		departmentId.Val.(uint64),
		role.Val.(uint8))

	if employeeCreateResult.Err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		r.Context().Done()
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}
