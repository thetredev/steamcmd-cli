package shared

type SharedSocket interface {
	SendBytes([]byte)
	SendMessage(string)
}
