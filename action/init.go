package action

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"crackmyd/common"
)

var ver = "v0.0.1"

// InitArgs initializes and resolves the input arguments.
func InitArgs() {
	verPtr := flag.Bool("version", false, usageMap["version"])
	flag.Usage = usage
	flag.Parse()

	if *verPtr {
		fmt.Printf("%s", ver)
		os.Exit(0)
	}

	obj := flag.Arg(0)
	if len(obj) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if !common.IsPathExist(obj) {
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(obj)
	if err != nil {
		fmt.Printf("Read file %s error: %s", obj, err.Error())
		os.Exit(2)
	}

	analyseFile(file)
}
