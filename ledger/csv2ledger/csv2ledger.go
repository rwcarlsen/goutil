
package main

import (
	"flag"
	"os"
	"log"
	
	"github.com/rwcarlsen/goutil/ledger"
)

var errout = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags)

var inacc = flag.String("inacc", "Assets:ISU Checking", "account these csv entries belong in")
var outacc = flag.String("outacc", "Expenses:Consumables:Food", "account to balance these csv entries")
var header = flag.Bool("header", true, "indicates if the csv file(s) has a header row")

func main() {
	flag.Parse()

	files := flag.Args()

	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			errout.Print(err)
			continue
		}

		journal, err := ledger.DecodeCsv(f, *header)
		if err != nil {
			errout.Print(err)
			continue
		}

		for _, entry := range journal {
			trans := entry.(*ledger.Transaction)
			trans.Posts[0].Account = *inacc
			p := &ledger.Posting{Account: *outacc}
			trans.Post(p)
		}

		if err := journal.WriteTo(os.Stdout); err != nil {
			errout.Print(err)
		}
	}
}

