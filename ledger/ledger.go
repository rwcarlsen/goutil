
package ledger

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"encoding/csv"
	"time"
)

type Entry interface {
	Print(w io.Writer) error
}

type Journal []Entry

func (j Journal) WriteTo(w io.Writer) error {
	for _, e := range j {
		fmt.Fprint(w, "\n")
		e.Print(w)
	}
	return nil
}

type Transaction struct {
	Date string
	Description string
	Status string
	Posts []*Posting
}

func (t *Transaction) Print(w io.Writer) error {
	fmt.Fprintf(w, "%v %v %v\n", t.Date, t.Status, t.Description)
	for _, p := range t.Posts {
		p.Print(w)
	}
	return nil
}

func (t *Transaction) Post(p *Posting) {
	t.Posts = append(t.Posts, p)
}

type Posting struct {
	Account string
	Amount string
	Comment []string
}

func (p *Posting) Print(w io.Writer) error {
	if len(p.Amount) == 0 {
		fmt.Fprintf(w, "    %v\n", p.Account)
	} else {
		fmt.Fprintf(w, "    %v    $%v\n", p.Account, p.Amount)
	}

	for _, comm := range p.Comment {
		fmt.Fprintf(w, "        ; %v\n", comm)
	}
	return nil
}

func DecodeTabDelim(r io.Reader) (Journal, error) {
	// "Trans #"	"Date"	"Cleared"	"Name"	"Memo"	"Account"	"Debit"	"Credit"
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	journal := Journal{}
	currTrans := &Transaction{}
	for i, l := range lines {
		cells := strings.Split(l, "\t")
		if len(cells) < 8 {
			fmt.Printf("skipped line %v: %v\n", i, l)
			continue
		}

		date := ""
		if cells[1] != "" {
			date = cells[1][6:] + "/" +  cells[1][:5]
		}
		status := stripQuotes(cells[2])
		name := stripQuotes(cells[3])
		memo := stripQuotes(cells[4])
		account := stripQuotes(cells[5])
		amt := cells[6]
		if len(cells[6]) == 0 {
			amt = "-" + cells[7]
		}

		if len(cells[0]) > 0 {
			if currTrans.Date != "" {
				journal = append(journal, currTrans)
				currTrans = &Transaction{}
			}
			currTrans.Date = date
			currTrans.Status = status
			currTrans.Description = name
		}

		comms := []string{}
		if len(name) > 0 {
			comms = append(comms, fmt.Sprintf("who: %v", name))
		}
		if len(memo) > 0 {
			comms = append(comms, fmt.Sprintf("memo: %v", memo))
		}

		p := &Posting{
			Account: account,
			Amount: amt,
			Comment: comms,
		}
		currTrans.Post(p)
	}
	journal = append(journal, currTrans)
	return journal, nil
}

func stripQuotes(str string) string {
	if len(str) < 2 {
		return str
	} else if str[0] == '"' && str[len(str)-1] == '"' {
		txt := str[1:len(str)-1]
		return txt
	}
	return str
}

func DecodeCsv(r io.Reader, header bool) (Journal, error) {
	// "Account","Date","Check","Description","Amount"
	csvr := csv.NewReader(r)
	records, err := csvr.ReadAll()
	if err != nil {
		return nil, err
	}

	if header {
		records = records[1:]
	}

	journal := Journal{}
	for _, rec := range records {
		date, err := time.Parse("01/02/2006", rec[1])
		if err != nil {
			fmt.Printf("rec[1]=%v\n", rec[1])
			return nil, err
		}
		p := &Posting{Account: "??", Amount: rec[4]}

		trans := &Transaction{
			Date: date.Format("2006/01/02"),
			Description: rec[3],
		}
		trans.Post(p)

		journal = append(journal, trans)
	}
	return journal, nil
}

