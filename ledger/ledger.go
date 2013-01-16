
package ledger

import (
	"io"
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
	Posts []*Post
}

func (t *Transaction) Print(w io.Writer) error {
	fmt.Fprintf(w, "%v %v %v\n", t.Date, t.Status, t.Description)
	for _, p := range t.Posts {
		p.Print(w)
	}
	return nil
}

type Post struct {
	Account string
	Amount string
	Comment string
}

func (p *Post) Print(w io.Writer) error {
	fmt.Fprintf(w, "    %v    %v", p.Account, p.Amount)
	if Comment != "" {
		fmt.Fprintf(w, "  ; %v", p.Comment)
	}
	fmt.Fprint(w, "\n")
	return nil
}

