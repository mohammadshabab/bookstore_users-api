package users

import (
	"fmt"

	"github.com/mohammadshabab/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mohammadshabab/bookstore_users-api/utils/errors"
	"github.com/mohammadshabab/bookstore_users-api/utils/mysql_utils"
)

const (
	//	indexUniqueEmail = "email_UNIQUE"
	errorNoRows           = "no rows in result set"
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users Where id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

// var (
// 	userDB = make(map[int64]*User)
// )

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	//user.DateCreated = date_utils.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)

	}

	// 	if strings.Contains(err.Error(), indexUniqueEmail) {
	// 		return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
	// 	}
	// 	return errors.NewInternalServerError(
	// 		fmt.Sprintf("error when trying to save user: %s", err.Error()))

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
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
func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
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

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFound(fmt.Sprintf("no users matching status %s, ", status))
	}
	return results, nil
}
