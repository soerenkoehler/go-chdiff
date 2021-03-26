module github.com/soerenkoehler/chdiff-go/test

go 1.16

require (
	github.com/soerenkoehler/chdiff-go/main v0.0.0-00010101000000-000000000000
	github.com/soerenkoehler/go-testutils v0.0.0-20210315022031-b4c6c1bf5cbd
)

replace github.com/soerenkoehler/chdiff-go/main => ../main

replace github.com/soerenkoehler/go-testutils => ../../go-testutils
