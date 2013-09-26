package access

import (
	// "log"
	"time"
)

const s_alpha, l_alpha = float64(0.1), float64(0.01)

type Access struct {
	lastAccess time.Time
	ShortMean  int64
	LongMean   int64
}

func NewAccess() *Access {
	var a Access
	a.lastAccess = time.Now()
	a.ShortMean = 0
	return &a
}
func (a *Access) Touch(average int) int64 {
	now := time.Now()
	// log.Printf("Average: %v", int64(average))
	// if a.ShortMean == 0 {
	// 	// log.Println("Building default access time")
	// 	a.ShortMean = now.UnixNano() - a.lastAccess.UnixNano()
	// 	a.LongMean = a.ShortMean
	// 	// log.Printf("Now: %v, Last: %v, result: %v", now.UnixNano(), a.lastAccess.UnixNano(), a.ShortMean)
	// 	a.lastAccess = now
	// 	return 0
	// }
	diff := now.UnixNano() - a.lastAccess.UnixNano()
	a.AddTime(diff, average)
	// a.ShortMean = int64(s_alpha*(float64(int64(average)*diff)) + (1-s_alpha)*float64(a.ShortMean))
	// a.LongMean = int64(l_alpha*(float64(int64(average)*diff)) + (1-l_alpha)*float64(a.LongMean))

	a.lastAccess = now
	return diff
}

func (a *Access) Caress() {
	a.lastAccess = time.Now()
	return
}
func (a *Access) AddTime(diff int64, average int) {
	if a.ShortMean == 0 {
		// log.Println("Building default access time")
		a.ShortMean = diff
		a.LongMean = a.ShortMean
		// log.Printf("Now: %v, Last: %v, result: %v", now.UnixNano(), a.lastAccess.UnixNano(), a.ShortMean)
		a.Caress()
		return
	}
	a.ShortMean = int64(s_alpha*(float64(int64(average)*diff)) + (1-s_alpha)*float64(a.ShortMean))
	a.LongMean = int64(l_alpha*(float64(int64(average)*diff)) + (1-l_alpha)*float64(a.LongMean))
	a.Caress()
	return
}

// func (a *Access) Touch() {
// 	// now := time.Now()
// 	// a.ShortMean = int64(s_alpha*(float64(now.UnixNano()-a.lastAccess.UnixNano())) + (1-s_alpha)*float64(a.ShortMean))
// 	// a.lastAccess = now
// 	a.Touch(1)
// }
