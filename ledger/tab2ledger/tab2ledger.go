
package main

import (
	"fmt"
	"os"
	"flag"
	"log"

	"github.com/rwcarlsen/goutil/ledger"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	out := flag.Arg(1)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	if err := ledger.Convert(f, w, &ledger.TabDelim{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success.")
}
