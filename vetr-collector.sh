#!/bin/bash

mkdir /tmp/aci-vetr-collections

# Fetch data from API
icurl -kG https://localhost//api/class/topSystem.json > /tmp/aci-vetr-collections/topSystem.json
icurl -kG https://localhost//api/class/eqptBoard.json > /tmp/aci-vetr-collections/eqptBoard.json
icurl -kG https://localhost//api/class/fabricNode.json > /tmp/aci-vetr-collections/fabricNode.json
icurl -kG https://localhost//api/class/fabricSetupP.json > /tmp/aci-vetr-collections/fabricSetupP.json
icurl -kG https://localhost//api/class/epLoopProtectP.json > /tmp/aci-vetr-collections/epLoopProtectP.json
icurl -kG https://localhost//api/class/epControlP.json > /tmp/aci-vetr-collections/epControlP.json
icurl -kG https://localhost//api/class/epIpAgingP.json > /tmp/aci-vetr-collections/epIpAgingP.json
icurl -kG https://localhost//api/class/infraSetPol.json > /tmp/aci-vetr-collections/infraSetPol.json
icurl -kG https://localhost//api/class/infraPortTrackPol.json > /tmp/aci-vetr-collections/infraPortTrackPol.json
icurl -kG https://localhost//api/class/coopPol.json > /tmp/aci-vetr-collections/coopPol.json
icurl -kG https://localhost//api/class/fvAEPg.json > /tmp/aci-vetr-collections/fvAEPg.json
icurl -kG https://localhost//api/class/fvRsBd.json > /tmp/aci-vetr-collections/fvRsBd.json
icurl -kG https://localhost//api/class/fvBD.json > /tmp/aci-vetr-collections/fvBD.json
icurl -kG https://localhost//api/class/fvCtx.json > /tmp/aci-vetr-collections/fvCtx.json
icurl -kG https://localhost//api/class/fvTenant.json > /tmp/aci-vetr-collections/fvTenant.json
icurl -kG https://localhost//api/class/fvSubnet.json > /tmp/aci-vetr-collections/fvSubnet.json
icurl -kG https://localhost//api/class/vzBrCP.json > /tmp/aci-vetr-collections/vzBrCP.json
icurl -kG https://localhost//api/class/vzFilter.json > /tmp/aci-vetr-collections/vzFilter.json
icurl -kG https://localhost//api/class/vzSubj.json > /tmp/aci-vetr-collections/vzSubj.json
icurl -kG https://localhost//api/class/vzRsSubjFiltAtt.json > /tmp/aci-vetr-collections/vzRsSubjFiltAtt.json
icurl -kG https://localhost//api/class/fvRsProv.json > /tmp/aci-vetr-collections/fvRsProv.json
icurl -kG https://localhost//api/class/fvRsCons.json > /tmp/aci-vetr-collections/fvRsCons.json
icurl -kG https://localhost//api/class/l3extOut.json > /tmp/aci-vetr-collections/l3extOut.json
icurl -kG https://localhost//api/class/l3extLNodeP.json > /tmp/aci-vetr-collections/l3extLNodeP.json
icurl -kG https://localhost//api/class/l3extRsNodeL3OutAtt.json > /tmp/aci-vetr-collections/l3extRsNodeL3OutAtt.json
icurl -kG https://localhost//api/class/l3extLIfP.json > /tmp/aci-vetr-collections/l3extLIfP.json
icurl -kG https://localhost//api/class/l3extInstP.json > /tmp/aci-vetr-collections/l3extInstP.json
icurl -kG https://localhost//api/class/isisDomPol.json > /tmp/aci-vetr-collections/isisDomPol.json
icurl -kG https://localhost//api/class/bgpRRNodePEp.json > /tmp/aci-vetr-collections/bgpRRNodePEp.json
icurl -kG https://localhost//api/class/l3IfPol.json > /tmp/aci-vetr-collections/l3IfPol.json
icurl -kG https://localhost//api/class/fabricNodeControl.json > /tmp/aci-vetr-collections/fabricNodeControl.json
icurl -kG https://localhost//api/class/fabricRsNodeCtrl.json > /tmp/aci-vetr-collections/fabricRsNodeCtrl.json
icurl -kG https://localhost//api/class/fabricRsLeNodePGrp.json > /tmp/aci-vetr-collections/fabricRsLeNodePGrp.json
icurl -kG https://localhost//api/class/fabricNodeBlk.json > /tmp/aci-vetr-collections/fabricNodeBlk.json
icurl -kG https://localhost//api/class/mcpIfPol.json > /tmp/aci-vetr-collections/mcpIfPol.json
icurl -kG https://localhost//api/class/infraRsMcpIfPol.json > /tmp/aci-vetr-collections/infraRsMcpIfPol.json
icurl -kG https://localhost//api/class/infraRsAccBaseGrp.json > /tmp/aci-vetr-collections/infraRsAccBaseGrp.json
icurl -kG https://localhost//api/class/infraRsAccPortP.json > /tmp/aci-vetr-collections/infraRsAccPortP.json
icurl -kG https://localhost//api/class/mcpInstPol.json > /tmp/aci-vetr-collections/mcpInstPol.json
icurl -kG https://localhost//api/class/infraAttEntityP.json > /tmp/aci-vetr-collections/infraAttEntityP.json
icurl -kG https://localhost//api/class/infraRsDomP.json > /tmp/aci-vetr-collections/infraRsDomP.json
icurl -kG https://localhost//api/class/infraRsVlanNs.json > /tmp/aci-vetr-collections/infraRsVlanNs.json
icurl -kG https://localhost//api/class/fvnsEncapBlk.json > /tmp/aci-vetr-collections/fvnsEncapBlk.json
icurl -kG https://localhost//api/class/firmwareRunning.json > /tmp/aci-vetr-collections/firmwareRunning.json
icurl -kG https://localhost//api/class/firmwareCtrlrRunning.json > /tmp/aci-vetr-collections/firmwareCtrlrRunning.json
icurl -kG https://localhost//api/class/pkiExportEncryptionKey.json > /tmp/aci-vetr-collections/pkiExportEncryptionKey.json
icurl -kG https://localhost//api/class/faultInst.json > /tmp/aci-vetr-collections/faultInst.json
icurl -kG https://localhost//api/class/fvcapRule.json > /tmp/aci-vetr-collections/fvcapRule.json
icurl -kG https://localhost//api/class/fvCEp.json -d 'rsp-subtree-include=count' > /tmp/aci-vetr-collections/fvCEp.json
icurl -kG https://localhost//api/class/fvIp.json -d 'rsp-subtree-include=count' > /tmp/aci-vetr-collections/fvIp.json
icurl -kG https://localhost//api/class/vnsCDev.json -d 'rsp-subtree-include=count' > /tmp/aci-vetr-collections/vnsCDev.json
icurl -kG https://localhost//api/class/vnsGraphInst.json -d 'rsp-subtree-include=count' > /tmp/aci-vetr-collections/vnsGraphInst.json
icurl -kG https://localhost//api/class/ctxClassCnt.json -d 'rsp-subtree-class=l2BD,fvEpP,l3Dom' > /tmp/aci-vetr-collections/ctxClassCnt.json
icurl -kG https://localhost//api/class/fabricHealthTotal.json > /tmp/aci-vetr-collections/fabricHealthTotal.json
icurl -kG https://localhost//api/class/topSystem.json -d 'rsp-subtree-include=health,no-scoped' > /tmp/aci-vetr-collections/healthInst.json
icurl -kG https://localhost//api/class/eqptcapacityVlanUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityVlanUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityPolUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityPolUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL2Usage5min.json > /tmp/aci-vetr-collections/eqptcapacityL2Usage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL2RemoteUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityL2RemoteUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL2TotalUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityL2TotalUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3Usage5min.json > /tmp/aci-vetr-collections/eqptcapacityL3Usage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3UsageCap5min.json > /tmp/aci-vetr-collections/eqptcapacityL3UsageCap5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3RemoteUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityL3RemoteUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3RemoteUsageCap5min.json > /tmp/aci-vetr-collections/eqptcapacityL3RemoteUsageCap5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3TotalUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityL3TotalUsage5min.json
icurl -kG https://localhost//api/class/eqptcapacityL3TotalUsageCap5min.json > /tmp/aci-vetr-collections/eqptcapacityL3TotalUsageCap5min.json
icurl -kG https://localhost//api/class/eqptcapacityMcastUsage5min.json > /tmp/aci-vetr-collections/eqptcapacityMcastUsage5min.json

# Zip result
zip -mj ~/aci-vetr-raw.zip /tmp/aci-vetr-collections/*.json

# Cleanup
rm -rf /tmp/aci-vetr-collections

echo Collection complete.
echo Provide Cisco Services the aci-vetr-raw.zip file.