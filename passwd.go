package utee

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswdHash generate password hash, compatible with PHP Yii framework
func PasswdHash(password string, cost ...int) (string, error) {
	// in Yii, the default cost is 13, but it's too slow, use the default value 10 in bcrypt
	realConst := bcrypt.DefaultCost
	if len(cost) > 0 {
		realConst = cost[0]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), realConst)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPasswd validate password，compatible with PHP Yii framework
func VerifyPasswd(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
