package interfaces

import "tcpServer/message"

type Hub interface {
	Unreg(uid string)
	Broadcast(msg *message.M)
}
