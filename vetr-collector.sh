#!/bin/bash

rm -rf /tmp/vetr-collector > /dev/null
mkdir /tmp/vetr-collector

# Fetch data from the API
icurl -kG https://localhost/api/class/topSystem.json > /tmp/vetr-collector/topSystem.json
icurl -kG https://localhost/api/class/eqptBoard.json > /tmp/vetr-collector/eqptBoard.json
icurl -kG https://localhost/api/class/fabricNode.json > /tmp/vetr-collector/fabricNode.json
icurl -kG https://localhost/api/class/infraWiNode.json > /tmp/vetr-collector/infraWiNode.json
icurl -kG https://localhost/api/class/fabricLink.json > /tmp/vetr-collector/fabricLink.json
icurl -kG https://localhost/api/class/ctrlrInst.json -d 'rsp-subtree=full' > /tmp/vetr-collector/ctrlrInst.json
icurl -kG https://localhost/api/class/fabricInst.json -d 'rsp-subtree=full' > /tmp/vetr-collector/fabricInst.json
icurl -kG https://localhost/api/class/infraInfra.json -d 'rsp-subtree=full' > /tmp/vetr-collector/infraInfra.json
icurl -kG https://localhost/api/class/fvTenant.json -d 'rsp-subtree=full' > /tmp/vetr-collector/fvTenant.json
icurl -kG https://localhost/api/class/vmmDomP.json -d 'rsp-subtree=full' > /tmp/vetr-collector/vmmDomP.json
icurl -kG https://localhost/api/class/physDomP.json -d 'rsp-subtree=full' > /tmp/vetr-collector/physDomP.json
icurl -kG https://localhost/api/class/l3extDomP.json -d 'rsp-subtree=full' > /tmp/vetr-collector/l3extDomP.json
icurl -kG https://localhost/api/class/l2extDomP.json -d 'rsp-subtree=full' > /tmp/vetr-collector/l2extDomP.json
icurl -kG https://localhost/api/class/firmwareRunning.json > /tmp/vetr-collector/firmwareRunning.json
icurl -kG https://localhost/api/class/maintUpgJob.json > /tmp/vetr-collector/maintUpgJob.json
icurl -kG https://localhost/api/class/firmwareCtrlrRunning.json > /tmp/vetr-collector/firmwareCtrlrRunning.json
icurl -kG https://localhost/api/class/pkiExportEncryptionKey.json > /tmp/vetr-collector/pkiExportEncryptionKey.json
icurl -kG https://localhost/api/class/faultInst.json > /tmp/vetr-collector/faultInst.json
icurl -kG https://localhost/api/class/apUiInfo.json > /tmp/vetr-collector/apUiInfo.json
icurl -kG https://localhost/api/class/apPlugin.json > /tmp/vetr-collector/apPlugin.json
icurl -kG https://localhost/api/class/fvcapRule.json > /tmp/vetr-collector/fvcapRule.json
icurl -kG https://localhost/api/class/fvCEp.json -d 'rsp-subtree-include=count' > /tmp/vetr-collector/fvCEp.json
icurl -kG https://localhost/api/class/fvIp.json -d 'rsp-subtree-include=count' > /tmp/vetr-collector/fvIp.json
icurl -kG https://localhost/api/class/vnsCDev.json -d 'rsp-subtree-include=count' > /tmp/vetr-collector/vnsCDev.json
icurl -kG https://localhost/api/class/vnsGraphInst.json -d 'rsp-subtree-include=count' > /tmp/vetr-collector/vnsGraphInst.json
icurl -kG https://localhost/api/class/ctxClassCnt.json -d 'rsp-subtree-class=l2BD,fvEpP,l3Dom' > /tmp/vetr-collector/ctxClassCnt.json
icurl -kG https://localhost/api/class/eqptcapacityVlanUsage5min.json > /tmp/vetr-collector/eqptcapacityVlanUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityPolUsage5min.json > /tmp/vetr-collector/eqptcapacityPolUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2Usage5min.json > /tmp/vetr-collector/eqptcapacityL2Usage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2RemoteUsage5min.json > /tmp/vetr-collector/eqptcapacityL2RemoteUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL2TotalUsage5min.json > /tmp/vetr-collector/eqptcapacityL2TotalUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3Usage5min.json > /tmp/vetr-collector/eqptcapacityL3Usage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3UsageCap5min.json > /tmp/vetr-collector/eqptcapacityL3UsageCap5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3RemoteUsage5min.json > /tmp/vetr-collector/eqptcapacityL3RemoteUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3RemoteUsageCap5min.json > /tmp/vetr-collector/eqptcapacityL3RemoteUsageCap5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3TotalUsage5min.json > /tmp/vetr-collector/eqptcapacityL3TotalUsage5min.json
icurl -kG https://localhost/api/class/eqptcapacityL3TotalUsageCap5min.json > /tmp/vetr-collector/eqptcapacityL3TotalUsageCap5min.json
icurl -kG https://localhost/api/class/eqptcapacityMcastUsage5min.json > /tmp/vetr-collector/eqptcapacityMcastUsage5min.json
icurl -kG https://localhost/api/class/fabricHealthTotal.json > /tmp/vetr-collector/fabricHealthTotal.json
icurl -kG https://localhost/api/class/topSystem.json -d 'rsp-subtree-include=health,no-scoped' > /tmp/vetr-collector/healthInst.json

# Zip result
zip -mj ~/aci-vetr-data.zip /tmp/vetr-collector/*.json

# Cleanup

rm -rf /tmp/vetr-collector

echo Collection complete
echo Please provide aci-vetr-data.zip to Cisco for analysis.
