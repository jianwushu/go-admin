package utils

import "strings"

// DesensitizePhone 手机号脱敏
// 13812345678 → 138****5678
func DesensitizePhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// DesensitizeEmail 邮箱脱敏
// test@example.com → t***@example.com
func DesensitizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	name := parts[0]
	if len(name) <= 1 {
		return email
	}
	return string(name[0]) + "***@" + parts[1]
}

// DesensitizeIDCard 身份证号脱敏
// 110101199001011234 → 110101********1234
func DesensitizeIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:6] + strings.Repeat("*", len(idCard)-10) + idCard[len(idCard)-4:]
}

// DesensitizeName 姓名脱敏
// 张三 → *三, 张三四 → *三四, 欧阳修 → **修
func DesensitizeName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	if len(runes) == 2 {
		return "*" + string(runes[1:])
	}
	return strings.Repeat("*", len(runes)-1) + string(runes[len(runes)-1:])
}

// DesensitizeBankCard 银行卡号脱敏
// 6222021234567890123 → 6222***********0123
func DesensitizeBankCard(cardNo string) string {
	if len(cardNo) < 8 {
		return cardNo
	}
	return cardNo[:4] + strings.Repeat("*", len(cardNo)-8) + cardNo[len(cardNo)-4:]
}

// DesensitizeAddress 地址脱敏（保留前6个字符）
func DesensitizeAddress(address string) string {
	runes := []rune(address)
	if len(runes) <= 6 {
		return address
	}
	return string(runes[:6]) + "******"
}
