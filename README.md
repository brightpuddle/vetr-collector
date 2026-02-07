<p align="center"> <img src="logo.png" height="96" border="0" alt="ACI vetR collector"> <p>

This tool collects data from the APIC to be used by Cisco Services in the ACI Health Check.

Binary releases are available [in the releases tab](https://github.com/brightpuddle/vetr-collector/releases/latest). It's recommended to always use the latest release unless you have a known requirement to use an earlier version.

Purpose
=======

This tool performs data collection for the ACI health check. This tool can be run from any computer with access to the APIC, including the APIC itself.

Once the collection is complete, the tool will create an `aci-vetr-data.zip` file in single-fabric mode. In multi-fabric mode, the tool produces one zip per fabric and an aggregate `aci-collection.zip` that contains all fabric archives. These files should be provided to the Cisco Services ACI consulting engineer for further analysis.

Note that in addition to Cisco Services analysis, this file can also be read by open source third party tools to review the configuration directly. See the [Third Party Tooling](#third-party-tooling) section for more details.

The tool also creates a log file that can be reviewed and/or provided to Cisco to troubleshoot any issues with the collection process. Note, that this file will only be available in a failure scenario; upon successful collection this file is bundled into the `aci-vetr-data.zip` file along with collection data.

How it works
============

The tool collects data from a number of endpoints on the APIC for configuration, current faults, scale-related data, etc. The results of these queries are archived in a zip file to be shared with Cisco. The tool currently has no interaction with the switches--all data is collected from the APIC, via the API.

The following file can be referenced to see the class queries performed by this tool:

https://github.com/brightpuddle/vetr-collector/blob/master/pkg/req/requests.go

**Note** that this file is part of the CI/CD process for this tool, so is always up to date with the latest query data.

Safety/Security
===============

-	All of the queries performed by this tool are also performed by the APIC GUI, so there is no more risk than clicking through the GUI.
-	Queries to the APIC are batched and throttled as to ensure reduced load on the APIC. Again, this results in less impact to the API than the GUI.
-	The APIC has internal safeguards to protect against excess API usage
-	API interaction in ACI has no impact on data forwarding behavior
-	This tool is open source and can be compiled manually with the Go compiler

This tool only collects the output of the afformentioned managed objects. Documentation on these endpoints is available in the [full API documentation](https://developer.cisco.com/site/apic-mim-ref-api/). Credentials are only used at the point of collection and are not stored in any way.

All data provided to Cisco will be maintained under Cisco's [data retention policy](https://www.cisco.com/c/en/us/about/trust-center/global-privacy-policy.html).

Lastly, the binary collector is not strictly required. The releases downloads also include a shell script, named `vetr-collector.sh`. This file can be copied up to the APIC using SCP, and run locally. The script uses icurl and zip to generate the same output as the binary collector. Note, that the script will need to be marked as executable to run on the APIC, i.e. `chmod +x vetr-collector.sh`. This is a more involved process and doesn't include the batching, throttling, and pagination capabilities of the binary collector, but can be used as an alternative collection mechanism if required.

Usage
=====

All command line parameters are optional; the tool will prompt for any missing information. Use the `--help` option to see this output from the CLI.

**Note** that only `url`, `username`, and `password` are typically required. The remainder of the options exist to work around uncommon connectivity challenges, e.g. a long RTT or slow response from the APIC.

## Single Fabric Mode (CLI)

The traditional CLI mode allows collecting from a single fabric:

```bash
./collector --url 10.1.1.1 --username admin --password cisco123
```

## Multi-Fabric Mode (Config File)

For collecting from multiple fabrics in parallel, use a YAML configuration file:

```bash
./collector --config fabrics.yaml
```

Example configuration file:

```yaml
global:
  username: admin
  password: cisco123
  confirm: true

fabrics:
  - name: production
    url: 10.1.1.1
  - name: staging
    url: 10.2.2.2
    username: staging-user
    password: staging-pass
  - name: development
    url: 10.3.3.3
```

  For a fully documented example with comments for every option, see [config-example.yaml](config-example.yaml).

### Config File Features

- **Parallel Collection**: All fabrics are collected simultaneously using goroutines
- **Global Settings**: Define common settings once in the `global` section
- **Per-Fabric Overrides**: Override any global setting per fabric
- **Flexible Naming**: 
  - If `name` is specified, output is `{name}.zip`
  - If no `name`, output is `{url}.zip`
- **Aggregate Archive**: After collecting all fabrics, the tool creates `aci-collection.zip` containing all per-fabric zip files
- **Fabric Context in Logs**: Each log message includes the fabric name for easy tracking
- **Validation**: Ensures fabric names/URLs are unique and required fields are present

### Supported Global/Fabric Settings

All CLI parameters can be used in the config file:
- `username` - APIC username
- `password` - APIC password
- `request_retry_count` - Times to retry failed requests (default: 3)
- `retry_delay` - Seconds to wait before retry (default: 10)
- `batch_size` - Max parallel requests (default: 7)
- `page_size` - Objects per page for large datasets (default: 1000)
- `confirm` - Skip confirmation prompts (default: false)
- `verbose` - Enable debug level logging (default: false)
- `class` - Collect single class (default: all)
- `query` - Query filters for single class

**Note**: `url` must be specified per fabric and is not supported as a global setting.

## Verbose Logging

Enable debug-level logging for detailed progress:

```bash
# Single fabric mode
./collector --url 10.1.1.1 --username admin --password cisco123 --verbose

# Multi-fabric mode
./collector --config fabrics.yaml --verbose

# Or set in config file
global:
  verbose: true
```

Debug logging shows:
- Individual query start/end times
- Detailed fetch progress for each class
- Fine-grained timing information

Standard Info logging shows:
- Batch completion messages
- Authentication status
- Major collection milestones

## Command Line Options

```
ACI vetR collector
version ...
Usage: collector [--url URL] [--username USERNAME] [--password PASSWORD] [--output OUTPUT] [--config CONFIG] [--request-retry-count REQUEST-RETRY-COUNT] [--retry-delay RETRY-DELAY] [--batch-size BATCH-SIZE] [--page-size PAGE-SIZE] [--confirm] [--verbose] [--class CLASS] [--query QUERY]

Options:
  --url URL              APIC hostname or IP address [env: ACI_URL]
  --username USERNAME    APIC username [env: ACI_USERNAME]
  --password PASSWORD    APIC password [env: ACI_PASSWORD]
  --output OUTPUT, -o OUTPUT
                         Output file [default: aci-vetr-data.zip]
  --config CONFIG, -c CONFIG
                         Path to YAML configuration file
  --request-retry-count REQUEST-RETRY-COUNT
                         Times to retry a failed request [default: 3]
  --retry-delay RETRY-DELAY
                         Seconds to wait before retry [default: 10]
  --batch-size BATCH-SIZE
                         Max request to send in parallel [default: 7]
  --page-size PAGE-SIZE
                         Object per page for large datasets [default: 1000]
  --confirm, -y          Skip confirmation
  --verbose, -v          Enable verbose (debug level) logging
  --class CLASS          Collect a single class [default: all]
  --query QUERY, -q QUERY
                         Query(s) to filter single class query
  --help, -h             display this help and exit
  --version              display version and exit
```

Performance and Troubleshooting
-------------------------------

In general the collector is expected to run very quickly and have no issues. That said, one error sometimes encountered is a class with too much data. As an example of this, suppose a fabric has a very large number of static path bindings. The collector queries objects by class, so all static path bindings will be requested in a single query, and when a response has too much data, instead of sending the response data, the APIC will respond with an error.

The collector addresses this with pagination. Pagination allows querying large datasets in "pages," which are groups of that object, e.g. static path binding 1-999, then 1000-1999, and so on. The actual byte size of a "page" of data will vary, as individual object sizes vary.

Reasonable defaults are provided to handle this; however, they may not work for every scenario. The two options are `--page-size` and `--batch-size`.

Page size defines how many objects will be sent back in each page query, so a page size of 1000 will try to query objects in pages of 1000 and 5000 groups of 5000. A larger page size means less total queries, so may be more performance, but at some point will run over the APIC's size limits.

Batch size determines how many queries are sent to the APIC before waiting for a response. The collector sends queries in parallel for faster performance; however, too many queries too quickly will be throttled and the APIC will refuse to respond. If you set `--batch-size 1` the collector will behave synchonously and wait for each query to complete before sending another. This will be slower then sending requests in parallel, but may be helpful for troubleshooting purposes.

Again, these and othe configurable settings should not generally need to be modified, but may be useful in corner cases with unusually large configurations, heavily loaded APICs, etc.

### Running code directly from source

Static binaries are provided for convenience and are generally preferred; however, if you'd like to run the code directly from source, e.g. for security auditing, this is also an option.

1.	[Install Go](https://go.dev/doc/install)
2.	Clone the repo
3.	`go mod download`
4.	`go run ./cmd/collector/*.go`

If on Windows, it's recommended to use Powershell or WSL to avoid issues with ANSI escape sequences and path slash direction.

Third Party Tooling
===================

The following tools can be used to visualize or analyze the vetR collection file. Note that these are not owned by Cisco Systems.

-	[vetR Summarizer](https://github.com/Tes3awy/vetr-summarizer) - Visualize and summarize vetR collection data through a web UI
-	[reQuery](https://github.com/brightpuddle/requery) - Run moquery-like queries against the collection file from the CLI
