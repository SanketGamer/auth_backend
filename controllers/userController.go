package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"user_management_system/helpers"
	"user_management_system/models"
	"user_management_system/notifier"
	"user_management_system/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//var repo=repository.NewUserRepository(database.DB)
var validate=validator.New()

//verify curr pass to db password/hash pass
func VerifyPassword(userPassword,HashPassword string)(bool,string){
	err:=bcrypt.CompareHashAndPassword([]byte(HashPassword),[]byte(userPassword))
	check:=true
	msg:=""
	if err!=nil{
	  msg=fmt.Sprintln("email or password is incorrect")
	  check=false
	}
	return check,msg
}

//hash the password
func HashPassword(password string)string{
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

//signup
func Signup(n *notifier.Notifier) gin.HandlerFunc{
	return func (c *gin.Context){

		var user models.User
		//parse the json
		err:=c.BindJSON(&user)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		//validate
		err=validate.Struct(user)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		//check email
		_,err=repository.FindByEmail(*user.Email)
		//already found
		if err==nil{
			c.JSON(400,gin.H{"error":"email already exist"})
			return
		}

		//password check
		if user.Password==nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"password required"})
			return
		}

		password:=HashPassword(*user.Password)
		user.Password=&password

		now:=time.Now()
		user.Created_at=now
		user.Updated_at=now

		user.ID=uuid.New()
		token, refresh_token, err := helpers.GenerateAllTokens(
           *user.Email,
           *user.First_Name,
           *user.Last_Name,
           user.ID.String(), 
           *user.User_type,  
		)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"check signup function"})
			return
		}

		user.Token=&token
		user.Refresh_token=&refresh_token

		err=repository.Create(&user)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not created"})
			return
		}

		n.Send(notifier.UserEvent{
			Type: notifier.EventUserCreated,
			UserID: user.ID.String(),
			Email: *user.Email,
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusOK,user)
	}
}

//login
func Login() gin.HandlerFunc{
		return func (c *gin.Context)  {
			ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
			defer cancel()
			var user models.User

			if err:=c.BindJSON(&user);err!=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
				return
			}

			founduser,err:=repository.FindByEmail((*user.Email))
			if err!=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
				return
			}

			passwordIdvalid,msg:=VerifyPassword(*user.Password,*founduser.Password)
			if !passwordIdvalid{
				c.JSON(http.StatusUnauthorized,gin.H{"error":msg})
				return
			}

		token, refreshToken, err := helpers.GenerateAllTokens(
	*founduser.Email,
	*founduser.First_Name,
	*founduser.Last_Name,
	founduser.ID.String(),   // uid
	*founduser.User_type,    // user_type
)
			if err!=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
	          return
			}

			helpers.UpdateAllToken(token,refreshToken,founduser.ID)

			updateUser,err:=repository.FindByUserID(ctx,founduser.ID)
			if err!=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated user"})
				return
			}
			c.JSON(http.StatusOK,updateUser)
		}
}

//show all user but only can see admin
func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		 log.Println("USER TYPE:", c.GetString("user_type"))

		//Authorization
		err:=helpers.CheckUserType(c,"ADMIN")
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":err.Error()})
			return
		}

		//Pagination
		recordPerPage,err:=strconv.Atoi(c.Query("recordPerPage"))
		if err!=nil || recordPerPage<1{
			recordPerPage=10;
		}
		page,err:=strconv.Atoi(c.Query("page"))
		if err!=nil || page<1{
			page=1
		}
		
		result,err:=repository.GetAllPaginated(ctx,page,recordPerPage)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured fetching users"})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

//show current user data
func Getuser() gin.HandlerFunc{
	return func (c *gin.Context){
		userIdstr:=c.Param("user_id")

		userId,err:=uuid.Parse(userIdstr)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid user id"})	
			return
		}

		//if user is admin allow admin can access all and if user_id maches allow his own data
		err=helpers.MatchUserTypeToUid(c,userIdstr)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		founduser,err:=repository.FindByUserID(ctx,userId)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,founduser)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userIdStr := c.Param("id")

		if userIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user id required",
			})
			return
		}

		// only ADMIN can delete users
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = repository.DeleteByUserID(ctx, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted successfully",
			"user_id": userId.String(),
		})
	}
}