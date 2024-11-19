package models

func ValidLuhn(orderNumber string) bool {
	var sum int
	alt := false

	for i := len(orderNumber) - 1; i >= 0; i-- {
		n := int(orderNumber[i] - '0')
		if n < 0 || n > 9 {
			return false
		}

		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}

	return sum%10 == 0
}
