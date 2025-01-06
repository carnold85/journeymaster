package libs

func Average(xs []float32) float32 {
	if len(xs) == 0 {
		return 0.0
	}
	total := float32(0.0)
	for _, v := range xs {
		total += v
	}
	return total / float32(len(xs))
}
