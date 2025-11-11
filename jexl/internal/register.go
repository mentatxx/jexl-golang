package internal

import "github.com/mentatxx/jexl-golang/jexl"

func init() {
	jexl.RegisterEngineBuilder(NewEngine)
}

// ExportNewEngine экспортирует NewEngine для использования в тестах
// через build tags. В обычной сборке эта функция не используется.
var ExportNewEngine = NewEngine
