module github.com/soerenkoehler/go-chdiff

go 1.18

require (
	github.com/alecthomas/kong v0.4.1
)

require (
	github.com/google/go-cmp v0.5.7
	github.com/soerenkoehler/go-testutils v0.0.2
)

require github.com/pkg/errors v0.9.1 // indirect

replace github.com/soerenkoehler/go-testutils => ../go-testutils
