package models

func ValidLuhn(orderNumber string) bool {
	sum := 0
	alt := false

	// Проходим с конца строки
	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit := int(orderNumber[i] - '0')
		if digit < 0 || digit > 9 {
			return false // Строка содержит нецифровой символ
		}
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}
	return sum%10 == 0
}
