package utils

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// BcryptCheck 使用 bcrypt 校验密码
// password: 用户输入的明文密码
// hashedPassword: 数据库中存储的加密密码
func BcryptCheck(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
