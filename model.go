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
