package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"log"

	"golang.hadleyso.com/msspr/src/auth"
	"golang.hadleyso.com/msspr/src/config"
	"golang.hadleyso.com/msspr/src/models"
	"golang.hadleyso.com/msspr/src/scenes"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {

	// Get inviter
	user, errUser := auth.GetUser(w, r)
	if errUser != nil {
		http.Redirect(w, r, "/500?error=GetUser+error", http.StatusSeeOther)
		return
	}

	var mailInfo *models.Mailbox

	emailAddress, err := mail.ParseAddress(user.EmailMigadu)
	if err != nil {
		mailInfo = nil
	} else {
		mailInfo = migaduGetMailbox(emailAddress.Address)
	}

	// Get flash message
	cMsg, errMsg := r.Cookie("FLASH_MESSAGE")
	cAtt, errAtt := r.Cookie("FLASH_ATTITUDE")
	var flash string
	var attitude string

	if errMsg == nil {
		flash, _ = url.QueryUnescape(cMsg.Value)
		http.SetCookie(w, &http.Cookie{Name: "FLASH_MESSAGE", Value: "", Path: "/", MaxAge: -1})
	}
	if errAtt == nil {
		attitude = cAtt.Value
		http.SetCookie(w, &http.Cookie{Name: "FLASH_ATTITUDE", Value: "", Path: "/", MaxAge: -1})
	}
	tmpl := template.Must(template.ParseFS(scenes.TemplateFS, "scenes/my-info.html", "scenes/base.html"))
	tmpl.ExecuteTemplate(w, "base",
		struct {
			Org           string
			User          *models.UserInfo
			MailboxInfo   *models.Mailbox
			PageTitle     string
			FlashMessage  string
			FlashAttitude string
		}{
			Org:           config.C.OrgName,
			PageTitle:     "SSPR",
			User:          user,
			MailboxInfo:   mailInfo,
			FlashMessage:  flash,
			FlashAttitude: attitude,
		},
	)

}

func migaduGetMailbox(email string) *models.Mailbox {
	local, domain, found := strings.Cut(email, "@")
	if !found {
		log.Println("migaduGetMailbox() email format invalid " + email)
		return nil
	}

	uri := "https://api.migadu.com/v1/domains/" + domain + "/mailboxes/" + local
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}

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
		log.Println("migaduGetMailbox() StatusCode " + strconv.Itoa(resp.StatusCode) + " " + uri)
		return nil
	}

	var mailbox models.Mailbox
	if err := json.NewDecoder(resp.Body).Decode(&mailbox); err != nil {
		log.Println("migaduGetMailbox() NewDecoder JSON error " + err.Error())
		return nil
	}

	return &mailbox
}
