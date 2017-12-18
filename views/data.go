package views

import (
	"log"
	"net/http"
	"photoz/models"
	"time"
)

const (
	// AlertLvlError displayed when error encountered
	AlertLvlError = "danger"
	// AlertLvlWarning displays a warning
	AlertLvlWarning = "warning"
	// AlertLvlInfo displays informative message
	AlertLvlInfo = "info"
	// AlertLvlSuccess displays upon successful action
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is given when arbitrary error occurs on the backend
	AlertMsgGeneric = "Something went wrong. Please try again and contact us if the problem persists"
)

// Data holds all the data we'll pass into our view
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type PublicError interface {
	error
	Public() string
}

// Alert is used to render Bootstrap Alert messages
type Alert struct {
	Level   string
	Message string
}

// SetAlert determines what type of error to display to the end user
func (d *Data) SetAlert(err error) {
	// see if error incoming error implements PublicError interface.
	// if so, ok resolves to true. pErr is then instantiated and caste
	// as Public Error
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}
	return &alert
}

// RedirectAlert accepts all the normal params for an http.Redirect
// and performs a redirect, but only after persisting the the provided
// alert in a cookie so it can be displayed when the new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)
}
