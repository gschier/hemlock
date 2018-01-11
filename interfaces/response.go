package interfaces

type Response interface {
	Cookie(name, value string) Response
	Status(status int) Response
	Data(data interface{}) Response
	View() View
}
