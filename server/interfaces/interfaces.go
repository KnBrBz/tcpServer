package interfaces

import "github.com/KnBrBz/message"

type Hub interface {
	Unreg(uid string)
	Broadcast(msg *message.M)
}

type Setup interface {
	Host() string
}
