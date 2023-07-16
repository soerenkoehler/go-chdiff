package common

import (
	"time"

	"github.com/lestrrat-go/strftime"
)

type Location struct {
	Path string
	Time time.Time
}

var LocationTimeFormat, _ = strftime.New("%Y-%m-%d %H-%M-%S")
