package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mohammadshabab/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mohammadshabab/bookstore_users-api/logger"
	"github.com/mohammadshabab/bookstore_utils-go/rest_errors"

	"github.com/mohammadshabab/bookstore_users-api/utils/mysql_utils"
)

const (
	//	indexUniqueEmail = "email_UNIQUE"
	errorNoRows                 = "no rows in result set"
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users Where id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

// var (
// 	userDB = make(map[int64]*User)
// )

func (user *User) Save() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	defer stmt.Close()
	//user.DateCreated = date_utils.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))

	}

	// 	if strings.Contains(err.Error(), indexUniqueEmail) {
	// 		return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
	// 	}
	// 	return errors.NewInternalServerError(
	// 		fmt.Sprintf("error when trying to save user: %s", err.Error()))

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}

	user.Id = userId
	return nil

	// currentUser := userDB[user.Id]
	// if currentUser != nil {
	// 	if currentUser.Email == user.Email {
	// 		return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", user.Email))
	// 	}
	// 	return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	// }

	// user.DateCreated = date_utils.GetNowString()
	// userDB[user.Id] = user

}
func (user *User) Get() rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	// if err := users_db.Client.Ping(); err != nil {
	// 	panic(err)
	// }
	// result := userDB[user.Id]
	// if result == nil {
	// 	return errors.NewNotFound(fmt.Sprintf("user %d not found", user.Id))
	// }
	// user.Id = result.Id
	// user.FirstName = result.FirstName
	// user.LastName = result.LastName
	// user.Email = result.Email
	// user.DateCreated = result.DateCreated
	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		//return errors.NewInternalServerError(err.Error())
		return rest_errors.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare Delete user statement", err)
		//return errors.NewInternalServerError(err.Error())
		return rest_errors.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to Delete user", err)
		//	mysql_utils.ParseError(err)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by Status statement", err)
		return nil, rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by Status", err)
		//return nil, errors.NewInternalServerError(err.Error())
		return nil, rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan user row into user struct", err)
			//return nil, mysql_utils.ParseError(err)
			return nil, rest_errors.NewInternalServerError("error when tying to gett user", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	return nil
}
