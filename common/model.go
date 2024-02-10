package common

import (
	"time"

	"github.com/lestrrat-go/strftime"
)

type Location struct {
	Path string
	Time time.Time
}

var Config struct {
	Exclude struct {
		Absolute     []string
		RootRelative []string
		Anywhere     []string
	}
}

var LocationTimeFormat, _ = strftime.New("%Y-%m-%d %H-%M-%S")
