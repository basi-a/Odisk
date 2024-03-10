package model

import (
	// "log"
	"log"
	g "odisk/global"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type Users struct {
	gorm.Model
	UserName  	string	`json:"username"`
	Password	string	
	Email		string  `json:"email" gorm:"uniqueIndex"`
	Permission 	string	`json:"Permission" gorm:"default:general"`	//general/userAdmin/s3Admin
}

func InitUser()  {
	g.DB.AutoMigrate(Users{})
	type AdminPermission struct {
		userAdmin 	string
		s3Admin 	string
	}
	permission := AdminPermission{
		userAdmin: "userAdmin",
		s3Admin: 	"s3Admin",
	}
	user := new(Users)
	if err := user.AddUser(
		g.Config.Server.Admin.UserAdmin.Username, 
		g.Config.Server.Admin.UserAdmin.Password, 
		g.Config.Server.Admin.UserAdmin.Email,
		&permission.userAdmin); err != nil {
			log.Fatalln("User administrator creation failed:",err)
		}
	if err := user.AddUser(
		g.Config.Server.Admin.S3Admin.Username, 
		g.Config.Server.Admin.S3Admin.Password, 
		g.Config.Server.Admin.S3Admin.Email,
		&permission.s3Admin); err != nil {
			log.Fatalln("S3 administrator creation failed:",err)
		}
}

// add a user with name password email
func (users *Users )AddUser(username, password, email string, permission *string)  error  {
	db := g.DB

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := Users{
		UserName: username,
		Password: string(hashedPassword),
		Email: email,
	}
	if permission != nil {
		user.Permission = *permission
	}
	err = db.Create(&user).Error
	if err != nil {
    	return err	
	}
	return nil
}
// del a user by email  
func (users *Users )DelUser(email  string) error {
	db := g.DB

	err := db.Delete(&Users{}, email).Error
	if err != nil {
		return err
	}
	return nil
}

// update user by email with name password and email
func (users *Users )UpdateUser(username, password, email string) error {
	db := g.DB
	var err error
	if password != "" {
		password, err = hashPassword(password)
		if err != nil {
			return err
		}
	}
	user := Users{
		UserName: username,
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
func (users *Users )ListUser() ([]Users, error) {
	db := g.DB

	var usersList []Users
	err := db.Find(&usersList).Error
	if err != nil {
		return nil, err
	}
	return usersList, nil
}

// get a user by email
func (users *Users )GetUser(email string) (*Users, error) {
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

func (user *Users)VerifyAccount(email, password string) (ok bool, err error ) {
	db := g.DB
	ok = true
	if err = db.Select("password").Where("email=?", email).Find(&user).Error; err != nil {
		return !ok, err
	}else  if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return !ok, err
	}
	return ok, nil
}