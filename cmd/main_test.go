package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	in := `{"Transaction":{"merchant":"Montreal Canadiens","amount":666,"time":"2019-02-13T11:00:00.000Z"}}`
	want := "\n{\"Account\":{},\"violations\":[\"Account-not-initialized\"]}\n\n"
	cd()
	run("make", "install")
	run("chmod", "+x", "authorizer")

	cmd := exec.Command("sh", "authorizer")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, in)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Trim(string(out), "")

	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func cd() {
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("could not change dir: %v", err)
	}
}

func run(name string, args ...string) {
	a := exec.Command( name, args...)
	a.Run()
}
