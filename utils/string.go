package utils

import (
	"strconv"

	 _ "github.com/shopspring/decimal"
)

func Stoi(str string) int {

	i, _ := strconv.Atoi(str)

	return i
}

func Ternary(statement bool, a, b interface{}) interface{} {
    if statement {
        return a
    }
    return b
}