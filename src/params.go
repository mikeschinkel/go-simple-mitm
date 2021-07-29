package simple_mitm

import (
	"flag"
	"fmt"
)

var PemPath string
var KeyPath string
var TestMode bool

func Initialize() {
	flag.StringVar(&PemPath,
		CertVar,
		CertFile,
		fmt.Sprintf("Filepath to %s", CertFile),
	)
	flag.StringVar(&KeyPath,
		KeyVar,
		KeyFile,
		fmt.Sprintf("Filepath to %s", KeyFile),
	)
	flag.BoolVar(&TestMode,
		ModeKey,
		false,
		"Is test mode?",
	)
}
