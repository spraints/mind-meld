package git

import (
	"log"
	"os/exec"
)

var Verbose bool

func dbg(cmd *exec.Cmd) {
	if Verbose {
		log.Print(cmd)
		//log.Printf("%v %v", cmd.Path, cmd.Args)
	}
}
