package models

import (
	"io"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func SetDebug(writer io.Writer) {
	boil.DebugMode = true
	boil.DebugWriter = writer
}
