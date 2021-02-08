package resource

// Usage : Description of arguments and options.
const Usage = `ChecksumDiff (build ${VERSION})

Usage:
    chdiff c PATH [-f FILE] [-m MODE]
    chdiff v PATH [-f FILE] [-m MODE]
    chdiff (-h | --help | --version)

Commands:
    c  Create digest file for directory PATH.
    v  Verify digest file for directory PATH.

Options:
    -f FILE    Use the given digest file.
               Default: PATH/.chdiff.<MODE>

    -m MODE    The checksum algorithm used.
               Values: SHA256, SHA512
               Default: SHA256

    -h --help  Show help.
    --version  Show version.`
