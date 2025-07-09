
## Overview

`bdf2csv` converts [Unix-like Artifacts Collector (UAC)](https://github.com/tclahr/uac) Linux bodyfiles into CSV format with human-readable timestamps in UTC.

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
