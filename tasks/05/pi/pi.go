package pi

import (
	"math"
	"sync"
)

func CalculatePi(concurrent, iterations int, gen RandomPointGenerator) float64 {
	acc := make(chan int, iterations)
	lim := make(chan int, iterations)
	var wg sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case lim <- 1:
					if x, y := gen.Next(); math.Sqrt(x*x+y*y) <= 1 {
						acc <- 1
					}
				default:
					return
				}
			}
		}()
	}
	wg.Wait()

	return float64(len(acc)*4) / float64(iterations)
}

type RandomPointGenerator interface {
	Next() (float64, float64)
}
