package interfaces

type Hub interface {
	Unreg(uid string)
	Broadcast(bytes []byte)
}
