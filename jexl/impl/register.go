package impl

import (
	"github.com/mentatxx/jexl-golang/jexl"
	"github.com/mentatxx/jexl-golang/jexl/internal"
)

func init() {
	jexl.RegisterEngineBuilder(internal.NewEngine)
}
