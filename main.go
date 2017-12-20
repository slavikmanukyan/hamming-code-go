package main

import "fmt"
import "strconv"

type Value byte
type Matrix [][]Value

func Multiply(m1, m2 Matrix) (m3 Matrix, ok bool) {
	rows, cols, extra := len(m1), len(m2[0]), len(m2)
	if len(m1[0]) != extra {
		return nil, false
	}
	m3 = make(Matrix, rows)
	for i := 0; i < rows; i++ {
		m3[i] = make([]Value, cols)
		for j := 0; j < cols; j++ {
			for k := 0; k < extra; k++ {
				m3[i][j] ^= m1[i][k] & m2[k][j]
			}
		}
	}
	return m3, true
}

func (m Matrix) String() string {
	rows := len(m)
	cols := len(m[0])
	out := "["
	for r := 0; r < rows; r++ {
		if r > 0 {
			out += ",\n "
		}
		out += "[ "
		for c := 0; c < cols; c++ {
			if c > 0 {
				out += ", "
			}
			out += fmt.Sprintf("%d", m[r][c])
		}
		out += " ]"
	}
	out += "]"
	return out
}

var ENCODE_MATRIX = Matrix{
	[]Value{0, 1, 1, 1, 0, 0, 0},
	[]Value{1, 0, 1, 0, 1, 0, 0},
	[]Value{1, 1, 0, 0, 0, 1, 0},
	[]Value{1, 1, 1, 0, 0, 0, 1}}

var DECODE_MATRIX = Matrix{
	[]Value{1, 0, 0, 0, 1, 1, 1},
	[]Value{0, 1, 0, 1, 0, 1, 1},
	[]Value{0, 0, 1, 1, 1, 0, 1}}

var ERROR_TABLE = map[string]int{
	"100": 0,
	"010": 1,
	"001": 2,
	"011": 3,
	"101": 4,
	"110": 5,
	"111": 6}

func generateDataMatrix(data byte) (dataMatrix Matrix) {
	dataMatrix = Matrix{[]Value{}}
	dataRow := []Value{0, 0, 0, 0}
	binaryData := strconv.FormatInt(int64(data), 2)
	for index, _ := range binaryData {
		dataRow[int(index)+4-len(binaryData)] = Value(binaryData[index] - 48)
	}
	dataMatrix[0] = dataRow
	return
}

func generateDecodeDataMatrix(data string) (decodeDataMatrix Matrix) {
	decodeDataMatrix = Matrix{}
	for _, ch := range data {
		row := []Value{Value(ch - 48)}
		decodeDataMatrix = append(decodeDataMatrix, row)
	}
	return
}

func matrixRowToString(m Matrix) string {
	str := ""
	for _, n := range m[0] {
		str += string(n + 48)
	}
	return str
}

func encode(data byte) string {
	dataMatrix := generateDataMatrix(data)
	encodeMatrix, _ := Multiply(dataMatrix, ENCODE_MATRIX)
	return matrixRowToString(encodeMatrix)
}

func getResultFromCorrected(data string) byte {
	parityRemoved := data[3:7]
	result, _ := strconv.ParseInt(parityRemoved, 2, 0)
	return byte(result)
}

func checkIfCorrect(decodeResult Matrix) bool {
	return decodeResult[0][0] == 0 && decodeResult[1][0] == 0 && decodeResult[2][0] == 0
}

func getErrorPosition(decodedResult Matrix) int {
	syndrome := string(decodedResult[0][0]+48) + string(decodedResult[1][0]+48) + string(decodedResult[2][0]+48)
	return ERROR_TABLE[syndrome]
}

func flipPosition(data string, index int) string {
	corrected := "0"
	if data[index] == '0' {
		corrected = "1"
	}
	return data[:index] + corrected + data[index+1:]
}

func decode(data string) byte {
	decodeDataMatrix := generateDecodeDataMatrix(data)
	decodeResult, _ := Multiply(DECODE_MATRIX, decodeDataMatrix)
	if checkIfCorrect(decodeResult) {
		return getResultFromCorrected(data)
	}
	errorPosition := getErrorPosition(decodeResult)
	corrected := flipPosition(data, errorPosition)
	return getResultFromCorrected(corrected)
}

func main() {
	fmt.Println(encode(10))        // 1010 --> 1011010
	fmt.Println(encode(11))        // 1011 --> 0101011
	fmt.Println(decode("1011010")) // 1011010 --> 10 correct
	fmt.Println(decode("0101011")) // 0101011 --> 11 correct
	fmt.Println(decode("1011110")) // 1011110 --> 10 corrected error
	fmt.Println(decode("0100011")) // 0100011 --> 11 corrected error
}
