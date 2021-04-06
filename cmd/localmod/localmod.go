package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"golang.org/x/mod/modfile"
)

var (
	commitArg = []string{"rev-list", "-1", "HEAD", "--abbrev-commit"}
)

func main() {
	mod := ""

	flag.StringVar(&mod, "mod", "go.mod", "go.mod file path")
	flag.Parse()

	readFile, err := ioutil.ReadFile(mod)
	if err != nil {
		panic(err)
	}

	file, err := modfile.Parse(mod, readFile, nil)
	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, v := range file.Replace {
		line := ""
		if v.New.Version != "" {
			line += v.New.Version
		} else {
			cmd := exec.Command("git", append([]string{"-C", v.New.Path}, commitArg...)...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				panic(err)
			}
			line += strings.TrimSpace(string(out))
		}
		line += "\t" + v.Old.Path + "\t" + v.New.Path
		_, _ = fmt.Fprintln(w, line)
	}
	_ = w.Flush()
}
