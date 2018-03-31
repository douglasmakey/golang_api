package models

import (
	"time"
	"encoding/hex"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/argon2"
	"github.com/labstack/gommon/log"
	"github.com/dgrijalva/jwt-go"

	"github.com/douglasmakey/backend_base/config"

)

type User struct {
	gorm.Model
	RoleID       uint   `gorm:"index;not null;default:'2'" json:"role_id,omitempty" valid:"int, required"`
	FirstName    string `gorm:"type:varchar(155);not null" json:"first_name,omitempty" valid:"required"`
	LastName     string `gorm:"type:varchar(155);not null" json:"last_name,omitempty" valid:"required"`
	Password     string `gorm:"type:varchar(128); not null" json:"password,omitempty" valid:"required"`
	Email        string `gorm:"type:varchar(100);unique_index" json:"email,omitempty" valid:"email,required"`
	RecoverToken string `gorm:"type:varchar(128); not null" json:"recover_token,omitempty" valid:"required"`
	Enabled      bool   `gorm:"default:'true'" json:"enabled,omitempty"`
}

type UserRegister struct {
	Email     string `json:"email" valid:"email,required"`
	Password1 string `json:"password1" valid:"required"`
	Password2 string `json:"password2" valid:"required"`
}

type UserLogged struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	RoleID    uint   `json:"role_id"`
	Jwt       string
}

// jwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	ID        uint   `json:"id"`
	RoleID    uint   `json:"role_id"`
	jwt.StandardClaims
}


func (u *User) SetPassword() {
	cfg := config.GetConfig()
	key  := argon2.Key([]byte(u.Password), cfg.Server.PasswordSalt, 3, 32*1024, 4, 32)
	u.Password = hex.EncodeToString(key)
}

func (u *User) generateUserJwt(origin *UserLogged) (error, string) {
	cfg := config.GetConfig()

	// Set custom claims
	claims := &JwtCustomClaims{
		origin.ID,
		origin.RoleID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(cfg.Server.JwtSecret)
	if err != nil {
		return err, ""
	}

	return nil, t
}

func (u *User) GenerateUserLogged() *UserLogged {

	userLog := new(UserLogged)
	userLog.ID = u.ID
	userLog.FirstName = u.FirstName
	userLog.LastName = u.LastName
	userLog.Email = u.Email
	userLog.RoleID = u.RoleID

	//Generate JWT
	err, jwt := u.generateUserJwt(userLog)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil
	}

	// Set JWT
	userLog.Jwt = jwt

	return userLog
}

