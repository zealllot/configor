package configor

func runeToInt(r []rune) (result int) {
	for i := len(r) - 1; i >= 0; i-- {
		ten := 1
		for j := 0; j < len(r)-1-i; j++ {
			ten = ten * 10
		}
		result = result + int(r[i]-'0')*ten
	}
	return
}

func runeToFloat64(r []rune) (result float64) {
	var beforePoint bool
	var lenPoint int
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] == 46 {
			beforePoint = true
			result = result * 0.1
			continue
		}
		if beforePoint {
			ten := 1.0
			for j := 0; j < len(r)-1-(lenPoint+1)-i; j++ {
				ten = ten * 10
			}
			result = result + float64(r[i]-'0')*ten
		} else {
			lenPoint++
			result = result*0.1 + float64(r[i]-'0')
		}
	}
	return
}
