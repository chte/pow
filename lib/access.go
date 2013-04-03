package access

import "time"

const s_alpha = float64(0.1)

type Access struct {
	lastAccess time.Time
	ShortMean  int64
}

func (a *Access) New() {
	a.lastAccess = time.Now()
	a.ShortMean = 0
}
func (a *Access) Touch() {
	now := time.Now()
	a.ShortMean = int64(s_alpha*(float64(now.UnixNano()-a.lastAccess.UnixNano())) + (1-s_alpha)*float64(a.ShortMean))
	a.lastAccess = now
}
