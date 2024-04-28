package model

import (
	"log"
	g "odisk/global"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	UserName   string `json:"username"`
	Password   string
	Email      string `json:"email" gorm:"uniqueIndex"`
	Permission string `json:"permission" gorm:"default:general"` //general/userAdmin/s3Admin
}

func AutoMigrateUser() {
	if g.DB.Migrator().HasTable(&Users{}) {
		return
	}
	g.DB.AutoMigrate(&Users{})
	type AdminPermission struct {
		userAdmin string
		s3Admin   string
	}
	permission := AdminPermission{
		userAdmin: "userAdmin",
		s3Admin:   "s3Admin",
	}

	s3Admin := Users{
		UserName:   g.Config.Server.Admin.S3Admin.Username,
		Password:   g.Config.Server.Admin.S3Admin.Password,
		Email:      g.Config.Server.Admin.S3Admin.Email,
		Permission: permission.s3Admin,
	}
	userAdmin := Users{
		UserName:   g.Config.Server.Admin.UserAdmin.Username,
		Password:   g.Config.Server.Admin.UserAdmin.Password,
		Email:      g.Config.Server.Admin.UserAdmin.Email,
		Permission: permission.userAdmin,
	}

	if _, err := s3Admin.AddUser(); err != nil {
		log.Println("S3 administrator creation failed:", err)
	}
	if _, err := userAdmin.AddUser(); err != nil {
		log.Println("S3 administrator creation failed:", err)
	}
}

// add a user with name password email
func (user *Users) AddUser() (userID uint, err error) {

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword

	// 创建新用户
	if err := g.DB.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (user *Users) DelUser() error {
	return g.DB.Delete(&user).Error
}

// update user by id with name password and email
func (user *Users) Update() error {
	// log.Println(user)
	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	// log.Println(user)
	return g.DB.Updates(&user).Error
}

// list all users
func ListUser() ([]Users, error) {

	var usersList []Users
	err := g.DB.Find(&usersList).Error
	if err != nil {
		return nil, err
	}
	return usersList, nil
}

// get a user by email
func (user *Users) GetUserByEmail() error {
	return g.DB.Where("email = ?", user.Email).First(&user).Error
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (user *Users) VerifyAccount() (ok bool, err error) {

	ok = true
	userWithhashedpassword := Users{}
	if err = g.DB.Select("password").Where("email=?", user.Email).Find(&userWithhashedpassword).Error; err != nil {
		return !ok, err
	} else if err = bcrypt.CompareHashAndPassword([]byte(userWithhashedpassword.Password), []byte(user.Password)); err != nil {
		return !ok, err
	}
	return ok, nil
}
