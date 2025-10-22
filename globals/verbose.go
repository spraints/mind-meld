package globals

import "github.com/spf13/pflag"

var Verbose bool

func AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&Verbose, "verbose", false, "Enable more output")
}
