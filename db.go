package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type LoggedUser struct {
	userId      string
	loggedAt    time.Time
	issuedToken string
}

var secretKey = []byte("secret-key")
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"user": &memdb.TableSchema{
			Name: "user",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
				"password": &memdb.IndexSchema{
					Name:    "password",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Password"},
				},
			},
		},
		"loggedUser": &memdb.TableSchema{
			Name: "loggedUser",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "issuedToken"},
				},
			},
		},
	},
}
var db *memdb.MemDB

func InitDB() {
	var err error
	db, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
}

func AddUser(email string, password string, name string) {
	txn := db.Txn(true)
	defer txn.Abort()

	// encode the password
	hashedPassword := encodePassword(password)

	user := &User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: hex.EncodeToString(hashedPassword),
		Name:     name,
	}

	if err := txn.Insert("user", user); err != nil {
		panic(err)
	}

	txn.Commit()
}

func Login(email string, password string) (string, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	hashedPassword := encodePassword(password)

	raw, err := txn.First("user", "password", hex.EncodeToString(hashedPassword))
	if err != nil {
		return "", err
	}
	if raw == nil || raw.(*User).Email != email {
		return "", errors.New("invalid credentials")
	}

	user := raw.(*User)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"email":  user.Email,
		"name":   user.Name,
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	loggedUser := &LoggedUser{
		userId:      user.ID,
		loggedAt:    time.Now(),
		issuedToken: tokenString,
	}

	txn = db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("loggedUser", loggedUser); err != nil {
		return "", err
	}

	txn.Commit()

	return tokenString, nil

}

func GetLoggedUserByToken(token string) (*User, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("loggedUser", "id", token)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, errors.New("invalid token")
	}

	rawUser, err := txn.First("user", "id", raw.(*LoggedUser).userId)
	if err != nil {
		return nil, err
	}

	return rawUser.(*User), err
}

func encodePassword(password string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}
