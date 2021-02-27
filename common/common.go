package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

// IsPathExist checks the path of file if is existed.
func IsPathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("IsPathExist %s is not existed: %s", path, err.Error())
			return false
		} else {
			fmt.Printf("IsPathExist os.Stat(%s) error: %s", path, err.Error())
			return false
		}
	} else {
		return true
	}
}

// MysqlPassword calculates the double hash of password by SHA1, equal to PASSWORD function in MySQL.
func MysqlPassword(password string) string {
	result := sha1.Sum([]byte(password))
	result = sha1.Sum(result[:])
	return hex.EncodeToString(result[:])
}
