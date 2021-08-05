package util

const (
	USD = "USD"
	IDR = "IDR"
	EUR = "EUR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, IDR, EUR:
		return true
	}
	return false
}
