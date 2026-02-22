package routes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"golang.hadleyso.com/msspr/src/auth"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
)

var appRelyingParty rp.RelyingParty

func authRoutes() {

	// App config
	SERVER_HOSTNAME := viper.GetString("SERVER_HOSTNAME")

	// OpenID Connect Client
	clientID := viper.GetString("CLIENT_ID")
	clientSecret := viper.GetString("CLIENT_SECRET")
	issuer := viper.GetString("OIDC_WELL_KNOWN")
	port := viper.GetString("OIDC_SERVER_PORT")
	scopes := strings.Split(viper.GetString("SCOPES"), " ")

	// OIDC URIs
	var redirectURI string
	if port != "" {
		redirectURI = fmt.Sprintf("%s:%v%v", SERVER_HOSTNAME, port, auth.CallbackPath)
	} else {
		redirectURI = fmt.Sprintf("%s%v", SERVER_HOSTNAME, auth.CallbackPath)
	}

	cookieHandler := httphelper.NewCookieHandler([]byte(viper.GetString("SESSION_KEY")), []byte(viper.GetString("SESSION_KEY")), httphelper.WithUnsecure())

	// Set Relying Party settings
	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
	}

	state := func() string {
		return uuid.New().String()
	}

	// OIDC RelyingParty Create
	type ctxKey struct{}
	logger := slog.Default()
	ctx := context.WithValue(context.Background(), ctxKey{}, logger)
	RelyingParty, idpInitErr := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, clientSecret, redirectURI, scopes, options...)
	if idpInitErr != nil {
		fmt.Printf("error creating provider %s", idpInitErr.Error()) // TODO: add logging
	}
	appRelyingParty = RelyingParty

	// Register routes
	Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if idpInitErr != nil {
			http.Error(w, "IdP Unavailable, try again later - service unavailable", http.StatusServiceUnavailable)
			return
		}
		rp.AuthURLHandler(state, appRelyingParty, rp.WithPromptURLParam("force"))(w, r)
	}).Methods("GET")
	Router.HandleFunc(auth.CallbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(auth.MarshallUserInfo), appRelyingParty)).Methods("GET", "POST")

	Router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if idpInitErr != nil {
			http.Error(w, "Failed to log out. IdP Unavailable, try again later - service unavailable", http.StatusServiceUnavailable)
			return
		}

		handleLogout(w, r)
	}).Methods("GET")

}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		})
	}

	http.Redirect(w, r, appRelyingParty.GetEndSessionEndpoint(), http.StatusSeeOther)

}
