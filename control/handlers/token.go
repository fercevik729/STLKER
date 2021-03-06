package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var encryptKey string

type User struct {
	gorm.Model
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type Claims struct {
	Name  string `json:"Name"`
	Admin bool   `json:"Admin"`
	jwt.StandardClaims
}

// Initialize the encryptkey
func init() {
	encryptKey = ReadEnvVar("KEY")
	if encryptKey == "" {
		panic(errors.New("couldn't retrieve KEY env variable"))
	}
}

// SignUp handles requests to /signup and adds new users to the db
func (c *ControlHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Destruct incoming request payload
	var (
		creds     User
		otherUser User
		err       error
	)
	data.FromJSON(&creds, r.Body)
	// If the credentials are empty don't sign them up
	if creds.Username == "" || creds.Password == "" {
		c.logHTTPError(w, "email and password cannot be empty", http.StatusBadRequest)
		return
	}
	ok, msg := validateUser(creds)
	if !ok {
		c.logHTTPError(w, msg, http.StatusBadRequest)
		return
	}
	// Create database connection
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Check to see if there are other users with that username
	db.Where("username=?", creds.Username).First(&otherUser)
	if !reflect.DeepEqual(otherUser, User{}) {
		c.logHTTPError(w, "a user with that username already exists", http.StatusBadRequest)
		return
	}

	// Encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logHTTPError(w, "couldn't encrypt password", http.StatusInternalServerError)
		return
	}
	// Assign credential to this hash
	creds.Password = string(hash)

	// Add credentials to database
	db.Model(&User{}).Create(&creds)

	// Status code to indicate successfully created user
	w.WriteHeader(http.StatusCreated)
	c.l.Println("[INFO] Signed up user:", creds.Username)
	data.ToJSON(fmt.Sprintf("Happy Investing! %s", creds.Username), w)

}

// LogIn handles requests to /login and creates JWTs for valid users
func (c *ControlHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	// Destructure incoming request payload
	var (
		usr   User
		dbUsr User
		admin bool
	)
	data.FromJSON(&usr, r.Body)

	// If the user object is empty it was a bad request
	if usr == (User{}) {
		c.logHTTPError(w, "Must provide username and password", http.StatusBadRequest)
		return
	}
	// Connect to database
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Find credentials for the email from db
	db.Model(&User{}).Where("username=?", usr.Username).Find(&dbUsr)

	// Compare the hashes
	err = bcrypt.CompareHashAndPassword([]byte(dbUsr.Password), []byte(usr.Password))
	if err != nil {
		c.logHTTPError(w, "passwords do not match", http.StatusUnauthorized)
		return
	}
	// If the username is "admin" set admin privileges
	if usr.Username == "admin" {
		admin = true
	}
	c.l.Println("[INFO] Logged in user:", usr.Username)
	// Set expiration time and claims
	expTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Name:  usr.Username,
		Admin: admin,
		StandardClaims: jwt.StandardClaims{
			// Set expiration time in unix
			ExpiresAt: expTime.Unix(),
		},
	}

	// Init access JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Make sure key is a byte array
	tokenStr, err := token.SignedString([]byte(encryptKey))
	if err != nil {
		c.logHTTPError(w, "couldn't create JWT", http.StatusInternalServerError)
		return
	}
	// Init refresh JWT
	rfClaims := &Claims{
		Name:  usr.Username,
		Admin: admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Add(120 * time.Minute).Unix(),
		},
	}
	rfToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rfClaims)
	rfTokenStr, err := rfToken.SignedString([]byte(encryptKey))
	if err != nil {
		c.logHTTPError(w, "couldn't create refresh JWT", http.StatusInternalServerError)
		return
	}
	// Set HTTP cookies for both tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expTime,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    rfTokenStr,
		Expires:  expTime.Add(120 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	resp := make(map[string]string, 0)
	resp["User"] = dbUsr.Username
	resp["Access-Token"] = tokenStr
	resp["Refresh-Token"] = rfTokenStr

	data.ToJSON(resp, w)
}

// LogOut deletes the token cookie
func (c *ControlHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	// Set max age to < 0 for token and refresh token cookies in order to delete them
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "refreshToken",
		MaxAge: -1,
	})
	c.l.Println("[INFO] Logged out user:", retrieveUsername(r))

}

// Refresh handles requests to /refresh and regenerates tokens if the current token
// is within a minute of expiry
func (c *ControlHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	status, claims := ValidateJWT(r, "refreshToken")
	if status != http.StatusOK {
		c.logHTTPError(w, "bad refresh token request", status)
		return
	}
	// Set new expiration time
	expTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expTime.Unix()
	// Create new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Make sure key is a byte array
	tokenStr, err := token.SignedString([]byte(encryptKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set HTTP cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expTime,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	// Send the token in the resp body as well
	data.ToJSON(map[string]string{
		"Access-Token": tokenStr,
	}, w)

}

func (c *ControlHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "couldn't connect to database", http.StatusInternalServerError)
		return
	}
	user := retrieveUsername(r)
	var deletedUser User
	db.Model(User{}).Where("username=?", user).Delete(&deletedUser)
	c.l.Println("[INFO] Deleted User", user)
}

// ValidateJWT checks if the JWT token in the request token is valid and returns an http status
// code depending on if it is along with a pointer to a claim struct
func ValidateJWT(r *http.Request, tokenName string) (int, *Claims) {
	var tknStr string
	tknStr = r.Header.Get("Authorization")

	// If there was no token in the header check the cookies
	if tknStr == "" {
		cookie, err := r.Cookie(tokenName)
		switch err {
		case nil:
		case http.ErrNoCookie:
			return http.StatusUnauthorized, nil
		default:
			return http.StatusBadRequest, nil
		}
		tknStr = cookie.Value
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(encryptKey), nil
	})

	switch err {
	case nil:
	case jwt.ErrSignatureInvalid:
		return http.StatusUnauthorized, nil
	default:
		return http.StatusBadRequest, nil
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, nil
	}
	return http.StatusOK, claims

}
