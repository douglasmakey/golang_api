package resources

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/gommon/log"


	"github.com/douglasmakey/backend_base/repositories"
	"github.com/douglasmakey/backend_base/models"
	"github.com/douglasmakey/backend_base/helpers"
)

type UserController struct{}

func (uc *UserController) Register(c echo.Context) error {
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return c.JSON(http.StatusNotFound, Error{
			Data: "Can't connect to database.",
		})
	}

	// Init UserRepository
	userRepo := repositories.NewUserRepo(db)

	// Init models
	user := new(models.User)
	userRegister := new(models.UserRegister)

	if err := c.Bind(userRegister); err != nil {
		log.Errorf("error : %v", err)
	}

	// Check if exists email
	exists := userRepo.Find(&user, "email", userRegister.Email)
	if exists {
		return c.JSON(http.StatusBadRequest, Error{
			Data: "Email already exists",
		})
	}

	//Validate User
	result, err := govalidator.ValidateStruct(userRegister)
	if err != nil {
		v := govalidator.ErrorsByField(err)
		return c.JSON(http.StatusBadRequest, Error{
			Data: v,
		})
	}

	//CheckPassword
	if userRegister.Password1 != userRegister.Password2 {
		return c.JSON(http.StatusBadRequest, Error{
			Data: "Passwords not match",
		})
	} else {
		user.Password = userRegister.Password1
		user.SetPassword()

	}

	//Set original model and values for defaults
	user.Email = userRegister.Email
	user.RecoverToken = helpers.GenerateTokenRecovery()

	save := userRepo.Save(&user)
	if !save {
		return c.JSON(http.StatusInternalServerError, Error{
			Data: "Contact to server. ",
		})
	}

	//Set userLogged for return
	userLogged := user.GenerateUserLogged()

	return c.JSON(http.StatusCreated, Response{
		Data:    userLogged,
		Success: result,
	})

}

func (uc *UserController) Login(c echo.Context) error {
	//Get DB
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return c.JSON(http.StatusNotFound, Error{
			Data: "Can't connect to database.",
		})
	}

	// Init UserRepository
	userRepo := repositories.NewUserRepo(db)
	user := new(models.User)

	//Bind form to Model.UserLogin
	if err := c.Bind(user); err != nil {
		log.Errorf("error: ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	// Validate params
	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, Error{
			Data: "Missing params",
		})
	}

	// Transform password
	user.SetPassword()

	success := userRepo.FindByCredentials(user)
	if !success {
		return c.JSON(http.StatusBadRequest, Error{
			Data: "Email or Password Invalid",
		})
	}

	userLogged := user.GenerateUserLogged()

	return c.JSON(http.StatusOK, Response{
		Data:    userLogged,
		Success: true,
	})

}
