#!python
import json
import os

TMP_FOLDER = "/tmp/vetr-collector"
SCRIPT_NAME = "vetr-collector.sh"

with open("./pkg/req/reqs.json") as f:
    reqs = json.loads(f.read())

f = open(SCRIPT_NAME, "w")

f.writelines(
    [
        "#!/bin/bash\n",
        "\n",
        "rm -rf %s > /dev/null\n" % TMP_FOLDER,
        "mkdir %s\n" % TMP_FOLDER,
        "\n",
        "# Fetch data from the API\n",
    ]
)

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
    f.writelines([cmd + "\n"])

f.writelines(
    [
        "\n",
        "# Zip result\n",
        "zip -mj ~/aci-vetr-data.zip %s/*.json\n" % TMP_FOLDER,
        "\n",
        "# Cleanup\n",
        "\n",
        "rm -rf %s\n" % TMP_FOLDER,
        "\n",
        "echo Collection complete\n",
        "echo Output writen to ~/aci-vetr-data.zip, i.e. user home folder\n",
        "echo Please provide aci-vetr-data.zip to Cisco for analysis.\n",
    ]
)
f.close()
os.chmod(SCRIPT_NAME, 0o722)
