package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenApi struct {
	AccessToken string `json:"token"`
}

type AutenticacaoController interface {
	GetToken(c *gin.Context)
	//Signin(context *gin.Context) //ou GetToken, faz a mesma coisa porque para obter o token é necessário ter um usuário e senha
	Verify(c *gin.Context)
	CheckAuthorized(context *gin.Context) (*Claims, error)
}

type autenticacaoController struct {
	//TODO: definir o service que será utilizado
}

func NewAutenticacaoController() AutenticacaoController {
	return &autenticacaoController{}
}

/**
http://localhost:9001/signin

Content-Type: application/json
Body:
	{
	"username": "user1",
	"password": "password1"
	}
*/
func (c *autenticacaoController) GetToken(context *gin.Context) {
	tokenApi := signin(context.Writer, context.Request)
	context.IndentedJSON(http.StatusOK, tokenApi)
}

func (c *autenticacaoController) Verify(context *gin.Context) {
	fmt.Println("#Verify Test")
	claims, err := verify(context)
	if err != nil {
		fmt.Println("Erro: ", err)
		context.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	fmt.Printf("Usuario %v autorizado", claims.Username)
	context.IndentedJSON(http.StatusFound, gin.H{"message": "Login sucess"})
}

func (c *autenticacaoController) CheckAuthorized(context *gin.Context) (*Claims, error) {
	return verify(context)
}

/**
http://localhost:9001/signin

Content-Type: application/json
Body:
	{
	"username": "user1",
	"password": "password1"
	}
*/
func (c *autenticacaoController) Signin(context *gin.Context) {
	tokenApi := signin(context.Writer, context.Request)
	context.IndentedJSON(http.StatusOK, tokenApi)
}

func verify(context *gin.Context) (*Claims, error) {
	token := context.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		fmt.Println("Token: ", token)

		return tokenVerify(token)
	}
	return &Claims{}, fmt.Errorf("Unauthorized")
}

func tokenVerify(tknStr string) (*Claims, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, jwt.ErrSignatureInvalid
		}
		return claims, err
	}

	if !tkn.Valid {
		return claims, fmt.Errorf("Unauthorized")
	}
	//sucess
	return claims, nil
}

//======================================================================================================
// Testes autorization com JWT
//======================================================================================================
var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func signin(w http.ResponseWriter, r *http.Request) TokenApi {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return TokenApi{}
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return TokenApi{}
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return TokenApi{}
	}

	fmt.Println("TOKEN_GEN: ", tokenString)

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return TokenApi{AccessToken: tokenString}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
