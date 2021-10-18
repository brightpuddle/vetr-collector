<p align="center">
<img src="logo.png" width="418" height="84" border="0" alt="ACI vetR collector">
<br/>
ACI health check data collector
<p>
<hr/>

This tool collects data from the APIC to be used by Cisco Services in the ACI
Health Check.

Binary releases are available
[in the releases tab](https://github.com/aci-vetr/collector/releases/latest).
Please always use the latest release unless you have a known requirement to use
an earlier version.

# Purpose

This tool performs data collection for the ACI health check. This tool can be
run from any computer with access to the APIC, including the APIC itself.

Once the collection is complete, the tool will create an `aci-vetr-data.zip`
file. This file should be provided to the Cisco Services ACI consulting engineer
for further analysis.

The tool also creates a log file that can be reviewed and/or provided to Cisco
to troubleshoot any issues with the collection process. Note, that this file
will only be available in a failure scenario; upon successful collection this
file is bundled into the `aci-vetr-data.zip` file along with collection data.

# How it works

The tool collects data from a number of endpoints on the APIC for configuration,
current faults, scale-related data, etc. The results of these queries are
archived in a zip file to be shared with Cisco. The tool currently has no
interaction with the switches--all data is collected from the APIC, via the API.

The following file can be referenced to see the class queries performed by this
tool:

https://github.com/aci-vetr/collector/blob/master/pkg/req/reqs.yaml

**Note** that this file is part of the CI/CD process for this tool, so is always
up to date with the latest query data.

# Safety

- All of the queries performed by this tool are also performed by the ACI GUI,
  so there is no more risk than clicking through the GUI.
- Queries to the APIC are batched and throttled as to ensure reduced load on the
  APIC. Again, this results in less impact to the API than the GUI.

# Security

This tool only collects the output of the afformentioned managed objects.
Documentation on these endpoints is available in the
[full API documentation](https://developer.cisco.com/site/apic-mim-ref-api/).
Credentials are only used at the point of collection and are not stored in any
way.

All data provided to Cisco will be maintained under Cisco's data retention
policy.

# Usage

All command line parameters are optional; the tool will prompt for any missing
information. Use the `--help` option to see this output from the CLI.

**Note** that only `apic`, `username`, and `password` are typically required.
The remainder of the options exist to work around uncommon connectivity
challenges, e.g. a long RTT or slow response from the APIC.

```
ACI vetR collector
version ...
Usage: collector [--apic APIC] [--username USERNAME] [--password PASSWORD] [--output OUTPUT] [--request-retry-count REQUEST-RETRY-COUNT] [--retry-delay RETRY-DELAY] [--batch-size BATCH-SIZE]

Options:
  --apic APIC, -a APIC   APIC hostname or IP address
  --username USERNAME, -u USERNAME
                         APIC username
  --password PASSWORD, -p PASSWORD
                         APIC password
  --output OUTPUT, -o OUTPUT
                         Output file [default: aci-vetr-data.zip]
  --request-retry-count REQUEST-RETRY-COUNT
                         Times to retry a failed request [default: 3]
  --retry-delay RETRY-DELAY
                         Seconds to wait before retry [default: 10]
  --batch-size BATCH-SIZE
                         Max request to send in parallel [default: 10]
  --help, -h             display this help and exit
  --version              display version and exit
```
