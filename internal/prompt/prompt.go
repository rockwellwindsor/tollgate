package prompt

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Response int

const (
	Yes     Response = iota // y — allow once
	No                      // n — deny
	AllTime                 // a — allow for this session
)

// dialogFn is the fallback when tty read returns EOF (non-interactive context).
// Replaced in tests to avoid spawning a real GUI dialog.
var dialogFn = osascriptDialog

func Ask(w io.Writer, r io.Reader, binary string, args []string) (Response, error) {
	fmt.Fprintf(w, "tollgate: allow %s %v? [y/n/a] ", binary, args)
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if errors.Is(err, io.EOF) {
		return dialogFn(binary, args)
	}
	if err != nil {
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

func osascriptDialog(binary string, args []string) (Response, error) {
	msg := fmt.Sprintf("tollgate: allow %s %s?", binary, strings.Join(args, " "))
	msg = strings.ReplaceAll(msg, `"`, `\"`)
	script := fmt.Sprintf(
		`display dialog "%s" buttons {"Deny", "Allow"} default button "Deny" with title "tollgate"`,
		msg,
	)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return No, nil
	}
	if strings.Contains(string(out), "Allow") {
		return Yes, nil
	}
	return No, nil
}
