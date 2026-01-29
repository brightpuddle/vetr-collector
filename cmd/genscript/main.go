// Command genscript generates the vetr-collector.sh shell script
// from the embedded request data.
package main

import (
"fmt"
"os"
"path/filepath"

"collector/pkg/req"
)

const (
tmpFolder  = "/tmp/vetr-collector"
)

func main() {
// Determine output path - should be at repo root
// When run via go generate from pkg/req, we need to go up two directories
scriptPath := "vetr-collector.sh"
if _, err := os.Stat("../../go.mod"); err == nil {
// We're in a subdirectory (e.g., pkg/req), write to repo root
scriptPath = "../../vetr-collector.sh"
}

// Get absolute path
absPath, err := filepath.Abs(scriptPath)
if err != nil {
fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
os.Exit(1)
}

// Open output file
f, err := os.Create(absPath)
if err != nil {
fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
os.Exit(1)
}
defer f.Close()

// Write script header
fmt.Fprintln(f, "#!/bin/bash")
fmt.Fprintln(f, "")
fmt.Fprintf(f, "rm -rf %s > /dev/null\n", tmpFolder)
fmt.Fprintf(f, "mkdir %s\n", tmpFolder)
fmt.Fprintln(f, "")
fmt.Fprintln(f, "# Fetch data from the API")

// Get requests
reqs, err := req.GetRequests()
if err != nil {
fmt.Fprintf(os.Stderr, "Error getting requests: %v\n", err)
os.Exit(1)
}

// Write icurl commands for each request
for _, r := range reqs {
cmd := fmt.Sprintf("icurl -kG https://localhost/api/class/%s.json", r.Class)

// Add query parameters if present
for k, v := range r.Query {
cmd += fmt.Sprintf(" -d '%s=%s'", k, v)
}

// Add output redirection
cmd += fmt.Sprintf(" > %s/%s.json", tmpFolder, r.Class)

fmt.Fprintln(f, cmd)
}

// Write script footer
fmt.Fprintln(f, "")
fmt.Fprintln(f, "# Zip result")
fmt.Fprintf(f, "zip -mj ~/aci-vetr-data.zip %s/*.json\n", tmpFolder)
fmt.Fprintln(f, "")
fmt.Fprintln(f, "# Cleanup")
fmt.Fprintln(f, "")
fmt.Fprintf(f, "rm -rf %s\n", tmpFolder)
fmt.Fprintln(f, "")
fmt.Fprintln(f, "echo Collection complete")
fmt.Fprintln(f, "echo Output writen to ~/aci-vetr-data.zip, i.e. user home folder")
fmt.Fprintln(f, "echo Please provide aci-vetr-data.zip to Cisco for analysis.")

// Make the script executable
if err := os.Chmod(absPath, 0o755); err != nil {
fmt.Fprintf(os.Stderr, "Error making script executable: %v\n", err)
os.Exit(1)
}

fmt.Printf("Generated %s successfully\n", absPath)
}
