package entities

import (
	"fmt"
	"go-webapp-starter/conf"
	Conf "go-webapp-starter/conf"
	HandlerContext "go-webapp-starter/context"
	Utils "go-webapp-starter/utils"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Account struct {
	AccountID    string    `json:"account_id" db:"account_id"`
	AdminID      string    `json:"admin_id" db:"admin_id"`
	Email        string    `json:"email" db:"email"`
	Salt         string    `json:"salt" db:"salt"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" db:"deleted_at"`
	FullName     string    `json:"full_name" db:"full_name"`
	Roles        string    `json:"roles" db:"roles"`
}

func (i *Account) PopulateFromPayload(payload *CreateAccountPayload) {
	i.AccountID = Utils.RandString(8)
	i.AdminID = i.AccountID // self
	i.Email = payload.Email
	i.FullName = payload.FullName
	i.Salt = uuid.Must(uuid.NewV4()).String()
	i.PasswordHash = Utils.HashPassword(payload.Password, i.Salt)
}

func (i *Account) Insert(hctx *HandlerContext.Context) {
	_, err := hctx.DB.Client.NamedExec(`INSERT INTO account (account_id, admin_id, email, full_name, salt, password_hash)
		VALUES (:account_id, :admin_id, :email, :full_name, :salt, :password_hash)`, i)
	Utils.CheckErr(err, "CreateAccountHandler", "Failed to insert account to database")
	_, err = hctx.DB.Client.Exec(`INSERT INTO account_role (account_id, role) VALUES (?,?)`, i.AccountID, "user")
	Utils.CheckErr(err, "CreateAccountHandler", "Failed to create role for account")
}

func (i *Account) ReadFromDB(accountIDOrEmail string, hctx *HandlerContext.Context) (exists bool) {
	var (
		err error
		q   = "SELECT account_id, admin_id, email, salt, password_hash, last_login, created_at, updated_at, deleted_at, full_name, GROUP_CONCAT(role) AS roles FROM account JOIN account_role USING (account_id) WHERE"
	)
	if strings.Contains(accountIDOrEmail, "@") {
		err = hctx.DB.Client.Get(i, q+" email = ?", accountIDOrEmail)
	} else {
		err = hctx.DB.Client.Get(i, q+" account_id = ?", accountIDOrEmail)
	}
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (i *Account) CreateSession(hctx *HandlerContext.Context) (sessionID string) {
	// If same account logs in while already having a running session, clean the old session before starting a new one.
	oldSessionID, _ := hctx.RD.Client.HGet(conf.SessionCookie+":all", i.AccountID).Result()
	if oldSessionID != "" {
		hctx.RD.Client.Del(conf.SessionCookie + ":" + oldSessionID)
	}
	sessionID = uuid.Must(uuid.NewV4()).String()
	hctx.RD.Client.Set(conf.SessionCookie+":"+sessionID, i.AccountID, time.Second*time.Duration(conf.SessionTTL))
	// Store an index of all running sessions
	hctx.RD.Client.HSet(conf.SessionCookie+":all", i.AccountID, sessionID)
	return
}

func (i *Account) ReadSession(r *http.Request, hctx *HandlerContext.Context) (exists bool, err error) {
	cookie, err := r.Cookie(Conf.SessionCookie)
	if err != nil {
		return false, nil
	}
	sessionID := cookie.Value
	accountID, err := hctx.RD.Client.Get(conf.SessionCookie + ":" + sessionID).Result()
	if err != nil {
		if err.Error() == "redis: nil" { // no record
			return false, nil
		}
		return false, err // other Redis error
	}
	exists = i.ReadFromDB(accountID, hctx)
	return exists, nil
}

func (i *Account) ToView(hctx *HandlerContext.Context) *AccountView {
	accountView := new(AccountView)
	accountView.ReadFromDB(i.AccountID, hctx)
	return accountView
}

func (i *Account) HasRole(role string) bool {
	for _, r := range strings.Split(i.Roles, ",") {
		if r == role {
			return true
		}
	}
	return false
}

func (i *Account) MustBe(role string, r *http.Request, hctx *HandlerContext.Context) {
	exists, err := i.ReadSession(r, hctx)
	Utils.CheckErr(err, "MustBe", "Failed to read session")
	if !exists || !i.HasRole(role) {
		panic(http.StatusUnauthorized)
	}
}

func (i *Account) MustHaveAccess(role string, r *http.Request, hctx *HandlerContext.Context) {
	exists, err := i.ReadSession(r, hctx)
	Utils.CheckErr(err, "MustBe", "Failed to read session")
	if !exists {
		panic(http.StatusUnauthorized)
	}
	if i.HasRole(role) {
		return
	}
	panic(http.StatusUnauthorized)
}

type AccountView struct {
	AccountID string `json:"account_id" db:"account_id"`
	AdminID   string `json:"admin_id" db:"admin_id"`
	Email     string `json:"email" db:"email"`
	LastLogin string `json:"last_login" db:"last_login"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
	FullName  string `json:"full_name" db:"full_name"`
	Roles     string `json:"roles" db:"roles"`
}

func (i *AccountView) ReadFromDB(accountIDOrEmail string, hctx *HandlerContext.Context) (exists bool) {
	var (
		err error
		q   = "SELECT account_id, admin_id, email, last_login, created_at, updated_at, deleted_at, full_name, roles FROM account_view WHERE"
	)
	if strings.Contains(accountIDOrEmail, "@") {
		err = hctx.DB.Client.Get(i, q+" email = ?", accountIDOrEmail)
	} else {
		err = hctx.DB.Client.Get(i, q+" account_id = ?", accountIDOrEmail)
	}
	if err != nil {
		return false
	}
	i.NormalizeTimes()
	return true

}

func (i *AccountView) NormalizeTimes() {
	i.LastLogin = Utils.TimeLayoutChange(i.LastLogin)
	i.CreatedAt = Utils.TimeLayoutChange(i.CreatedAt)
	i.UpdatedAt = Utils.TimeLayoutChange(i.UpdatedAt)
	i.DeletedAt = Utils.TimeLayoutChange(i.DeletedAt)
}

type LoginPayload struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	RememberForTwoWeeks bool   `json:"remember_for_two_weeks"`
}

func (i *LoginPayload) Parse(r *http.Request) (validationErrors string, err error) {
	var (
		bodyBytes            []byte
		validationErrorsList []string
	)
	bodyBytes, err = Utils.GetBody(r, i)
	if err != nil {
		return "", err
	}
	// Email
	if !strings.Contains(string(bodyBytes), `"email":`) {
		validationErrorsList = append(validationErrorsList, "email:Email is required")
	} else {
		if isValid := Utils.IsEmailValid(i.Email); !isValid {
			validationErrorsList = append(validationErrorsList, fmt.Sprintf("email:\"%s\" is not a valid email address", i.Email))
		}
	}
	// Password
	if !strings.Contains(string(bodyBytes), `"password":`) {
		validationErrorsList = append(validationErrorsList, "password:Password is required")
	}
	return strings.Join(validationErrorsList, ","), nil
}

type CreateAccountPayload struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

func (i *CreateAccountPayload) Parse(r *http.Request) (validationErrors string, err error) {
	var (
		bodyBytes            []byte
		validationErrorsList []string
	)
	bodyBytes, err = Utils.GetBody(r, i)
	if err != nil {
		return "", err
	}
	// Email
	if !strings.Contains(string(bodyBytes), `"email":`) {
		validationErrorsList = append(validationErrorsList, "email:Email is required")
	} else {
		if isValid := Utils.IsEmailValid(i.Email); !isValid {
			validationErrorsList = append(validationErrorsList, fmt.Sprintf("email:\"%s\" is not a valid email address", i.Email))
		}
	}
	// Password
	if !strings.Contains(string(bodyBytes), `"password":`) {
		validationErrorsList = append(validationErrorsList, "password:Password is required")
	} else {
		if isValid := Utils.IsPasswordValid(i.Password); !isValid {
			validationErrorsList = append(validationErrorsList, fmt.Sprint("password:password must be at least 8 characters long, have lower and upper case letters and digits"))
		}
	}
	return strings.Join(validationErrorsList, ","), nil
}
