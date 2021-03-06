package action

import (
	"flag"
	"fmt"
	"os"

	"crackmyd/common"
)

var ver = "v0.0.2"

// InitArgs initializes and resolves the input arguments.
func InitArgs() {
	verPtr := flag.Bool("version", false, usageMap["version"])
	pwdPtr := flag.String("password", "", usageMap["password"])
	sufPtr := flag.String("suffix", "", usageMap["suffix"])
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
		analyseFile(obj)
	} else {
		flag.Usage()
	}
}
