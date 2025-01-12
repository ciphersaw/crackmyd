package action

import (
	"fmt"
	"os"

	"crackmyd/common"

	flag "github.com/spf13/pflag"
)

var ver = "v0.0.2"

// InitArgs initializes and resolves the input arguments.
func InitArgs() {
	verPtr := flag.BoolP("version", "v", false, usageMap["version"])
	pwdPtr := flag.StringP("password", "p", "", usageMap["password"])
	sufPtr := flag.StringP("suffix", "s", "", usageMap["suffix"])
	flag.Usage = usage
	flag.Parse()

	if *verPtr {
		fmt.Printf("%s", ver)
		os.Exit(0)
	}

	if *pwdPtr != "" {
		if !common.IsPathExist(*pwdPtr) {
			os.Exit(1)
		}
		PwdMode = "assign"
		PwdFile = *pwdPtr
	}

	if *sufPtr != "" {
		if !common.IsPathExist(*sufPtr) {
			os.Exit(1)
		}
		SufMode = "assign"
		SufFile = *sufPtr
	}

	obj := flag.Arg(0)
	if len(obj) > 0 {
		if !common.IsPathExist(obj) {
			os.Exit(1)
		}
		result, err := analyseFile(obj)
		if err != nil {
			os.Exit(2)
		}
		printUserMYD(result)
	} else {
		flag.Usage()
	}
}
