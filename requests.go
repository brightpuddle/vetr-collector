package main

import (
	"fmt"

	"collector/aci"
)

// Mod modifies a aci Request
type Mod = func(*aci.Req)

// Request is an HTTP request.
type Request struct {
	class  string // MO class
	path   string // Request path
	prefix string // Prefix for the DB
	mods   []Mod  // Request modifiers, e.g. query parameters
	filter string // Result filter (default to #.{class}.attributes)
}

func getRequests() []*Request {
	reqs := []*Request{
		/************************************************************
		Infrastructure
		************************************************************/
		{class: "topSystem"},    // All devices
		{class: "eqptBoard"},    // APIC hardware
		{class: "fabricNode"},   // Switch hardware
		{class: "fabricSetupP"}, // Pods (fabric setup policy)

		/************************************************************
		Fabric-wide settings
		************************************************************/
		{class: "epLoopProtectP"},    // EP loop protection policy
		{class: "epControlP"},        // Rogue EP control policy
		{class: "epIpAgingP"},        // IP aging policy
		{class: "infraSetPol"},       // Fabric-wide settings
		{class: "infraPortTrackPol"}, // Port tracking policy
		{class: "coopPol"},           // COOP group policy

		/************************************************************
		Tenants
		************************************************************/
		// Primary constructs
		{class: "fvAEPg"},   // EPG
		{class: "fvRsBd"},   // EPG --> BD
		{class: "fvBD"},     // BD
		{class: "fvCtx"},    // VRF
		{class: "fvTenant"}, // Tenant
		{class: "fvSubnet"}, // Subnet

		// Contracts
		{class: "vzBrCP"},          // Contract
		{class: "vzFilter"},        // Filter
		{class: "vzSubj"},          // Subject
		{class: "vzRsSubjFiltAtt"}, // Subject --> filter
		{class: "fvRsProv"},        // EPG --> contract provided
		{class: "fvRsCons"},        // EPG --> contract consumed

		// L3outs
		{class: "l3extOut"},            // L3out
		{class: "l3extLNodeP"},         // L3 node profile
		{class: "l3extRsNodeL3OutAtt"}, // Node profile --> Node
		{class: "l3extLIfP"},           // L3 interface profile
		{class: "l3extInstP"},          // External EPG

		/************************************************************
		Fabric Policies
		************************************************************/
		{class: "isisDomPol"},         // ISIS policy
		{class: "bgpRRNodePEp"},       // BGP route reflector nodes
		{class: "l3IfPol"},            // L3 interface policy
		{class: "fabricNodeControl"},  // node control (Dom, netflow,etc)
		{class: "fabricRsNodeCtrl"},   // node policy group --> node control
		{class: "fabricRsLeNodePGrp"}, // leaf --> leaf node policy group
		{class: "fabricNodeBlk"},      // Node block

		/************************************************************
		Fabric Access
		************************************************************/
		// MCP
		{class: "mcpIfPol"},          // MCP inteface policy
		{class: "infraRsMcpIfPol"},   // MCP pol --> policy group
		{class: "infraRsAccBaseGrp"}, // policy group --> host port selector
		{class: "infraRsAccPortP"},   // int profile --> node profile

		{class: "mcpInstPol"}, // MCP global policy

		// AEP/domain/VLANs
		{class: "infraAttEntityP"}, // AEP
		{class: "infraRsDomP"},     // AEP --> domain
		{class: "infraRsVlanNs"},   // Domain --> VLAN pool
		{class: "fvnsEncapBlk"},    // VLAN encap block

		/************************************************************
		Admin/Operations
		************************************************************/
		{class: "firmwareRunning"},        // Switch firmware
		{class: "firmwareCtrlrRunning"},   // Controller firmware
		{class: "pkiExportEncryptionKey"}, // Crypto key

		/************************************************************
		Live State
		************************************************************/
		{class: "faultInst"}, // Faults
		{class: "fvcapRule"}, // Capacity rules

		{ // Endpoint count
			class:  "fvCEp",
			filter: "#.moCount.attributes",
			mods:   []Mod{aci.Query("rsp-subtree-include", "count")},
		},
		{ // IP count
			class:  "fvIp",
			filter: "#.moCount.attributes",
			mods:   []Mod{aci.Query("rsp-subtree-include", "count")},
		},

		{ // L4-L7 container count
			class:  "vnsCDev",
			filter: "#.moCount.attributes",
			mods:   []Mod{aci.Query("rsp-subtree-include", "count")},
		},

		{ // L4-L7 service graph count
			class:  "vnsGraphInst",
			filter: "#.moCount.attributes",
			mods:   []Mod{aci.Query("rsp-subtree-include", "count")},
		},

		{ // MO count by node
			class: "ctxClassCnt",
			mods:  []Mod{aci.Query("rsp-subtree-class", "l2BD,fvEpP,l3Dom")},
		},

		// Fabric health
		{class: "fabricHealthTotal"}, // Total and per-pod health scores
		{ // Per-device health stats
			class:  "topSystem",
			prefix: "healthInst",
			mods:   []Mod{aci.Query("rsp-subtree-include", "health,no-scoped")},
			filter: "#.healthInst.attributes",
		},

		// Switch capacity
		{class: "eqptcapacityVlanUsage5min"},        // VLAN
		{class: "eqptcapacityPolUsage5min"},         // TCAM
		{class: "eqptcapacityL2Usage5min"},          // L2 local
		{class: "eqptcapacityL2RemoteUsage5min"},    // L2 remote
		{class: "eqptcapacityL2TotalUsage5min"},     // L2 total
		{class: "eqptcapacityL3Usage5min"},          // L3 local
		{class: "eqptcapacityL3UsageCap5min"},       // L3 local cap
		{class: "eqptcapacityL3RemoteUsage5min"},    // L3 remote
		{class: "eqptcapacityL3RemoteUsageCap5min"}, // L3 remote cap
		{class: "eqptcapacityL3TotalUsage5min"},     // L3 total
		{class: "eqptcapacityL3TotalUsageCap5min"},  // L3 total cap
		{class: "eqptcapacityMcastUsage5min"},       // Multicast
	}

	for _, req := range reqs {
		if req.filter == "" {
			req.filter = fmt.Sprintf("#.%s.attributes", req.class)
		}
		if req.path == "" {
			req.path = "/api/class/" + req.class
		}
		if req.prefix == "" {
			req.prefix = req.class
		}
	}
	return reqs
}
