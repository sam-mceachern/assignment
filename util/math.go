package util

import (
	"fmt"
	"strconv"
	"strings"
)

func CheckNumberIsRoundedTo(number float32, decimalPlaces int) error {
	amountStr := strconv.FormatFloat(float64(number), 'f', -1, 32)
	amountSplit := strings.Split(amountStr, ".")
	if len(amountSplit) < 2 {
		return nil
	}

	if len(amountSplit[1]) > decimalPlaces {
		return fmt.Errorf("number has too many decimal places: %d", len(amountSplit[1]))
	}

	return nil
}
