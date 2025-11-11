package jexl

// Sandbox описывает ограничения на доступ к методам/свойствам.
// Аналог org.apache.commons.jexl3.introspection.JexlSandbox.
type Sandbox struct {
	whitelist map[string]map[string]bool
	blacklist map[string]map[string]bool
}

// NewSandbox создаёт пустой Sandbox.
func NewSandbox() *Sandbox {
	return &Sandbox{
		whitelist: make(map[string]map[string]bool),
		blacklist: make(map[string]map[string]bool),
	}
}

// Allow добавляет элемент белого списка.
func (s *Sandbox) Allow(className, member string) {
	if s == nil {
		return
	}
	if _, ok := s.whitelist[className]; !ok {
		s.whitelist[className] = make(map[string]bool)
	}
	s.whitelist[className][member] = true
}

// Deny добавляет элемент чёрного списка.
func (s *Sandbox) Deny(className, member string) {
	if s == nil {
		return
	}
	if _, ok := s.blacklist[className]; !ok {
		s.blacklist[className] = make(map[string]bool)
	}
	s.blacklist[className][member] = true
}

// Allowed проверяет, разрешён ли доступ.
func (s *Sandbox) Allowed(className, member string) bool {
	if s == nil {
		return true
	}
	if denies, ok := s.blacklist[className]; ok && denies[member] {
		return false
	}
	if allows, ok := s.whitelist[className]; ok {
		return allows[member]
	}
	return true
}
