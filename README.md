<p align="center">
<img src="logo.png" width="418" height="84" border="0" alt="ACI vetR collector">
<br/>
ACI health check data collector
<p>
<hr/>

This tool collects data from the APIC to be used by Cisco Services in the ACI Health Check.

Binary releases are available [in the releases tab](https://github.com/brightpuddle/vetr-collector/releases). Please always use the latest release unless you have a known requirement to use an earlier version.

# Purpose

This tool performs data collection for the ACI health check. This tool can be run from any computer with access to the APIC, including the APIC itself.

Once the collection is complete, the tool will create a `aci-vetr-data.zip` file. This file should be provided to the Cisco Services ACI consulting engineer for further analysis.

The tool also creates a log file that can be reviewed and/or provided to Cisco to troubleshoot any issues with the collection process. Note, that this file will only be available in a failure scenario; upon successful collection this file is bundled into the `aci-vetr-data.zip` file along with collection data.

# How it works

The tool collects data from a number of endpoints on the APIC for configuration, current faults, scale-related data, etc. The results of these queries are archived in a zip file to be shared with Cisco. The tool currently has no interaction with the switches--all data is collected from the APIC, via the API.

The following API managed objects are queried by this tool:

```
/api/class/topSystem
/api/class/eqptBoard
/api/class/fabricNode
/api/class/fabricSetupP
/api/class/epLoopProtectP
/api/class/epControlP
/api/class/epIpAgingP
/api/class/infraSetPol
/api/class/infraPortTrackPol
/api/class/coopPol
/api/class/fvAEPg
/api/class/fvRsBd
/api/class/fvBD
/api/class/fvCtx
/api/class/fvTenant
/api/class/fvSubnet
/api/class/vzBrCP
/api/class/vzFilter
/api/class/vzSubj
/api/class/vzRsSubjFiltAtt
/api/class/fvRsProv
/api/class/fvRsCons
/api/class/l3extOut
/api/class/l3extLNodeP
/api/class/l3extRsNodeL3OutAtt
/api/class/l3extLIfP
/api/class/l3extInstP
/api/class/isisDomPol
/api/class/bgpRRNodePEp
/api/class/l3IfPol
/api/class/fabricNodeControl
/api/class/fabricRsNodeCtrl
/api/class/fabricRsLeNodePGrp
/api/class/fabricNodeBlk
/api/class/mcpIfPol
/api/class/infraRsMcpIfPol
/api/class/infraRsAccBaseGrp
/api/class/infraRsAccPortP
/api/class/mcpInstPol
/api/class/infraAttEntityP
/api/class/infraRsDomP
/api/class/infraRsVlanNs
/api/class/fvnsEncapBlk
/api/class/firmwareRunning
/api/class/firmwareCtrlrRunning
/api/class/pkiExportEncryptionKey
/api/class/faultInst
/api/class/fvcapRule
/api/class/fvCEp
/api/class/fvIp
/api/class/vnsCDev
/api/class/vnsGraphInst
/api/class/ctxClassCnt
/api/class/fabricHealthTotal
/api/class/topSystem
/api/class/eqptcapacityVlanUsage5min
/api/class/eqptcapacityPolUsage5min
/api/class/eqptcapacityL2Usage5min
/api/class/eqptcapacityL2RemoteUsage5min
/api/class/eqptcapacityL2TotalUsage5min
/api/class/eqptcapacityL3Usage5min
/api/class/eqptcapacityL3UsageCap5min
/api/class/eqptcapacityL3RemoteUsage5min
/api/class/eqptcapacityL3RemoteUsageCap5min
/api/class/eqptcapacityL3TotalUsage5min
/api/class/eqptcapacityL3TotalUsageCap5min
/api/class/eqptcapacityMcastUsage5min
```

# Security

This tool only collects the output of the afformentioned managed objects. Documentation on these endpoints is available in the [full API documentation](https://developer.cisco.com/site/apic-mim-ref-api/). Credentials are only used at the point of collection and are not stored in any way.

All data provided to Cisco will be maintained under Cisco's data retention policy.

# Usage

All command line parameters are optional; the tool will prompt for any missing information. This is a command line tool, but can be run directly from the Windows/Mac/Linux GUI if desired--the tool will pause once complete, before closing the terminal.

```
Usage: aci-vetr-c [--ip IP] [--username USERNAME] [--password PASSWORD] [--output OUTPUT] [--debug]

Options:
  --apic APIC, -a APIC   APIC hostname or IP address
  --username USERNAME, -u USERNAME
                         APIC username
  --password PASSWORD, -p PASSWORD
                         APIC password
  --output OUTPUT, -o OUTPUT
                         Output file [default: aci-vetr-data.zip]
  --icurl                Write requests to icurl script
  --help, -h             display this help and exit
  --version              display version and exit
```
