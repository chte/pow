package problem

import (
	"../access"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"math/rand"
	"strconv"
)

type CpuInfo struct {
	Load, Avg float64
}

type Param struct {
	Local, Global *access.Access
	Cpu           CpuInfo
}

var BaseDifficulty = Difficulty{2, 16}
var ZeroDifficulty = Difficulty{0, 0}
var GetDifficulty = firstmodel

const cpu_thres = 40.0

type Difficulty struct {
	Zeroes   int
	Problems int
}

func (d *Difficulty) multiply(f int) *Difficulty {
	r := *d
	r.Problems *= f
	// log.Printf("Mulitplied %v by %v to get %v", *d, f, r)
	for r.Problems > 256 {
		r.Problems /= 256
		r.Zeroes++
	}
	log.Printf("Multiplied %v by %v to get %v", *d, f, r)
	return &r
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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func firstmodel(p Param) Difficulty {
	log.Printf("Base diff: %v", BaseDifficulty)
	if math.Max(p.Cpu.Load, p.Cpu.Avg) < cpu_thres {
		return ZeroDifficulty
	}
	if p.Local.ShortMean > max(p.Global.ShortMean, p.Global.LongMean) {
		if p.Local.ShortMean > 5*p.Global.ShortMean && math.Max(p.Cpu.Load, p.Cpu.Avg) < cpu_thres+0.2 {
			return ZeroDifficulty
		}
		return BaseDifficulty
	}
	// log.Printf("Multiplying")
	return *BaseDifficulty.multiply(1 + int((math.Max(p.Cpu.Avg, cpu_thres)-cpu_thres)*float64(int(5*p.Local.ShortMean/(p.Global.ShortMean+1)))))
}
func simpleonoff(p Param) Difficulty {
	// log.Printf("Base diff: %v", BaseDifficulty)
	// log.Printf("cpu load:%v", cpu_load)
	if p.Cpu.Load < cpu_thres {
		return ZeroDifficulty
	}
	return BaseDifficulty
}
func base(p Param) Difficulty {
	return BaseDifficulty
}
func zero(p Param) Difficulty {
	return ZeroDifficulty
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
