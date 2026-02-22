package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.hadleyso.com/msspr/src/config"
	"golang.hadleyso.com/msspr/src/models"
	"golang.hadleyso.com/msspr/src/scenes"
)

var (
	CallbackPath                             = "/auth/callback"
	SessionCookieStore *sessions.CookieStore = nil
)

// Generate HTTP error code and render login page to redirect
func UnauthorizedLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(scenes.TemplateFS, "scenes/login.html"))
	http.SetCookie(w, &http.Cookie{Name: "FLASH_PATH", Value: r.RequestURI, Path: "/", MaxAge: 300})
	w.WriteHeader(http.StatusUnauthorized)
	tmpl.Execute(w, nil)
}

// Sets cookie with user data after pulling from OIDC
func MarshallUserInfo(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty, info *oidc.UserInfo) {
	if SessionCookieStore == nil {
		SessionCookieStore = sessions.NewCookieStore([]byte(viper.GetString("SESSION_KEY")))
	}
	data, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user models.UserInfo

	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get unique Migadu
	var odicStringMap map[string]any
	json.Unmarshal(data, &odicStringMap)
	if val, ok := odicStringMap[config.C.MigaduAttribute]; ok {
		user.EmailMigadu = val.(string)
	} else {
		user.EmailMigadu = ""
	}

	var session_age = int(tokens.ExpiresIn)
	if viper.IsSet("SESSION_AGE") {
		session_age = viper.GetInt("SESSION_AGE")
	}

	session_SSPR_IDP_IDENTITY, _ := SessionCookieStore.Get(r, "SSPR_IDP_IDENTITY")
	session_SSPR_IDP_IDENTITY.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   session_age,
		HttpOnly: true,
	}

	session_SSPR_IDP_IDENTITY.Values["IDP"] = &user
	session_SSPR_IDP_IDENTITY.Save(r, w)

	session_SSPR_APP_AUTH, _ := SessionCookieStore.Get(r, "SSPR_APP_AUTH")
	session_SSPR_APP_AUTH.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   session_age,
		HttpOnly: true,
	}
	session_SSPR_APP_AUTH.Values["AUTHENTICATED"] = true
	session_SSPR_APP_AUTH.Save(r, w)

	FLASH_PATH, errCookie := r.Cookie("FLASH_PATH")
	if errCookie != nil {
		http.SetCookie(w, &http.Cookie{Name: "FLASH_PATH", Value: "", Path: "/", MaxAge: 0})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	http.SetCookie(w, &http.Cookie{Name: "FLASH_PATH", Value: "", Path: "/", MaxAge: 0})
	http.Redirect(w, r, FLASH_PATH.Value, http.StatusSeeOther)

}

// Returns user data from existing session
func GetUser(w http.ResponseWriter, r *http.Request) (*models.UserInfo, error) {
	if SessionCookieStore == nil {
		SessionCookieStore = sessions.NewCookieStore([]byte(viper.GetString("SESSION_KEY")))
	}

	// User data
	session_SSPR_IDP_IDENTITY, _ := SessionCookieStore.Get(r, "SSPR_IDP_IDENTITY")
	user, err := session_SSPR_IDP_IDENTITY.Values["IDP"].(*models.UserInfo)
	if !err {
		http.Error(w, "Error getting user from session", http.StatusInternalServerError)
		return user, fmt.Errorf("Error getting user from session")
	}

	return user, nil
}

// Check if request has valid user session
func ValidateSession(w http.ResponseWriter, r *http.Request) bool {
	if SessionCookieStore == nil {
		SessionCookieStore = sessions.NewCookieStore([]byte(viper.GetString("SESSION_KEY")))
	}

	session_SSPR_APP_AUTH, _ := SessionCookieStore.Get(r, "SSPR_APP_AUTH")

	if session_SSPR_APP_AUTH.Values["AUTHENTICATED"] != true {
		UnauthorizedLogin(w, r)
		return false
	}
	return true
}

// Check if request has valid session
func MiddleValidateSession(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if SessionCookieStore == nil {
			SessionCookieStore = sessions.NewCookieStore([]byte(viper.GetString("SESSION_KEY")))
		}

		session_APP_AUTH, _ := SessionCookieStore.Get(r, "SSPR_APP_AUTH")

		if session_APP_AUTH.Values["AUTHENTICATED"] != true {
			UnauthorizedLogin(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
