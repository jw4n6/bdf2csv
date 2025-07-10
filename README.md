
## Overview

`bdf2csv` converts [Unix-like Artifacts Collector (UAC)](https://github.com/tclahr/uac) Linux bodyfiles into CSV format with human-readable timestamps in UTC.

## Usage
```bash
bdf2csv -i <bodyfile> -o <csvfile> [options]

Options:
  -i string         Input bodyfile path (required)
  -o string         Output CSV file path (required)
  -e                Keep timestamps in epoch format only (default is human-readable)
  -v                Show version and exit
  -h                Show this help message

Example:
  bdf2csv -i bodyfile.txt -o bodyfile.csv
  bdf2csv -i bodyfile.txt -o bodyfile.csv -e
```

## Installation

### Using go install (Recommended)
```bash
go install github.com/jw4n6/bdf2csv@latest
```

### Manual Installation
```bash
git clone https://github.com/jw4n6/bdf2csv.git
cd bdf2csv
go build -o bdf2csv
```

## Format

**Column Structure:**
1. **0** - Placeholder (always "0" in UAC bodyfiles)
2. **Name** - File/directory path
3. **Inode** - Inode number
4. **Mode** - File permissions (e.g., "drwxr-xr-x")
5. **UID** - User ID
6. **GID** - Group ID
7. **Size** - File size in bytes
8. **ATime** - Access time
9. **MTime** - Modify time
10. **CTime** - Change time
11. **CrTime** - Creation time

NOTE: in order to use the `bdf2csv` binary, make sure the GOBIN is part of your PATH env variable:

```bash
$ export GOBIN=`go env GOPATH`/bin
$ export PATH=$PATH:$GOBIN
```
