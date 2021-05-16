#!python
import yaml

TMP_FOLDER = "/tmp/vetr-collector"

reqs = yaml.load(open("reqs.yaml"), Loader=yaml.Loader)

print("#!/bin/bash")
print("")
print("mkdir %s" % TMP_FOLDER)
print("rm -rf %s > /dev/null" % TMP_FOLDER)
print("")
print("# Fetch data from the API")

for req in reqs:
    # icurl command
    cmd = ["icurl -kG"]
    cls = req["class"]

    # url
    cmd.append("https://localhost/api/class/%s.json" % cls)
    for k, v in req.get("query", {}).items():
        cmd.append("-d '%s=%s'" % (k, v))

    # Parse the output
    flt = req.get("filter", ".imdata[].%s.attributes" % cls)
    if "#" in flt:
        flt = ".%s" % flt.replace(".#", "[]")
    cmd.append("| jq '%s'" % flt)

    # redirect output to file
    pfx = req.get("prefix", cls)
    cmd.append("> %s/%s.json" % (TMP_FOLDER, pfx))

    cmd = " ".join(cmd)
    print(cmd)

print("")
print("# Zip result")
print("zip -mj ~/aci-vetr-data.zip %s/*.json" % TMP_FOLDER)
print("")
print("#Cleanup")
print("")
print("rm -rf %s" % TMP_FOLDER)
print("")
print("echo Collection complete")
print("echo Please provice aci-vetr-data.zip to Cisco Services for analysis.")
