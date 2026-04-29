package helpers

import (
	"log"
	"os"
	"time"
	"user_management_system/database"
	"user_management_system/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type SignedDetails struct{
Email string
First_name string
Last_name string
Uid string
User_type string
jwt.StandardClaims
}

func init(){
	err:=godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading .env file")
	}
}

var db *gorm.DB
func InitDB(database *gorm.DB) {
	db = database
}

var secret_key=os.Getenv("SECRET_KEY")

//generate 2 tokens
func GenerateAllTokens(email, firstName, lastName, uid, userType string) (string, string, error) {
	//payload
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	//NewWithclaims generates token object means it sets header then asign payload(claims) and generates singnature(SignedString : data+secret key)  
	//header need - alg,type 
	//SigningMethodHS256 generates a algorithm
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret_key))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret_key))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

//find the user by uuid and update 3 fields
func UpdateAllToken(signedtoken,SignedRefreshToken string,userId uuid.UUID)error{
	 err := database.DB.Model(&models.User{}).
        Where("id = ?", userId).
        Updates(map[string]interface{}{
            "token":         signedtoken,
            "refresh_token": SignedRefreshToken,
            "updated_at":    time.Now(),
        }).Error

    return err
}

//token is valid or not
func ValidateToken(signedToken string) (claims *SignedDetails,msg string){
	token,err:=jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		//this function gives secret key to jwt
		func (token *jwt.Token) (interface{},error) {
			return []byte(secret_key),nil
		},
	)
	if err!=nil{
		msg=err.Error()
		return
	}
	claims,ok:=token.Claims.(*SignedDetails)
	if !ok{
		msg="the token is invalid"
		return
	}
	if claims.ExpiresAt<time.Now().Unix(){
		msg="token is expired"
		return
	}
	return claims,msg
}