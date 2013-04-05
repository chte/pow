package problem

import (
	"../access"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
)

var BaseDifficulty = Difficulty{2, 16}
var ZeroDifficulty = Difficulty{0, 0}

type Difficulty struct {
	Zeroes   int
	Problems int
}

type Problem struct {
	Seed, Solution int
}

func NewProblem() Problem {
	return Problem{Seed: int(rand.Uint32())}
}

func init_zeroes(s string) (num int) {
	for _, c := range s {
		if c != '0' {
			return
		}
		num++
	}
	return
}
func ConstructProblemSet(d Difficulty) []Problem {
	p := make([]Problem, d.Problems)
	for i, _ := range p {
		p[i] = NewProblem()
	}
	return p
}
func GetDifficulty(local, global *access.Access, cpu_load float64) Difficulty {
	if cpu_load < 0.4 {
		return ZeroDifficulty
	}
	return BaseDifficulty
}
func Verify(local, received []Problem, d Difficulty) bool {
	if len(received) < d.Problems {
		return false
	}
	ha := sha256.New()
	var sha string
	for i := 0; i < d.Problems; i++ {
		ha.Reset()
		ha.Write([]byte(strconv.Itoa(received[i].Solution)))
		ha.Write([]byte(strconv.Itoa(local[i].Seed)))
		sha = hex.EncodeToString(ha.Sum(nil))
		// fmt.Printf("Response solution: %v\n Calc Solution: %v\n", msg.Problems[i].Solution, sha)
		if init_zeroes(sha) < d.Zeroes {
			return false
		}
	}
	return true
}
