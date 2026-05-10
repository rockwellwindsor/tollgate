package prompt

import (
	"fmt"
	"io"
)

type Response int

const (
	Yes     Response = iota // y — allow once
	No                      // n — deny
	AllTime                 // a — allow for this session
)

func Ask(w io.Writer, r io.Reader, binary string, args []string) (Response, error) {
	fmt.Fprintf(w, "tollgate: allow %s %v? [y/n/a] ", binary, args)
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return No, err
	}
	switch buf[0] {
	case 'y', 'Y':
		return Yes, nil
	case 'a', 'A':
		return AllTime, nil
	default:
		return No, nil
	}
}
