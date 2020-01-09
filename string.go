package tool
import (
	"sort"
	"strings"
	"github.com/Efabien/cognitive_types"
)

type setManipulation func( portion []string, from int, to int, start int)
type portionManipulation func( portion []string, from int, to int)

func Levenshtein(str1, str2 string) int {
	rune1 := []rune(str1)
	rune2 := []rune(str2)
	s1len := len(rune1)
	s2len := len(rune2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func Same(a, b string, degree int) bool {
	var ref int;
	if len(a) <= 4 || len(b) <= 4 {
		ref = 0;
	} else {
		ref = degree;
	}

	return Levenshtein(a, b) <= ref
}

func ExactMatch(a, b []string, degree int) bool {
	if len(a) != len(b) {
		return false
	}
	return Every(a, func(item string, index int) bool {
		return Same(item, b[index], degree)
	})
}

func Every(a []string, callback func(b string, index int)bool) bool {
	istrue := true
	for i := 0; i < len(a); i++ {
		istrue = callback(a[i], i)
		if istrue == false {
			break
		}
	}
	return istrue
}

func Some(a []string, callback func(b string, index int)bool) bool {
	isFalse := false
	for i := 0; i < len(a); i++ {
		isFalse = callback(a[i], i)
		if isFalse == true {
			break
		}
	}
	return isFalse
}

func Filter(input[]string, callback func(item string, index int)bool)(result[]string) {
	for index, current := range input {
		if callback(current, index) {
			result = append(result, current)
		}
	}
	return
}

func PortionReading(
	tab []string,
	interval int,
	callback portionManipulation) {
	for i := 0; i < len(tab) - interval + 1; i++ {
		callback(tab[i: i + interval], i, i + interval)
	}
}

func LongestSet(struc [][]string, callback setManipulation) {
	sort.SliceStable(struc, func( i, j int) bool { return len(struc[i]) > len(struc[j]) })
	for i := 0; i < len(struc) - 1; i++ {
		field := struc[i + 1:]
		for scope := len(struc[i + 1]); scope > 0; scope -- {
			PortionReading(struc[i], scope, func(portion []string, from int, to int) {
				for k := 0; k < len(field); k ++ {
					PortionReading(field[k], scope, func(actual []string, actualFrom int, actualTo int) {
						if ExactMatch(portion, actual, 2) {
							callback(portion, from, to, i)
						}
					})
				}
			})
		}
	}
}

func AjustSet(struc [][]string) {
	LongestSet(struc, func(portion []string, from int, to int, start int) {
		if from + to > len(portion) - 1 {
			struc[start] = struc[start][:from]
		} else {
			struc[start] = append(struc[start][:from], struc[start][from + to:]...)
		}
	})
}

func Arrayify(tab []string) (result [][]string) {
	for i := 0; i < len(tab); i ++ {
		result = append(result, strings.Fields(tab[i]))
	}
	return
}

func Precompute(intents cognitivetypes.Raw) cognitivetypes.Intents{
	result := make(cognitivetypes.Intents)
	for key, value := range intents {
		for _, content := range value {
			input := cognitivetypes.Intent { Texts: Arrayify(content), Treshold: 0.5 }
			result[key] = input
		}
	}
	return result
}