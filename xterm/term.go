package xterm

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
)

func StdinRead() (byte, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}

	var b [1]byte
	_, err = os.Stdin.Read(b[:])
	term.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func StdinReadPassword(msg ...any) string {
	var pass []byte
	var err error
	for len(pass) == 0 {
		fmt.Print(msg...)
		fmt.Print(": ")
		//syscall.Stdin doesn't work on Windows
		//pass, err = terminal.ReadPassword(syscall.Stdin)
		pass, err = terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		fmt.Println()
	}
	return string(pass)
}

func StdinReadPassword2(prompt1, prompt2 string) *string {
	for i := 0; i < 3; i++ {
		pass1 := StdinReadPassword(prompt1)
		pass2 := StdinReadPassword(prompt2)
		if pass1 == pass2 {
			return &pass1
		}
	}
	return nil
}

func StdinConfirmInput(answer string) bool {
	answer = strings.TrimSpace(answer)
	if answer == "" {
		panic("answer cannot be empty")
	}
	fmt.Printf("Enter '%s' to confirm: ", answer)
	var actual string
	fmt.Scanln(&actual)
	return actual == answer
}
