package resource

// Usage : Description of arguments and options.
const Usage = `ChecksumDiff (build ${VERSION})

Usage:
  chdiff (c|v) [PATH] [-f FILE] [-m MODE]
  chdiff (-h | --help | --version)
  chdiff

Commands:
  c  Create digest file for PATH.
  v  Verify digest file for PATH.

Options:
 -f FILE  Use the given digest file. The default is
          * PATH/.chdiff.MODE if PATH is a directory
          * PATH.chdiff.MODE if PATH is a file

 -m MODE  The checksum algorithm used SHA256, SHA512 [default: SHA256]
 -h --help  Show help.
 --version  Show version.

Remarks:
    * Calling chdiff without parameters is equivalent to: chdiff v .`
