package main

import (
	"bufio"
	"collector"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/alexflint/go-arg"
	"golang.org/x/term"
)

const resultZip = "aci-vetr-data.zip"

// input collects CLI input.
func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.Trim(input, "\r\n")
}

// Args are command line parameters.
type Args struct {
	APIC              string `arg:"-a"                    help:"APIC hostname or IP address"`
	Username          string `arg:"-u"                    help:"APIC username"`
	Password          string `arg:"-p"                    help:"APIC password"`
	Output            string `arg:"-o"                    help:"Output file"`
	RequestRetryCount int    `arg:"--request-retry-count" help:"Times to retry a failed request" default:"3"`
	RetryDelay        int    `arg:"--retry-delay"         help:"Seconds to wait before retry"    default:"10"`
	BatchSize         int    `arg:"--batch-size"          help:"Max request to send in parallel" default:"10"`
	Confirm           bool   `arg:"-y"                    help:"Skip confirmation"`
}

// Description is the CLI description string.
func (Args) Description() string {
	return "ACI vetR collector"
}

// Version is the CLI version string.
func (Args) Version() string {
	return fmt.Sprintf("version %s", collector.Version)
}

// NewArgs collects the CLI args and creates a new 'Args'.
func newArgs() Args {
	args := Args{Output: resultZip}
	arg.MustParse(&args)

	if args.APIC == "" {
		args.APIC = input("APIC IP:")
	}
	if args.Username == "" {
		args.Username = input("Username:")
	}
	if args.Password == "" {
		fmt.Print("Password: ")
		pwd, _ := term.ReadPassword(int(syscall.Stdin))
		args.Password = string(pwd)
		fmt.Println()
	}
	return args
}
