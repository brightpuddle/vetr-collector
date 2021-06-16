#!python
import os

import yaml

TMP_FOLDER = "/tmp/vetr-collector"
SCRIPT_NAME = "vetr-collector.sh"

reqs = yaml.load(open("reqs.yaml"), Loader=yaml.Loader)

f = open(SCRIPT_NAME, "w")

f.write("#!/bin/bash")
f.write("")
f.write("rm -rf %s > /dev/null" % TMP_FOLDER)
f.write("mkdir %s" % TMP_FOLDER)
f.write("")
f.write("# Fetch data from the API")

for req in reqs:
    # icurl command
    cmd = ["icurl -kG"]
    cls = req["class"]

    # url
    cmd.append("https://localhost/api/class/%s.json" % cls)
    for k, v in req.get("query", {}).items():
        cmd.append("-d '%s=%s'" % (k, v))

    # redirect output to file
    pfx = req.get("prefix", cls)
    cmd.append("> %s/%s.json" % (TMP_FOLDER, pfx))

    cmd = " ".join(cmd)
    f.write(cmd)

f.write("")
f.write("# Zip result")
f.write("zip -mj ~/aci-vetr-data.zip %s/*.json" % TMP_FOLDER)
f.write("")
f.write("#Cleanup")
f.write("")
f.write("rm -rf %s" % TMP_FOLDER)
f.write("")
f.write("echo Collection complete")
f.write("echo Please provide aci-vetr-data.zip to Cisco for analysis.")
f.close()
os.chmod(SCRIPT_NAME, 0o722)
