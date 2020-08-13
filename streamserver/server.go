package streamserver

import "sync"

var (
	streamserver *StreamServer
	streamserverlock sync.Mutex
)

func GetStreamServer() *StreamServer  {
	if streamserver == nil{
		streamserverlock.Lock()
		defer streamserverlock.Unlock()
		if streamserver == nil{
			streamserver = NewStreamServer()
		}
	}
	return streamserver
}

func StartStreamServer() error  {
	return GetStreamServer().StartServer()
}

func StopStreamserver()  {
	GetStreamServer().StopServer()
}
