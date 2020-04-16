package handlers

import (
	"encoding/json"
	"fmt"
	"go-webapp-starter/conf"
	HandlerContext "go-webapp-starter/context"
	Entities "go-webapp-starter/entities"
	Utils "go-webapp-starter/utils"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func accountExists(email string, hctx *HandlerContext.Context) bool {
	var response int
	var sql = fmt.Sprintf("SELECT EXISTS(SELECT account_id FROM account WHERE email = ?)")
	err := hctx.DB.Client.Get(&response, sql, email)
	Utils.CheckErr(err, "accountExists", "Failed to fetch SELECT EXISTS query")
	return response == 1
}

type loginBody struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"auto_login"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	var (
		err           error
		payload       = new(Entities.LoginPayload)
		account       = new(Entities.Account)
		accountExists bool
		sessionID     string
		cookie        http.Cookie
	)
	if !Utils.ParsePayload(w, r, payload) {
		return
	}
	accountExists = account.ReadFromDB(payload.Email, hctx)
	if !accountExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Email %s is not registered", payload.Email)
		return
	}
	if account.PasswordHash != Utils.HashPassword(payload.Password, account.Salt) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Incorrect password")
		return
	}
	_, err = hctx.DB.Client.Exec("UPDATE account SET last_login = ? WHERE account_id = ?", time.Now(), account.AccountID)
	Utils.CheckErr(err, "LoginHandler", "Failed to update last_login field of an account")
	sessionID = account.CreateSession(hctx)
	cookie = http.Cookie{
		Name:     conf.SessionCookie,
		Domain:   conf.SessionDomain,
		Path:     "/",
		HttpOnly: true,
		Value:    sessionID}
	if payload.RememberForTwoWeeks {
		cookie.Expires = time.Now().Add(time.Hour * 24 * 14)
	} else {
		cookie.Expires = time.Now().Add(time.Hour * 24)
	}
	http.SetCookie(w, &cookie)
	responseBytes, _ := json.Marshal(account.ToView(hctx))
	w.Write(responseBytes)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	var cookie = http.Cookie{
		Name:    conf.SessionCookie,
		Domain:  conf.SessionDomain,
		Path:    "/",
		Value:   "",
		Expires: time.Unix(0, 0)}
	http.SetCookie(w, &cookie)
}

func WhoamiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	var (
		err           error
		account       = new(Entities.Account)
		accountExists bool
	)
	accountExists, err = account.ReadSession(r, hctx)
	Utils.CheckErr(err, "WhoamiHandler", "Failed to read session")
	if !accountExists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	responseBytes, _ := json.Marshal(account.ToView(hctx))
	w.Write(responseBytes)
}

func ListAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	var (
		err           error
		user          = new(Entities.Account)
		accounts      []*Entities.AccountView
		responseBytes []byte
	)
	user.MustBe("superadmin", r, hctx)
	err = hctx.DB.Client.Select(&accounts, "SELECT account_id, admin_id, email, last_login, created_at, updated_at, deleted_at, full_name, roles FROM account_view")
	Utils.CheckErr(err, "ListAccountHandler", "Failed to fetch all accounts from database")
	for i, _ := range accounts {
		accounts[i].NormalizeTimes()
	}
	responseBytes, _ = json.Marshal(accounts)
	w.Write(responseBytes)
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	var (
		payload       = new(Entities.CreateAccountPayload)
		account       = new(Entities.Account)
		responseBytes []byte
	)
	if !Utils.ParsePayload(w, r, payload) {
		return
	}
	if accountExists(payload.Email, hctx) {
		w.WriteHeader(409)
		return
	}
	account.PopulateFromPayload(payload)
	account.Insert(hctx)
	responseBytes, _ = json.Marshal(account.ToView(hctx))
	w.Write(responseBytes)
}
