package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"monkiato/guardian/internal/auth"
	"monkiato/guardian/internal/models"
	"monkiato/guardian/internal/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//DefaultTokenExpirationHours token duration
const DefaultTokenExpirationHours int64 = 2

//DefaultRedirectURL redirect url used by default
const DefaultRedirectURL string = "http://localhost/"

//DefaultDomainName domain name
const DefaultDomainName string = "localhost"

//Auth objects required by the auth routes
type Auth struct {
	db                   *gorm.DB
	logger               *log.Logger
	authHandler          *auth.Handler
	tokenExpirationHours int64
}

//NewAuth auth constructor
func NewAuth(_db *gorm.DB, _logger *log.Logger) *Auth {
	secretKey, found := os.LookupEnv("SECRET_KEY")
	if !found {
		log.Fatal("SECRET_KEY not found")
	}
	tokenExpirationHours := DefaultTokenExpirationHours
	hours, found := os.LookupEnv("TOKEN_EXPIRATION_HOURS")
	if found {
		tokenExpirationHours, _ = strconv.ParseInt(hours, 10, 64)
	}

	authHandler := auth.NewHandler(getDomainName(), secretKey)
	return &Auth{
		db:                   _db,
		logger:               _logger,
		authHandler:          authHandler,
		tokenExpirationHours: tokenExpirationHours,
	}
}

//AddRoutes add all routers related to the users
func (a *Auth) AddRoutes(router *mux.Router) {
	a.logger.Print("adding auth routes...")
	router.HandleFunc("/login", a.authenticateHandler).Methods(http.MethodPost)
	router.HandleFunc("/logout", a.logoutHandler).Methods(http.MethodGet)
	router.HandleFunc("/validate", a.validateHandler)
	router.HandleFunc("/signin", a.signinHandler).Methods(http.MethodPost)
	router.HandleFunc("/approve", a.approveHandler).Methods(http.MethodPost)
	router.HandleFunc("/me", a.getMeHandler).Methods(http.MethodGet)
}

func (a *Auth) getTokenExpirationTime() time.Duration {
	return time.Hour * time.Duration(a.tokenExpirationHours)
}

func (a *Auth) isValidUser(user models.User) bool {
	return len(user.Username) > 0 &&
		len(user.Email) > 0 &&
		len(user.Name) > 0 &&
		len(user.Lastname) > 0 &&
		len(user.Password) > 0
}

