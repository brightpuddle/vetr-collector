package main

import (
	"collector/pkg/config"

	"github.com/alexflint/go-arg"
)

const resultZip = "aci-vetr-data.zip"

var version = "(dev)"

// Args are command line parameters.
type Args struct {
	URL               string            `arg:"--url,env:ACI_URL"           help:"APIC hostname or IP address"`
	Username          string            `arg:"--username,env:ACI_USERNAME" help:"APIC username"`
	Password          string            `arg:"--password,env:ACI_PASSWORD" help:"APIC password"`
	Output            string            `arg:"-o"                          help:"Output file"`
	ConfigFile        string            `arg:"-c,--config"                 help:"Path to YAML configuration file"`
	RequestRetryCount int               `arg:"--request-retry-count"       help:"Times to retry a failed request"       default:"3"`
	RetryDelay        int               `arg:"--retry-delay"               help:"Seconds to wait before retry"          default:"10"`
	BatchSize         int               `arg:"--batch-size"                help:"Max request to send in parallel"       default:"7"`
	PageSize          int               `arg:"--page-size"                 help:"Object per page for large datasets"    default:"1000"`
	Confirm           bool              `arg:"-y"                          help:"Skip confirmation"`
	Verbose           bool              `arg:"-v,--verbose"                help:"Enable verbose (debug level) logging"`
	Class             string            `arg:"--class"                     help:"Collect a single class"                default:"all"`
	Query             map[string]string `arg:"-q"                          help:"Query(s) to filter single class query"`
}

// Description is the CLI description string.
func (Args) Description() string {
	return "ACI vetR collector"
}

// Version is the CLI version string.
func (Args) Version() string {
	return version
}

// readArgs collects the CLI args and returns a config.Config.
func readArgs() (*config.Config, error) {
	args := Args{Output: resultZip}
	arg.MustParse(&args)

	if args.ConfigFile != "" {
		cfg, err := config.ParseConfig(args.ConfigFile)
		if err != nil {
			return nil, err
		}
		if err := cfg.NormalizeAndPrompt(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg := config.New()
	requestRetryCount := args.RequestRetryCount
	retryDelay := args.RetryDelay
	batchSize := args.BatchSize
	pageSize := args.PageSize
	confirm := args.Confirm
	verbose := args.Verbose

	cfg.Global.Verbose = args.Verbose
	cfg.Fabrics = []config.FabricConfig{{
		URL:               args.URL,
		Output:            args.Output,
		Username:          args.Username,
		Password:          args.Password,
		RequestRetryCount: &requestRetryCount,
		RetryDelay:        &retryDelay,
		BatchSize:         &batchSize,
		PageSize:          &pageSize,
		Confirm:           &confirm,
		Verbose:           &verbose,
		Class:             args.Class,
		Query:             args.Query,
	}}

	if err := cfg.NormalizeAndPrompt(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
