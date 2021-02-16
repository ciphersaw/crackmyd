package action

import (
	"bytes"
	"os"
	"regexp"
	"strings"

	"crackmyd/common"
	"github.com/olekukonko/tablewriter"
)

// weakPasswordMap stores the hash of plaintext in weakPasswordList.
var weakPasswordMap = map[string]string{}

// mysqlPwdReg matches the information block containing host, user and password.
var mysqlPwdReg = regexp.MustCompile(`(.{3})[^\x20-\x7e]{1}([%\w\.\\]+)[^\x20-\x7e]{1}(\w+)\*(\w{40})`)

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

	// Check if the password is equal to the hash of user.
	if password == common.MysqlPassword(u.user) {
		u.plaintext = u.user
		return
	}

	// Check if the password is equal to the hash of weakPasswordList.
	for plaintext, hash := range weakPasswordMap {
		if password == hash {
			u.plaintext = plaintext
			return
		}
	}

	// Check if the password is equal to the hash of combination of user and suffix.
	for _, suffix := range userSuffixList {
		combo := u.user + suffix
		if password == common.MysqlPassword(combo) {
			u.plaintext = combo
			return
		}
	}
}

// analyseFile analyses the user.MYD file input, where in order to extract the host, user and password,
// and find out the plaintext of password finally.
func analyseFile(file []byte) {
	var result []userMYD

	records := mysqlPwdReg.FindAllSubmatch(file, -1)
	for _, r := range records {
		// If the 3 bytes before the control character of host do not contain /xFB, regard the record as invalid.
		if !bytes.Contains(r[1], []byte{251}) {
			continue
		}

		u := userMYD{
			host:     string(r[2]),
			user:     string(r[3]),
			password: string(r[4]),
		}
		u.crack()
		result = append(result, u)
	}

	printUserMYD(result)
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
