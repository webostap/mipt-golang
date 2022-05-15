package pi

func CalculatePi(concurrent, iterations int, gen RandomPointGenerator) float64 {

}

type RandomPointGenerator interface {
	Next() (float64, float64)
}
