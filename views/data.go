package views

import "../models"

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
