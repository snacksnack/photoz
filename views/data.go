package views

const (
	// AlertLvlError displayed when error encountered
	AlertLvlError = "danger"
	// AlertLvlWarning displays a warning
	AlertLvlWarning = "warning"
	// AlertLvlInfo displays informative message
	AlertLvlInfo = "info"
	// AlertLvlSuccess displays upon successful action
	AlertLvlSuccess = "success"
)

// Data holds all the data we'll pass into our view
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert is used to render Bootstrap Alert messages
type Alert struct {
	Level   string
	Message string
}
