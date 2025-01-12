package action

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"crackmyd/common"

	"github.com/olekukonko/tablewriter"
)

// PwdMode indicates the mode of method 2 in crack function.
var PwdMode = "default"

// PwdFile records the path of dictionary of passwords in crack function.
var PwdFile = ""

// SufMode indicates the mode of method 3 in crack function.
var SufMode = "default"

// SufFile records the path of dictionary of suffixes in crack function.
var SufFile = ""

// weakPasswordMap stores the hash of plaintext in weakPasswordList.
var weakPasswordMap = map[string]string{}

// mysqlPwdReg1 matches the information block containing host, user and password, used for MySQL v4.1~v5.6.
var mysqlPwdReg1 = regexp.MustCompile(`(.{3})[^\x20-\x7e]{1}([%\w\.\\]+)[^\x20-\x7e]{1}(\w+)\*(\w{40})`)

// mysqlPwdReg2 matches the information block containing host, user and password, used for MySQL v5.7+.
var mysqlPwdReg2 = regexp.MustCompile(`(.{3})[^\x20-\x7e]{1}([%\w\.\\]+)[^\x20-\x7e]{1}(\w+)[\x01]{30}[^\x20-\x7e]{1}mysql_native_password\)[\x00]{1}\*(\w{40})`)

func init() {
	// Calculate and store the hash of plaintext in weakPasswordList.
	for _, plaintext := range weakPasswordList {
		weakPasswordMap[plaintext] = common.MysqlPassword(plaintext)
	}
}

// userMYD describes the main elements in user.MYD file, as well as mysql.user table in MySQL,
// which includes host, user, password, and plaintext that could be enumerated by crack function.
type userMYD struct {
	host      string
	user      string
	password  string
	plaintext string
}

// crack enumerates the plaintext of password by brute force attack.
func (u *userMYD) crack() {
	password := strings.ToLower(u.password)

	// Strategy 1: Same As User Name
	// Check if the password is equal to the hash of user.
	if password == common.MysqlPassword(u.user) {
		u.plaintext = u.user
		return
	}

	// Strategy 2: Simple Guess
	// Check if the password is equal to the hash of the default weakPasswordList, or the user-defined passwords in dictionary.
	if PwdMode == "default" {
		for plaintext, hash := range weakPasswordMap {
			if password == hash {
				u.plaintext = plaintext
				return
			}
		}
	} else if PwdMode == "assign" {
		hit, plaintext := assignPasswordDict(PwdFile, password)
		if hit {
			u.plaintext = plaintext
			return
		}
	}

	// Strategy 3: Suffix Combo
	// Check if the password is equal to the hash of combination of user with the default suffixes, or the user-defined suffixes in dictionary.
	if SufMode == "default" {
		for _, suffix := range userSuffixList {
			combo := u.user + suffix
			if password == common.MysqlPassword(combo) {
				u.plaintext = combo
				return
			}
		}
	} else if SufMode == "assign" {
		hit, plaintext := assignSuffixDict(SufFile, password, u.user)
		if hit {
			u.plaintext = plaintext
			return
		}
	}

}

// assignPasswordDict enumerates the plaintext by a user-defined dictionary of passwords.
func assignPasswordDict(obj, password string) (hit bool, plaintext string) {
	file, err := os.Open(obj)
	if err != nil {
		fmt.Printf("assignPasswordDict os.Open(%s) error: %s\n", obj, err.Error())
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		plaintext = scanner.Text()
		if password == common.MysqlPassword(plaintext) {
			return true, plaintext
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("assignPasswordDict scanner.Scan(%s) error: %s\n", obj, err.Error())
		os.Exit(2)
	}

	return false, ""
}

// assignSuffixDict enumerates the plaintext by a user-defined dictionary of suffixes that will be combined with user.
func assignSuffixDict(obj, password, user string) (hit bool, combo string) {
	file, err := os.Open(obj)
	if err != nil {
		fmt.Printf("assignSuffixDict os.Open(%s) error: %s\n", obj, err.Error())
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		combo = user + scanner.Text()
		if password == common.MysqlPassword(combo) {
			return true, combo
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("assignSuffixDict scanner.Scan(%s) error: %s\n", obj, err.Error())
		os.Exit(2)
	}

	return false, ""
}

// analyseFile analyses the user.MYD file input, where in order to extract the host, user and password,
// and find out the plaintext of password finally.
func analyseFile(obj string) (result []userMYD, err error) {
	content, err := os.ReadFile(obj)
	if err != nil {
		fmt.Printf("analyseFile os.ReadFile(%s) error: %s", obj, err.Error())
		return nil, err
	}

	records := extractRecords(content)
	for _, r := range records {
		r.crack()
		result = append(result, r)
	}

	return result, nil
}

// extractRecords extracts records from the content of user.MYD, including host, user and password.
func extractRecords(content []byte) (records []userMYD) {
	matches := mysqlPwdReg1.FindAllSubmatch(content, -1)
	for _, m := range matches {
		// If the 3 bytes before the control character of host do not contain [0xFB], regard the record as invalid.
		if !bytes.Contains(m[1], []byte{0xFB}) {
			continue
		}
		u := userMYD{
			host:     string(m[2]),
			user:     string(m[3]),
			password: string(m[4]),
		}
		records = append(records, u)
	}

	matches = mysqlPwdReg2.FindAllSubmatch(content, -1)
	for _, m := range matches {
		// If the 3 bytes before the control character of host are not equal to [0xFF 0x13 0xFC], regard the record as invalid.
		if !bytes.Equal(m[1], []byte{0xFF, 0x13, 0xFC}) {
			continue
		}
		u := userMYD{
			host:     string(m[2]),
			user:     string(m[3]),
			password: string(m[4]),
		}
		records = append(records, u)
	}

	return
}

// printUserMYD outputs the result of analyseFile in table format.
func printUserMYD(uList []userMYD) {
	var data [][]string

	for _, u := range uList {
		uSlice := []string{u.host, u.user, u.password, u.plaintext}
		data = append(data, uSlice)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"host", "user", "password", "plaintext"})
	table.AppendBulk(data)
	table.Render()
}
