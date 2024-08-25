package adjust

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	sqlbuilder "user/internal/database/sql"
	"user/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Db *sql.DB
}


func (u *Database) LogIn(ctx context.Context, req *models.LogInRequest) (*models.LogInResponse, error) {
	query, args, err := sqlbuilder.LogIn(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var password string
	if err := u.Db.QueryRow(query, args...).Scan(&password); err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(password)
	fmt.Println(req.Password)
	return &models.LogInResponse{Status: u.ComparePassword(password, req.Password)}, nil
}

func (u *Database) CreateUser(ctx context.Context, req *models.RegisterUserRequest) (*models.GeneralResponse, error) {
	req.Password = u.Hashing(req.Password)
	query, args, err := sqlbuilder.Create(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var id int

	if err := u.Db.QueryRow(query, args...).Scan(&id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: fmt.Sprintf("user created with this id %v", id)}, nil
}

func (u *Database) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.GetUserResponse, error) {
	query, args, err := sqlbuilder.Get(req)
	if err != nil {
		log.Println(err)
		return nil, errors.New("there is no such user")
	}
	var res models.GetUserResponse
	if err := u.Db.QueryRow(query, args...).Scan(&res.ID, &res.Username, &res.Age, &res.Email, &res.LogOut); err != nil {
		log.Println(err)
		return nil, err
	}
	return &res, nil
}

func (u *Database) LastInserted(ctx context.Context, req *models.LastInsertedUser) (*models.GetUserResponse, error) {
	query, args, err := sqlbuilder.LastInserted()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var res models.GetUserResponse
	if err := u.Db.QueryRow(query, args...).Scan(&res.ID, &res.Username, &res.Age, &res.Email, &res.LogOut); err != nil {
		log.Println(err)
		return nil, err
	}
	return &res, nil
}

func (u *Database) UpdateUser(ctx context.Context, req *models.UpdateUserRequest) (*models.GeneralResponse, error) {
	req.Password = u.Hashing(req.Password)
	query, args, err := sqlbuilder.Update(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var id int

	if err := u.Db.QueryRow(query, args...).Scan(&id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: fmt.Sprintf("user updated with this id %v", id)}, nil
}

func (u *Database) LogOut(ctx context.Context, req *models.GetUserRequest) (*models.GeneralResponse, error) {
	query, args, err := sqlbuilder.LogOut(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: "Succesfully Logged Out"}, nil
}

func (u *Database) DeletUser(ctx context.Context, req *models.GetUserRequest) (*models.GeneralResponse, error) {
	query, args, err := sqlbuilder.LogOut(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: "Succesfully Deleted"}, nil
}


func (u *Database) ComparePassword(hashed, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *Database) Hashing(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

