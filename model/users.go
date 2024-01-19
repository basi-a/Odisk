package model

import (
	g "odisk/global"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type Users struct {
	gorm.Model
	Name  		string	`json:"name"`
	Password	string	
	Email		string  `json:"email" gorm:"uniqueIndex"`
}

func InitUser()  {
	g.DB.AutoMigrate(Users{})
}
// add a user with name password email
func AddUser(name, password, email string)  error  {
	db := g.DB

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user := Users{
		Name: name,
		Password: string(hashedPassword),
		Email: email,
	}
	err = db.Create(&user).Error
	if err != nil {
    return err
	}
	return nil
}
// del a user by email  
func DelUser(email  string) error {
	db := g.DB

	err := db.Delete(&Users{}, email).Error
	if err != nil {
		return err
	}
	return nil
}

// update user by email with name password and email
func UpdateUser(name, password, email string) error {
	db := g.DB
	var err error
	if password != "" {
		password, err = hashPassword(password)
		if err != nil {
			return err
		}
	}
	user := Users{
		Name:     name,
		Password: password,
		Email:    email,
	}
	err = db.Model(&Users{}).Where("email = ?", email).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

// list all users
func ListUser() ([]Users, error) {
	db := g.DB

	var users []Users
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// get a user by email
func GetUser(email string) (*Users, error) {
	db := g.DB

	var user Users
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}