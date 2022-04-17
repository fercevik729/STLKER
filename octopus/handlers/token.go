package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// TODO: make this more secure
const jwtKey = "secret"

type Credentials struct {
	gorm.Model
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Claims struct {
	Email string `json:"Email"`
	jwt.StandardClaims
}

// SignIn handles requests to /signin and creates JWTs for valid users
func (c *ControlHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	c.l.Println("[INFO] Handle Sign In")
	// Destruct incoming request payload
	var (
		creds   Credentials
		dbCreds Credentials
	)
	data.FromJSON(&creds, r.Body)

	// If the creds object is empty it was a bad request
	if creds == (Credentials{}) {
		c.LogHTTPError(w, "Must provide username and password", http.StatusBadRequest)
		return
	}
	// Connect to database
	db, err := NewGormDBConn("stlker.db")
	if err != nil {
		c.LogHTTPError(w, "couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Migrate schema
	db.AutoMigrate(&Credentials{})
	db.Model(&Credentials{}).Where("email=?", creds.Email).Find(&dbCreds)

	// If the passwords' hashes don't match inform the client
	hash, err := encrypt([]byte(creds.Password), jwtKey)
	if err != nil {
		c.LogHTTPError(w, "couldn't encrypt password", http.StatusInternalServerError)
		return
	}
	if hash != dbCreds.Password {
		c.LogHTTPError(w, "passwords do not match", http.StatusUnauthorized)
		return
	}
	// Set expiration time and claims
	expTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			// Set expiration time in unix
			ExpiresAt: expTime.Unix(),
		},
	}

	// Init JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		c.LogHTTPError(w, "couldn't create JWT", http.StatusInternalServerError)
		return
	}
	// Set HTTP cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expTime,
	})
}

// SignUp handles requests to /signup and adds new users to the db
func (c *ControlHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	c.l.Println("[INFO] Handle Sign Up")
	// Destruct incoming request payload
	var (
		creds Credentials
		err   error
	)
	data.FromJSON(&creds, r.Body)
	// If the credentials are empty don't sign them up
	if creds.Email == "" || creds.Password == "" {
		c.LogHTTPError(w, "email and password cannot be empty", http.StatusBadRequest)
		return
	}

	// Encrypt the password
	creds.Password, err = encrypt([]byte(creds.Password), jwtKey)
	if err != nil {
		c.LogHTTPError(w, "couldn't encrypt password", http.StatusInternalServerError)
	}

	// Add credentials to database
	db, err := NewGormDBConn("stlker.db")
	if err != nil {
		c.LogHTTPError(w, "couldn't connect to database", http.StatusInternalServerError)
		return
	}
	db.Model(&Credentials{}).Create(creds)

}

// Refresh handles requests to /refresh and regenerates tokens if the current token
// is within a minute of expiry
func (c *ControlHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	status, claims := ValidateJWT(r)
	if status != http.StatusOK {
		c.LogHTTPError(w, "bad refresh token request", status)
	}
	// Set new expiration time
	expTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expTime.Unix()
	// Create new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tknStr, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set HTTP cookie with token
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tknStr,
		Expires: expTime,
	})

}

// Encrypt encrypts a string using a 16 byte long key
func encrypt(weakText []byte, key string) (string, error) {

	// Check the length
	if len(key) < 16 {
		return "", errors.New("key is too short")
	}

	// Create new cipher
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Encrypt text
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, weakText, nil)), nil

}

func ValidateJWT(r *http.Request) (int, *Claims) {
	cookie, err := r.Cookie("token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return http.StatusUnauthorized, nil
		default:
			return http.StatusBadRequest, nil
		}

	}
	tknStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		switch err {
		case jwt.ErrSignatureInvalid:
			return http.StatusUnauthorized, nil
		default:
			return http.StatusBadRequest, nil
		}
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, nil
	}
	return http.StatusOK, claims

}

/*
// Decrypt decrypts a string using a 16 byte long key
func decrypt(strongText []byte, key string) (string, error) {

	// Check the length
	if len(key) < 16 {
		return "", errors.New("key is too short")
	}
	// Create cipher
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Decrypt the file contents
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(strongText) < nonceSize {
		return "", err
	}

	nonce, strongText := strongText[:nonceSize], strongText[nonceSize:]
	weakText, err := gcm.Open(nil, nonce, strongText, nil)
	if err != nil {
		return "", err
	}
	return string(weakText), nil
}
*/
