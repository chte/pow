package access

import (
	// "log"
	"time"
)

const s_alpha = float64(0.05)

type Access struct {
	lastAccess time.Time
	ShortMean  int64
}

func NewAccess() *Access {
	var a Access
	a.lastAccess = time.Now()
	a.ShortMean = 0
	return &a
}
func (a *Access) Touch(average int) {
	now := time.Now()
	// log.Printf("Average: %v", int64(average))
	if a.ShortMean == 0 {
		// log.Println("Building default access time")
		a.ShortMean = now.UnixNano() - a.lastAccess.UnixNano()
		// log.Printf("Now: %v, Last: %v, result: %v", now.UnixNano(), a.lastAccess.UnixNano(), a.ShortMean)
		a.lastAccess = now
		return
	}
	a.ShortMean = int64(s_alpha*(float64(int64(average)*(now.UnixNano()-a.lastAccess.UnixNano()))) + (1-s_alpha)*float64(a.ShortMean))
	a.lastAccess = now
}

// func (a *Access) Touch() {
// 	// now := time.Now()
// 	// a.ShortMean = int64(s_alpha*(float64(now.UnixNano()-a.lastAccess.UnixNano())) + (1-s_alpha)*float64(a.ShortMean))
// 	// a.lastAccess = now
// 	a.Touch(1)
// }
