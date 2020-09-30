package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/alexflint/go-arg"
	"golang.org/x/crypto/ssh/terminal"
)

// input collects CLI input.
func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.Trim(input, "\r\n")
}

// Args are command line parameters.
type Args struct {
	APIC        string `arg:"-a" help:"APIC hostname or IP address"`
	Username    string `arg:"-u" help:"APIC username"`
	Password    string `arg:"-p" help:"APIC password"`
	Output      string `arg:"-o" help:"Output file"`
	WriteScript bool   `help:"Write requests to icurl script"`
	ReadRaw     string `help:"Read raw data from manually collection" placeholder:"FILE"`
}

// Description is the CLI description string.
func (Args) Description() string {
	return "ACI vetR collector"
}

// Version is the CLI version string.
func (Args) Version() string {
	return "version " + version
}

// NewArgs collects the CLI args and creates a new 'Args'.
func newArgs() (Args, error) {
	args := Args{Output: resultZip}
	arg.MustParse(&args)

	switch {
	case args.WriteScript || args.ReadRaw != "":
		return args, nil
	default:
		if args.APIC == "" {
			args.APIC = input("APIC IP:")
		}
		if args.Username == "" {
			args.Username = input("Username:")
		}
		if args.Password == "" {
			fmt.Print("Password: ")
			pwd, _ := terminal.ReadPassword(int(syscall.Stdin))
			args.Password = string(pwd)
		}
	}
	return args, nil
}
