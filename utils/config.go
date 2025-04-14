package utils

import "sync"

var Mutex = &sync.Mutex{} // Protect OnlineUsers map