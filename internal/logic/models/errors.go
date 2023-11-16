package models

import "errors"

var (
	ErrCouldNotFindResult  = errors.New("could not find result")
	ErrNoExchangeRateFound = errors.New("could not find exchange rate")
)
