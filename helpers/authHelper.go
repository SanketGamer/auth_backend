package helpers


import (
	"errors"
	"github.com/gin-gonic/gin"
)

//check whether the current user has the required role
func CheckUserType(c *gin.Context,role string) error{
	userType:=c.GetString("user_type")
	if userType!=role{
		return errors.New("Unauthorize to access this resource")
	}
	return nil
}

//It ensures a normal USER can only access their own data
func MatchUserTypeToUid(c *gin.Context,userId string) (err error){
	userType:=c.GetString("user_type")
    uid:=c.GetString("uid")
	err=nil
	if userType=="USER" && userId!=uid{	
	err=errors.New("Unauthorized to access this resource")
	return err
	}
	err=CheckUserType(c,userType)
	return err
}