func (a *Auth) signinHandler(w http.ResponseWriter, req *http.Request) {
	var err error

	user := models.User{
		Username:      req.FormValue("username"),
		Email:         req.FormValue("email"),
		Name:          req.FormValue("name"),
		Lastname:      req.FormValue("lastname"),
		Password:      req.FormValue("password"),
		ApprovalToken: rand.String(10),
	}

	// validate user data coming from client
	if !a.isValidUser(user) {
		fmt.Println("invalid user data received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create cookie with its JWT token internally
	expiresAt := time.Now().Add(a.getTokenExpirationTime()).Unix()
	sessionUser := models.SessionUser{user.Username, user.Name, user.Lastname}
	cookie, token, err := a.authHandler.CreateCookie(&sessionUser, expiresAt)
	if err != nil {
		fmt.Println("unable to create the cookie", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Save user in DB
	user.Token = token
	err = models.CreateUser(a.db, &user)
	if err != nil {
		fmt.Println("unable to create user data", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("created user with token: %s, approval token is: %s", token, user.ApprovalToken)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}

func (a *Auth) authenticateHandler(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	loginUser := models.LoginUser{
		Username: username,
		Password: password,
	}

	// Check user exists
	dbUser, err := models.GetUser(a.db, loginUser.Username)
	if err != nil {
		fmt.Println("User not found", loginUser.Username, err)
		a.redirect(w, req, err)
		return
	}

	// check if password is matching
	if !dbUser.Approved || dbUser.Password != password {
		fmt.Println("Not Approved or password doesn't match", loginUser.Password, dbUser.Password)
		a.redirect(w, req, errors.New("Password doesn't match"))
		return
	}

	// create cookie with its JWT token internally
	expiresAt := time.Now().Add(a.getTokenExpirationTime()).Unix()
	sessionUser := models.SessionUser{dbUser.Username, dbUser.Name, dbUser.Lastname}
	cookie, token, err := a.authHandler.CreateCookie(&sessionUser, expiresAt)
	if err != nil {
		fmt.Println("unable to create the cookie", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update token in user db record
	dbUser.Token = token
	err = models.UpdateUser(a.db, dbUser)
	if err != nil {
		fmt.Println("can't update user", err)
		a.redirect(w, req, errors.New("can't update user"))
		return
	}

	http.SetCookie(w, cookie)
	// check if there's a redirect to send it to that url directly
	if redirect, ok := req.URL.Query()["redirect"]; ok && len(redirect) > 0 {
		http.Redirect(w, req, redirect[0], 301)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (a *Auth) logoutHandler(w http.ResponseWriter, req *http.Request) {
	// get cookie
	cookie, _, err := a.authHandler.ReadCookie(req)
	if err != nil {
		fmt.Println("cookie not found")
		w.WriteHeader(http.StatusOK)
		return
	}

	// invalidate cookie
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusNoContent)
}

func (a *Auth) approveHandler(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	approvalToken := req.FormValue("approval-token")
	loginUser := models.LoginUser{
		Username: username,
	}

	// Check user exists
	dbUser, err := models.GetUser(a.db, loginUser.Username)
	if err != nil {
		fmt.Println("User not found", loginUser.Username, err)
		a.redirect(w, req, err)
		return
	}

	// check if approval token is matching
	if dbUser.ApprovalToken != approvalToken {
		fmt.Println("Approval token doesn't match", approvalToken, dbUser.ApprovalToken)
		a.redirect(w, req, errors.New("Approval token doesn't match"))
		return
	}

	// mark user as approved
	dbUser.Approved = true
	err = models.UpdateUser(a.db, dbUser)
	if err != nil {
		fmt.Println("can't update user", err)
		a.redirect(w, req, errors.New("can't update user"))
		return
	}

	// check if there's a redirect to send it to that url directly
	if redirect, ok := req.URL.Query()["redirect"]; ok && len(redirect) > 0 {
		http.Redirect(w, req, redirect[0], 301)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (a *Auth) redirect(w http.ResponseWriter, req *http.Request, err error) {
	//fmt.Println(err.Error())
	host := req.Header.Get("X-Forwarded-Host")
	proto := req.Header.Get("X-Forwarded-Proto")
	uri := req.Header.Get("X-Forwarded-Uri")
	redirectURL, found := os.LookupEnv("REDIRECT_URL")
	if !found {
		redirectURL = DefaultRedirectURL
	}
	http.Redirect(w, req, fmt.Sprintf("%s?redirect=%s://%s%s", redirectURL, proto, host, uri), 301)
}

func (a *Auth) validateHandler(w http.ResponseWriter, req *http.Request) {
	// get cookie & token
	_, jwtToken, err := a.authHandler.ReadCookie(req)
	if err != nil {
		fmt.Println("cookie not found")
		a.redirect(w, req, err)
		return
	}

	if jwtToken.Token == nil || !jwtToken.Token.Valid {
		fmt.Println("invalid token")
		a.redirect(w, req, err)
		return
	}

	sessionUser, err := jwtToken.GetUser()
	if err != nil {
		fmt.Println("unable to get sessionUser for given token")
		a.redirect(w, req, err)
		return
	}

	//get user data from db to confirm the token was assigned to this user
	dbUser, err := models.GetUser(a.db, sessionUser.Username)
	if err != nil {
		fmt.Println("unable to get user from DB, username:", sessionUser.Username)
		a.redirect(w, req, err)
		return
	}

	token, err := jwtToken.ToString()
	if err != nil {
		fmt.Println("unable to obtain token as a string")
		a.redirect(w, req, err)
		return
	}

	if !dbUser.Approved || dbUser.Token != token {
		fmt.Println("User not approved or mismatched token")
		a.redirect(w, req, errors.New("User not approved or mismatched token"))
		return
	}

	//all good
	w.Header().Set("X-Forwarded-User", dbUser.Username)
	w.WriteHeader(http.StatusOK)
}

func (a *Auth) getMeHandler(w http.ResponseWriter, req *http.Request) {
	// get cookie & token
	_, jwtToken, err := a.authHandler.ReadCookie(req)
	if err != nil {
		fmt.Println("cookie not found")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if jwtToken.Token == nil || !jwtToken.Token.Valid {
		fmt.Println("invalid token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	sessionUser, err := jwtToken.GetUser()
	if err != nil {
		fmt.Println("unable to get sessionUser for given token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//get user data from db to confirm the token was assigned to this user
	dbUser, err := models.GetUser(a.db, sessionUser.Username)
	if err != nil {
		fmt.Println("unable to get user from DB, username:", sessionUser.Username)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	data, _ := json.Marshal(dbUser)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getDomainName() string {
	domainName, found := os.LookupEnv("DOMAIN_NAME")
	if !found {
		domainName = DefaultDomainName
	}
	return domainName
}
