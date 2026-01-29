// Package req contains the collector requests
package req

//go:generate go run ../../cmd/genscript/main.go

import "collector/pkg/aci"

// Mod modifies an aci Request
type Mod = func(*aci.Req)

// Request is an HTTP request.
type Request struct {
	Class string            // MO class
	Query map[string]string // Query parameters
}

// Requests contains all the ACI API requests to execute
var Requests = []Request{
	{Class: "faultInst"},
	{Class: "eqptFlash"},
	{Class: "topSystem"},
	{Class: "isisDomPol"},
	{Class: "fabricSetupP"},
	{Class: "eqptStorage"},
	{Class: "pkiExportEncryptionKey"},
	{Class: "fvCtx"},
	{Class: "fabricNode"},
	{Class: "fabricLink"},
	{Class: "infraRsVlanNs"},
	{Class: "l2extDomP"},
	{Class: "l3extDomP"},
	{Class: "physDomP"},
	{Class: "l2extRsL2DomAtt"},
	{Class: "l3extRsL3DomAtt"},
	{Class: "fvRsDomAtt"},
	{Class: "maintMaintGrp"},
	{Class: "maintRsMgrpp"},
	{Class: "maintMaintP"},
	{Class: "fabricNodeBlk"},
	{Class: "firmwareRunning"},
	{Class: "firmwareCtrlrRunning"},
	{Class: "fabricHealthTotal"},
	{
		Class: "healthInst",
		Query: map[string]string{
			"query-target-filter": "wcard(healthInst.dn,\"^sys/health$\")",
		},
	},
	{Class: "infraSetPol"},
	{Class: "fvRsPathAtt"},
	{Class: "fvTenant"},
	{Class: "fvBD"},
	{Class: "vzBrCP"},
	{Class: "fvAEPg"},
	{Class: "l3extOut"},
	{Class: "epLoopProtectP"},
	{Class: "bgpRRNodePEp"},
	{Class: "eqptcapacityFSPartition"},
	{Class: "fvSubnet"},
	{Class: "fvRsBd"},
	{Class: "l3extRsPathL3OutAtt"},
	{Class: "ipv4Addr"},
	{Class: "ipv6Addr"},
	{Class: "eqptExtCh"},
	{Class: "coopPol"},
	{Class: "mcpInstPol"},
	{Class: "l3extLNodeP"},
	{Class: "l3extRsNodeL3OutAtt"},
	{Class: "fvnsVlanInstP"},
	{Class: "fvnsEncapBlk"},
	{Class: "infraAttEntityP"},
	{Class: "infraRsDomP"},
	{Class: "infraRsFuncToEpg"},
	{Class: "infraPortTrackPol"},
	{Class: "fabricRsTimePol"},
	{Class: "datetimePol"},
	{Class: "datetimeNtpProv"},
	{Class: "fcDomP"},
	{Class: "vmmDomP"},
	{Class: "infraRsAttEntP"},
	{Class: "apPlugin"},
	{Class: "l3IfPol"},
	{Class: "fabricExplicitGEp"},
	{Class: "fabricNodePEp"},
	{Class: "epControlP"},
	{Class: "fvRsCtx"},
	{Class: "fabricRsLeNodePGrp"},
	{Class: "fabricRsNodeCtrl"},
	{Class: "fabricNodeControl"},
	{Class: "fabricRsSpNodePGrp"},
	{Class: "configRsRemotePath"},
	{Class: "fvcapRule"},
	{
		Class: "fvCEp",
		Query: map[string]string{
			"rsp-subtree-include": "count",
		},
	},
	{Class: "vzFilter"},
	{Class: "fabricCtrlrConfigP"},
	{Class: "l3extInstP"},
	{Class: "fvRsCons"},
	{Class: "mcpIfPol"},
	{Class: "infraRsMcpIfPol"},
	{Class: "infraRsAccBaseGrp"},
	{Class: "infraRsAccPortP"},
	{
		Class: "ctxClassCnt",
		Query: map[string]string{
			"rsp-subtree-class": "l2BD,fvEpP,l3Dom",
		},
	},
	{Class: "eqptcapacityVlanUsage5min"},
	{Class: "eqptcapacityL2Usage5min"},
	{Class: "eqptcapacityL2RemoteUsage5min"},
	{Class: "eqptcapacityL2TotalUsage5min"},
	{Class: "eqptcapacityL3Usage5min"},
	{Class: "eqptcapacityL3RemoteUsage5min"},
	{Class: "eqptcapacityL3TotalUsage5min"},
	{Class: "eqptcapacityL3TotalUsageCap5min"},
	{Class: "eqptcapacityPolUsage5min"},
	{Class: "infraWiNode"},
	{Class: "epIpAgingP"},
	{Class: "eqptFt"},
	{Class: "eqptFC"},
	{Class: "eqptSupC"},
	{Class: "eqptPsu"},
	{Class: "eqptLC"},
	{Class: "eqptSysC"},
	{Class: "cdpAdjEp"},
	{Class: "lldpAdjEp"},
}

// GetRequests returns normalized requests
func GetRequests() ([]Request, error) {
	return Requests, nil
}
