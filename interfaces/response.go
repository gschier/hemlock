package interfaces

type Response interface {
	Cookie(name, value string) Response
	Status(status int) Response
	Data(data interface{}) Response
	Dataf(format string, a ...interface{}) Response
	View() View
}
