package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func TestAsk_EOFFallsBackToDialog(t *testing.T) {
	orig := dialogFn
	defer func() { dialogFn = orig }()
	dialogFn = func(_ string, _ []string) (Response, error) { return No, nil }

	resp, err := Ask(&bytes.Buffer{}, strings.NewReader(""), "git", []string{"push"})
	if err != nil {
		t.Fatalf("Ask() error = %v, want nil", err)
	}
	if resp != No {
		t.Errorf("Ask() = %v, want No", resp)
	}
}

func TestAsk_EOFDialogAllowResponse(t *testing.T) {
	orig := dialogFn
	defer func() { dialogFn = orig }()
	dialogFn = func(_ string, _ []string) (Response, error) { return Yes, nil }

	resp, err := Ask(&bytes.Buffer{}, strings.NewReader(""), "git", []string{"push"})
	if err != nil {
		t.Fatalf("Ask() error = %v, want nil", err)
	}
	if resp != Yes {
		t.Errorf("Ask() = %v, want Yes", resp)
	}
}
