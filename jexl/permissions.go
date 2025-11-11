package jexl

// Permissions контролируют доступный набор пакетов и классов.
// Структура служит портом org.apache.commons.jexl3.introspection.JexlPermissions.
type Permissions struct {
	allowed []string
	denied  []string
}

// NewPermissions создаёт новый набор разрешений.
func NewPermissions(allowed, denied []string) *Permissions {
	return &Permissions{
		allowed: append([]string(nil), allowed...),
		denied:  append([]string(nil), denied...),
	}
}

// Allowed возвращает список разрешённых шаблонов.
func (p *Permissions) Allowed() []string {
	if p == nil {
		return nil
	}
	return append([]string(nil), p.allowed...)
}

// Denied возвращает список запрещённых шаблонов.
func (p *Permissions) Denied() []string {
	if p == nil {
		return nil
	}
	return append([]string(nil), p.denied...)
}

var (
	// PermissionsRestricted соответствует набору RESTRICTED в Java-версии.
	PermissionsRestricted = NewPermissions(
		[]string{
			"java.lang", "java.util", "java.math",
			"java.time", "jexl", "jexl3",
		},
		[]string{
			"java.lang.System",
			"java.lang.Runtime",
			"java.lang.Process",
			"java.lang.ProcessBuilder",
			"java.lang.Thread",
			"java.lang.ClassLoader",
		},
	)

	// PermissionsUnrestricted соответствует UNRESTRICTED.
	PermissionsUnrestricted = NewPermissions(nil, nil)
)

// Clone создаёт копию набора разрешений.
func (p *Permissions) Clone() *Permissions {
	if p == nil {
		return nil
	}
	return NewPermissions(p.allowed, p.denied)
}
