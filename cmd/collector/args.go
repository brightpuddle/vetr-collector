package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/alexflint/go-arg"
	"golang.org/x/term"
)

const resultZip = "aci-vetr-data.zip"

var version = "(dev)"

// input collects CLI input.
func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.Trim(input, "\r\n")
}

// Args are command line parameters.
type Args struct {
	URL               string            `arg:"--url,env:ACI_URL" help:"APIC hostname or IP address"`
	Username          string            `arg:"--username,env:ACI_USERNAME" help:"APIC username"`
	Password          string            `arg:"--password,env:ACI_PASSWORD" help:"APIC password"`
	Output            string            `arg:"-o" help:"Output file"`
	ConfigFile        string            `arg:"-c,--config" help:"Path to YAML configuration file"`
	RequestRetryCount int               `arg:"--request-retry-count" help:"Times to retry a failed request" default:"3"`
	RetryDelay        int               `arg:"--retry-delay" help:"Seconds to wait before retry" default:"10"`
	BatchSize         int               `arg:"--batch-size" help:"Max request to send in parallel" default:"7"`
	PageSize          int               `arg:"--page-size" help:"Object per page for large datasets" default:"1000"`
	Confirm           bool              `arg:"-y" help:"Skip confirmation"`
	Verbose           bool              `arg:"-v,--verbose" help:"Enable verbose (debug level) logging"`
	Class             string            `arg:"--class" help:"Collect a single class" default:"all"`
	Query             map[string]string `arg:"-q" help:"Query(s) to filter single class query"`
}

// Description is the CLI description string.
func (Args) Description() string {
	return "ACI vetR collector"
}

// Version is the CLI version string.
func (Args) Version() string {
	return version
}

// NewArgs collects the CLI args and creates a new 'Args'.
func newArgs() Args {
	args := Args{Output: resultZip}
	arg.MustParse(&args)

	if args.URL == "" {
		args.URL = input("APIC IP:")
	}
	args.URL, _ = strings.CutPrefix(args.URL, "http://")
	args.URL, _ = strings.CutPrefix(args.URL, "https://")

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
