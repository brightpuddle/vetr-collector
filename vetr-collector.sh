#!/bin/bash

rm -rf /tmp/vetr-collector > /dev/null
mkdir /tmp/vetr-collector

# Fetch data from the API
icurl -kG https://localhost/api/class/faultInst.json > /tmp/vetr-collector/faultInst.json
icurl -kG https://localhost/api/class/eqptFlash.json > /tmp/vetr-collector/eqptFlash.json
icurl -kG https://localhost/api/class/topSystem.json > /tmp/vetr-collector/topSystem.json
icurl -kG https://localhost/api/class/isisDomPol.json > /tmp/vetr-collector/isisDomPol.json
icurl -kG https://localhost/api/class/fabricSetupP.json > /tmp/vetr-collector/fabricSetupP.json
icurl -kG https://localhost/api/class/eqptStorage.json > /tmp/vetr-collector/eqptStorage.json
icurl -kG https://localhost/api/class/pkiExportEncryptionKey.json > /tmp/vetr-collector/pkiExportEncryptionKey.json
icurl -kG https://localhost/api/class/fvCtx.json > /tmp/vetr-collector/fvCtx.json
icurl -kG https://localhost/api/class/fabricNode.json > /tmp/vetr-collector/fabricNode.json
icurl -kG https://localhost/api/class/fabricLink.json > /tmp/vetr-collector/fabricLink.json
icurl -kG https://localhost/api/class/infraRsVlanNs.json > /tmp/vetr-collector/infraRsVlanNs.json
icurl -kG https://localhost/api/class/l2extDomP.json > /tmp/vetr-collector/l2extDomP.json
icurl -kG https://localhost/api/class/l3extDomP.json > /tmp/vetr-collector/l3extDomP.json
icurl -kG https://localhost/api/class/physDomP.json > /tmp/vetr-collector/physDomP.json
icurl -kG https://localhost/api/class/l2extRsL2DomAtt.json > /tmp/vetr-collector/l2extRsL2DomAtt.json
icurl -kG https://localhost/api/class/l3extRsL3DomAtt.json > /tmp/vetr-collector/l3extRsL3DomAtt.json
icurl -kG https://localhost/api/class/fvRsDomAtt.json > /tmp/vetr-collector/fvRsDomAtt.json
icurl -kG https://localhost/api/class/maintMaintGrp.json > /tmp/vetr-collector/maintMaintGrp.json
icurl -kG https://localhost/api/class/maintRsMgrpp.json > /tmp/vetr-collector/maintRsMgrpp.json
icurl -kG https://localhost/api/class/maintMaintP.json > /tmp/vetr-collector/maintMaintP.json
icurl -kG https://localhost/api/class/fabricNodeBlk.json > /tmp/vetr-collector/fabricNodeBlk.json
icurl -kG https://localhost/api/class/firmwareRunning.json > /tmp/vetr-collector/firmwareRunning.json
icurl -kG https://localhost/api/class/firmwareCtrlrRunning.json > /tmp/vetr-collector/firmwareCtrlrRunning.json
icurl -kG https://localhost/api/class/fabricHealthTotal.json > /tmp/vetr-collector/fabricHealthTotal.json
icurl -kG https://localhost/api/class/healthInst.json -d 'query-target-filter=wcard(healthInst.dn,"^sys/health$")' > /tmp/vetr-collector/healthInst.json
icurl -kG https://localhost/api/class/infraSetPol.json > /tmp/vetr-collector/infraSetPol.json
icurl -kG https://localhost/api/class/fvRsPathAtt.json > /tmp/vetr-collector/fvRsPathAtt.json
icurl -kG https://localhost/api/class/fvTenant.json > /tmp/vetr-collector/fvTenant.json
icurl -kG https://localhost/api/class/fvBD.json > /tmp/vetr-collector/fvBD.json
icurl -kG https://localhost/api/class/vzBrCP.json > /tmp/vetr-collector/vzBrCP.json
icurl -kG https://localhost/api/class/fvAEPg.json > /tmp/vetr-collector/fvAEPg.json
icurl -kG https://localhost/api/class/l3extOut.json > /tmp/vetr-collector/l3extOut.json
icurl -kG https://localhost/api/class/epLoopProtectP.json > /tmp/vetr-collector/epLoopProtectP.json
icurl -kG https://localhost/api/class/bgpRRNodePEp.json > /tmp/vetr-collector/bgpRRNodePEp.json
icurl -kG https://localhost/api/class/eqptcapacityFSPartition.json > /tmp/vetr-collector/eqptcapacityFSPartition.json
icurl -kG https://localhost/api/class/fvSubnet.json > /tmp/vetr-collector/fvSubnet.json
icurl -kG https://localhost/api/class/fvRsBd.json > /tmp/vetr-collector/fvRsBd.json
icurl -kG https://localhost/api/class/l3extRsPathL3OutAtt.json > /tmp/vetr-collector/l3extRsPathL3OutAtt.json
icurl -kG https://localhost/api/class/ipv4Addr.json > /tmp/vetr-collector/ipv4Addr.json
icurl -kG https://localhost/api/class/ipv6Addr.json > /tmp/vetr-collector/ipv6Addr.json
icurl -kG https://localhost/api/class/eqptExtCh.json > /tmp/vetr-collector/eqptExtCh.json
icurl -kG https://localhost/api/class/coopPol.json > /tmp/vetr-collector/coopPol.json
icurl -kG https://localhost/api/class/mcpInstPol.json > /tmp/vetr-collector/mcpInstPol.json
icurl -kG https://localhost/api/class/l3extLNodeP.json > /tmp/vetr-collector/l3extLNodeP.json
icurl -kG https://localhost/api/class/l3extRsNodeL3OutAtt.json > /tmp/vetr-collector/l3extRsNodeL3OutAtt.json
icurl -kG https://localhost/api/class/fvnsVlanInstP.json > /tmp/vetr-collector/fvnsVlanInstP.json
icurl -kG https://localhost/api/class/fvnsEncapBlk.json > /tmp/vetr-collector/fvnsEncapBlk.json
icurl -kG https://localhost/api/class/infraAttEntityP.json > /tmp/vetr-collector/infraAttEntityP.json
icurl -kG https://localhost/api/class/infraRsDomP.json > /tmp/vetr-collector/infraRsDomP.json
icurl -kG https://localhost/api/class/infraRsFuncToEpg.json > /tmp/vetr-collector/infraRsFuncToEpg.json
icurl -kG https://localhost/api/class/infraPortTrackPol.json > /tmp/vetr-collector/infraPortTrackPol.json
icurl -kG https://localhost/api/class/fabricRsTimePol.json > /tmp/vetr-collector/fabricRsTimePol.json
icurl -kG https://localhost/api/class/datetimePol.json > /tmp/vetr-collector/datetimePol.json
icurl -kG https://localhost/api/class/datetimeNtpProv.json > /tmp/vetr-collector/datetimeNtpProv.json
icurl -kG https://localhost/api/class/fcDomP.json > /tmp/vetr-collector/fcDomP.json
icurl -kG https://localhost/api/class/vmmDomP.json > /tmp/vetr-collector/vmmDomP.json
icurl -kG https://localhost/api/class/infraRsAttEntP.json > /tmp/vetr-collector/infraRsAttEntP.json
icurl -kG https://localhost/api/class/apPlugin.json > /tmp/vetr-collector/apPlugin.json
icurl -kG https://localhost/api/class/l3IfPol.json > /tmp/vetr-collector/l3IfPol.json
icurl -kG https://localhost/api/class/fabricExplicitGEp.json > /tmp/vetr-collector/fabricExplicitGEp.json
icurl -kG https://localhost/api/class/fabricNodePEp.json > /tmp/vetr-collector/fabricNodePEp.json
icurl -kG https://localhost/api/class/epControlP.json > /tmp/vetr-collector/epControlP.json
icurl -kG https://localhost/api/class/fvRsCtx.json > /tmp/vetr-collector/fvRsCtx.json
icurl -kG https://localhost/api/class/fabricRsLeNodePGrp.json > /tmp/vetr-collector/fabricRsLeNodePGrp.json
icurl -kG https://localhost/api/class/fabricRsNodeCtrl.json > /tmp/vetr-collector/fabricRsNodeCtrl.json
icurl -kG https://localhost/api/class/fabricNodeControl.json > /tmp/vetr-collector/fabricNodeControl.json
icurl -kG https://localhost/api/class/fabricRsSpNodePGrp.json > /tmp/vetr-collector/fabricRsSpNodePGrp.json
icurl -kG https://localhost/api/class/configRsRemotePath.json > /tmp/vetr-collector/configRsRemotePath.json
icurl -kG https://localhost/api/class/fvcapRule.json > /tmp/vetr-collector/fvcapRule.json
icurl -kG https://localhost/api/class/fvCEp.json -d 'rsp-subtree-include=count' > /tmp/vetr-collector/fvCEp.json
icurl -kG https://localhost/api/class/vzFilter.json > /tmp/vetr-collector/vzFilter.json
icurl -kG https://localhost/api/class/fabricCtrlrConfigP.json > /tmp/vetr-collector/fabricCtrlrConfigP.json
icurl -kG https://localhost/api/class/l3extInstP.json > /tmp/vetr-collector/l3extInstP.json
icurl -kG https://localhost/api/class/fvRsCons.json > /tmp/vetr-collector/fvRsCons.json
icurl -kG https://localhost/api/class/mcpIfPol.json > /tmp/vetr-collector/mcpIfPol.json
icurl -kG https://localhost/api/class/infraRsMcpIfPol.json > /tmp/vetr-collector/infraRsMcpIfPol.json
icurl -kG https://localhost/api/class/infraRsAccBaseGrp.json > /tmp/vetr-collector/infraRsAccBaseGrp.json
icurl -kG https://localhost/api/class/infraRsAccPortP.json > /tmp/vetr-collector/infraRsAccPortP.json
icurl -kG https://localhost/api/class/ctxClassCnt.json -d 'rsp-subtree-class=l2BD,fvEpP,l3Dom' > /tmp/vetr-collector/ctxClassCnt.json
icurl -kG https://localhost/api/class/eqptcapacityVlanUsage5min.json > /tmp/vetr-collector/eqptcapacityVlanUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2Usage5min.json > /tmp/vetr-collector/eqptcapacityL2Usage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2RemoteUsage5min.json > /tmp/vetr-collector/eqptcapacityL2RemoteUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2TotalUsage5min.json > /tmp/vetr-collector/eqptcapacityL2TotalUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3Usage5min.json > /tmp/vetr-collector/eqptcapacityL3Usage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3RemoteUsage5min.json > /tmp/vetr-collector/eqptcapacityL3RemoteUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3TotalUsage5min.json > /tmp/vetr-collector/eqptcapacityL3TotalUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3TotalUsageCap5min.json > /tmp/vetr-collector/eqptcapacityL3TotalUsageCap5min.json
icurl -kG https://localhost/api/class/eqptcapacityPolUsage5min.json > /tmp/vetr-collector/eqptcapacityPolUsage5min.json
icurl -kG https://localhost/api/class/infraWiNode.json > /tmp/vetr-collector/infraWiNode.json
icurl -kG https://localhost/api/class/epIpAgingP.json > /tmp/vetr-collector/epIpAgingP.json
icurl -kG https://localhost/api/class/eqptFt.json > /tmp/vetr-collector/eqptFt.json
icurl -kG https://localhost/api/class/eqptFC.json > /tmp/vetr-collector/eqptFC.json
icurl -kG https://localhost/api/class/eqptSupC.json > /tmp/vetr-collector/eqptSupC.json
icurl -kG https://localhost/api/class/eqptPsu.json > /tmp/vetr-collector/eqptPsu.json
icurl -kG https://localhost/api/class/eqptLC.json > /tmp/vetr-collector/eqptLC.json
icurl -kG https://localhost/api/class/eqptSysC.json > /tmp/vetr-collector/eqptSysC.json
icurl -kG https://localhost/api/class/cdpAdjEp.json > /tmp/vetr-collector/cdpAdjEp.json
icurl -kG https://localhost/api/class/lldpAdjEp.json > /tmp/vetr-collector/lldpAdjEp.json

# Zip result
zip -mj ~/aci-vetr-data.zip /tmp/vetr-collector/*.json

# Cleanup

rm -rf /tmp/vetr-collector

echo Collection complete
echo Output writen to ~/aci-vetr-data.zip, i.e. user home folder
echo Please provide aci-vetr-data.zip to Cisco for analysis.
