package main

import (
	"flag"
	"github.com/mattkimber/gandalf/magica"
	"log"
	"strings"
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		obj, err := magica.FromFile(arg)
		if err != nil {
			log.Fatalf("could not read file: %s", err)
		}

		outFile := strings.Replace(arg, ".vox", "-out.vox", 1)
		err = obj.SaveToFile(outFile)
		if err != nil {
			log.Fatalf("could not write file: %s", err)
		}
	}
}