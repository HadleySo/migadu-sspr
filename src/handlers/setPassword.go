package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"golang.hadleyso.com/msspr/src/auth"
	"golang.hadleyso.com/msspr/src/config"
	"golang.hadleyso.com/msspr/src/models"
)

func SetPasswd(w http.ResponseWriter, r *http.Request) {

	// Get inviter
	user, errUser := auth.GetUser(w, r)
	if errUser != nil {
		http.Redirect(w, r, "/500?error=GetUser+error", http.StatusSeeOther)
		return
	}

	var mailInfo *models.Mailbox

	// Check email address format
	emailAddress, err := mail.ParseAddress(user.EmailMigadu)
	if err != nil {
		mailInfo = nil
	} else {
		mailInfo = migaduGetMailbox(emailAddress.Address)
	}

	// Get new password
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	newPassword := r.FormValue("new_password")

	// Set
	setErr := migaduSetPassword(mailInfo.Address, newPassword)

	// Response handle
	if setErr != nil {
		http.SetCookie(w, &http.Cookie{Name: "FLASH_MESSAGE", Value: url.QueryEscape(setErr.Error()), Path: "/", MaxAge: 300})
		http.SetCookie(w, &http.Cookie{Name: "FLASH_ATTITUDE", Value: "BAD", Path: "/", MaxAge: 300})
	} else {
		http.SetCookie(w, &http.Cookie{Name: "FLASH_MESSAGE", Value: url.QueryEscape("Password set, allow 5 minutes to update"), Path: "/", MaxAge: 300})
		http.SetCookie(w, &http.Cookie{Name: "FLASH_ATTITUDE", Value: "GOOD", Path: "/", MaxAge: 300})
	}
	http.Redirect(w, r, "/my/", http.StatusFound)

}

func migaduSetPassword(email string, password string) error {
	local, domain, found := strings.Cut(email, "@")
	if !found {
		log.Println("migaduSetPassword() email format invalid " + email)
		return fmt.Errorf("Set password error - email format invalid")
	}

	// Password into JSON
	if len(password) < 5 {
		return fmt.Errorf("Password too short")
	}
	data := map[string]string{
		"password": password,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Set password error, no details could be provided")
	}

	// Create PUT
	uri := "https://api.migadu.com/v1/domains/" + domain + "/mailboxes/" + local
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.C.MigaduAPIuser, config.C.MigaduAPIkey)

	client := &http.Client{
		Timeout: 5 * time.Second, // 5‑second timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("migaduSetpassword resp.StatusCode: " + resp.Status)
		return fmt.Errorf("Unable to set password due to service provider error")
	}
	return nil
}
