package jexl

// Node представляет узел AST.
// Аналог org.apache.commons.jexl3.parser.JexlNode.
type Node interface {
	// Children возвращает дочерние узлы.
	Children() []Node
	// String возвращает строковое представление узла.
	String() string
	// SourceText возвращает исходный текст узла.
	SourceText() string
}

// ScriptNode представляет корневой узел скрипта или выражения.
// Аналог org.apache.commons.jexl3.parser.ASTJexlScript.
type ScriptNode struct {
	info       *Info
	children   []Node
	source     string
	pragmas    map[string]any
	features   *Features
	variables  []string
	parameters []string
}

// NewScriptNode создаёт новый ScriptNode.
func NewScriptNode(info *Info, source string, features *Features) *ScriptNode {
	return &ScriptNode{
		info:     info,
		source:   source,
		features: features,
		children: []Node{},
		pragmas:  make(map[string]any),
	}
}

// Children возвращает дочерние узлы.
func (s *ScriptNode) Children() []Node {
	return s.children
}

// AddChild добавляет дочерний узел.
func (s *ScriptNode) AddChild(child Node) {
	if child != nil {
		s.children = append(s.children, child)
	}
}

// SetChild устанавливает дочерний узел по индексу.
func (s *ScriptNode) SetChild(index int, child Node) {
	if index >= 0 && index < len(s.children) && child != nil {
		s.children[index] = child
	}
}

// String возвращает строковое представление.
func (s *ScriptNode) String() string {
	return s.source
}

// SourceText возвращает исходный текст.
func (s *ScriptNode) SourceText() string {
	return s.source
}

// Info возвращает информацию о скрипте.
func (s *ScriptNode) Info() *Info {
	return s.info
}

// Pragmas возвращает pragma директивы.
func (s *ScriptNode) Pragmas() map[string]any {
	return s.pragmas
}

// SetPragma устанавливает pragma.
func (s *ScriptNode) SetPragma(key string, value any) {
	if s.pragmas == nil {
		s.pragmas = make(map[string]any)
	}
	s.pragmas[key] = value
}

// Features возвращает features скрипта.
func (s *ScriptNode) Features() *Features {
	return s.features
}

// Variables возвращает список переменных.
func (s *ScriptNode) Variables() []string {
	return s.variables
}

// SetVariables устанавливает список переменных.
func (s *ScriptNode) SetVariables(vars []string) {
	s.variables = vars
}

// Parameters возвращает список параметров.
func (s *ScriptNode) Parameters() []string {
	return s.parameters
}

// SetParameters устанавливает список параметров.
func (s *ScriptNode) SetParameters(params []string) {
	s.parameters = params
}

// LiteralNode представляет литерал (число, строка, bool, null).
type LiteralNode struct {
	value  any
	source string
}

// NewLiteralNode создаёт новый LiteralNode.
func NewLiteralNode(value any, source string) *LiteralNode {
	return &LiteralNode{
		value:  value,
		source: source,
	}
}

// Children возвращает пустой список (литералы не имеют детей).
func (l *LiteralNode) Children() []Node {
	return nil
}

// String возвращает строковое представление.
func (l *LiteralNode) String() string {
	return l.source
}

// SourceText возвращает исходный текст.
func (l *LiteralNode) SourceText() string {
	return l.source
}

// Value возвращает значение литерала.
func (l *LiteralNode) Value() any {
	return l.value
}

// IdentifierNode представляет идентификатор (переменную).
type IdentifierNode struct {
	name   string
	source string
}

// NewIdentifierNode создаёт новый IdentifierNode.
func NewIdentifierNode(name, source string) *IdentifierNode {
	return &IdentifierNode{
		name:   name,
		source: source,
	}
}

// Info возвращает информацию об узле (для ошибок).
func (i *IdentifierNode) Info() *Info {
	return NewInfoAt(i.source, 1, 1)
}

// Children возвращает пустой список.
func (i *IdentifierNode) Children() []Node {
	return nil
}

// String возвращает строковое представление.
func (i *IdentifierNode) String() string {
	return i.source
}

// SourceText возвращает исходный текст.
func (i *IdentifierNode) SourceText() string {
	return i.source
}

// Name возвращает имя идентификатора.
func (i *IdentifierNode) Name() string {
	return i.name
}

// BinaryOpNode представляет бинарную операцию.
type BinaryOpNode struct {
	op     string
	left   Node
	right  Node
	source string
}

// NewBinaryOpNode создаёт новый BinaryOpNode.
func NewBinaryOpNode(op string, left, right Node, source string) *BinaryOpNode {
	return &BinaryOpNode{
		op:     op,
		left:   left,
		right:  right,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (b *BinaryOpNode) Children() []Node {
	return []Node{b.left, b.right}
}

// String возвращает строковое представление.
func (b *BinaryOpNode) String() string {
	return b.source
}

// SourceText возвращает исходный текст.
func (b *BinaryOpNode) SourceText() string {
	return b.source
}

// Op возвращает операцию.
func (b *BinaryOpNode) Op() string {
	return b.op
}

// Left возвращает левый операнд.
func (b *BinaryOpNode) Left() Node {
	return b.left
}

// Right возвращает правый операнд.
func (b *BinaryOpNode) Right() Node {
	return b.right
}

// UnaryOpNode представляет унарную операцию.
type UnaryOpNode struct {
	op      string
	operand Node
	source  string
}

// NewUnaryOpNode создаёт новый UnaryOpNode.
func NewUnaryOpNode(op string, operand Node, source string) *UnaryOpNode {
	return &UnaryOpNode{
		op:      op,
		operand: operand,
		source:  source,
	}
}

// Children возвращает дочерние узлы.
func (u *UnaryOpNode) Children() []Node {
	return []Node{u.operand}
}

// String возвращает строковое представление.
func (u *UnaryOpNode) String() string {
	return u.source
}

// SourceText возвращает исходный текст.
func (u *UnaryOpNode) SourceText() string {
	return u.source
}

// Op возвращает операцию.
func (u *UnaryOpNode) Op() string {
	return u.op
}

// Operand возвращает операнд.
func (u *UnaryOpNode) Operand() Node {
	return u.operand
}

// PropertyAccessNode представляет доступ к свойству объекта (obj.prop).
type PropertyAccessNode struct {
	object   Node
	property Node
	source   string
}

// NewPropertyAccessNode создаёт новый PropertyAccessNode.
func NewPropertyAccessNode(object, property Node, source string) *PropertyAccessNode {
	return &PropertyAccessNode{
		object:   object,
		property: property,
		source:   source,
	}
}

// Children возвращает дочерние узлы.
func (p *PropertyAccessNode) Children() []Node {
	return []Node{p.object, p.property}
}

// String возвращает строковое представление.
func (p *PropertyAccessNode) String() string {
	return p.source
}

// SourceText возвращает исходный текст.
func (p *PropertyAccessNode) SourceText() string {
	return p.source
}

// Object возвращает объект.
func (p *PropertyAccessNode) Object() Node {
	return p.object
}

// Property возвращает свойство.
func (p *PropertyAccessNode) Property() Node {
	return p.property
}

// IndexAccessNode представляет доступ к элементу массива/мапы (arr[index]).
type IndexAccessNode struct {
	object Node
	index  Node
	source string
}

// NewIndexAccessNode создаёт новый IndexAccessNode.
func NewIndexAccessNode(object, index Node, source string) *IndexAccessNode {
	return &IndexAccessNode{
		object: object,
		index:  index,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (i *IndexAccessNode) Children() []Node {
	return []Node{i.object, i.index}
}

// String возвращает строковое представление.
func (i *IndexAccessNode) String() string {
	return i.source
}

// SourceText возвращает исходный текст.
func (i *IndexAccessNode) SourceText() string {
	return i.source
}

// Object возвращает объект.
func (i *IndexAccessNode) Object() Node {
	return i.object
}

// Index возвращает индекс.
func (i *IndexAccessNode) Index() Node {
	return i.index
}

// MethodCallNode представляет вызов метода или функции (obj.method(args) или func(args)).
type MethodCallNode struct {
	target Node // может быть nil для функций верхнего уровня
	method Node // имя метода или функции
	args   []Node
	source string
}

// NewMethodCallNode создаёт новый MethodCallNode.
func NewMethodCallNode(target, method Node, args []Node, source string) *MethodCallNode {
	return &MethodCallNode{
		target: target,
		method: method,
		args:   args,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (m *MethodCallNode) Children() []Node {
	children := []Node{}
	if m.target != nil {
		children = append(children, m.target)
	}
	children = append(children, m.method)
	children = append(children, m.args...)
	return children
}

// String возвращает строковое представление.
func (m *MethodCallNode) String() string {
	return m.source
}

// SourceText возвращает исходный текст.
func (m *MethodCallNode) SourceText() string {
	return m.source
}

// Target возвращает целевой объект (может быть nil).
func (m *MethodCallNode) Target() Node {
	return m.target
}

// Method возвращает имя метода или функции.
func (m *MethodCallNode) Method() Node {
	return m.method
}

// Args возвращает аргументы.
func (m *MethodCallNode) Args() []Node {
	return m.args
}

// AssignmentNode представляет присваивание (target = value).
type AssignmentNode struct {
	target Node
	value  Node
	source string
}

// NewAssignmentNode создаёт новый AssignmentNode.
func NewAssignmentNode(target, value Node, source string) *AssignmentNode {
	return &AssignmentNode{
		target: target,
		value:  value,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (a *AssignmentNode) Children() []Node {
	return []Node{a.target, a.value}
}

// String возвращает строковое представление.
func (a *AssignmentNode) String() string {
	return a.source
}

// SourceText возвращает исходный текст.
func (a *AssignmentNode) SourceText() string {
	return a.source
}

// Target возвращает целевой узел присваивания.
func (a *AssignmentNode) Target() Node {
	return a.target
}

// Value возвращает значение присваивания.
func (a *AssignmentNode) Value() Node {
	return a.value
}

// TernaryNode представляет тернарный оператор (condition ? trueExpr : falseExpr).
type TernaryNode struct {
	condition Node
	trueExpr  Node
	falseExpr Node
	source    string
}

// NewTernaryNode создаёт новый TernaryNode.
func NewTernaryNode(condition, trueExpr, falseExpr Node, source string) *TernaryNode {
	return &TernaryNode{
		condition: condition,
		trueExpr:  trueExpr,
		falseExpr: falseExpr,
		source:    source,
	}
}

// Children возвращает дочерние узлы.
func (t *TernaryNode) Children() []Node {
	return []Node{t.condition, t.trueExpr, t.falseExpr}
}

// String возвращает строковое представление.
func (t *TernaryNode) String() string {
	return t.source
}

// SourceText возвращает исходный текст.
func (t *TernaryNode) SourceText() string {
	return t.source
}

// Condition возвращает условие.
func (t *TernaryNode) Condition() Node {
	return t.condition
}

// TrueExpr возвращает выражение для true.
func (t *TernaryNode) TrueExpr() Node {
	return t.trueExpr
}

// FalseExpr возвращает выражение для false.
func (t *TernaryNode) FalseExpr() Node {
	return t.falseExpr
}

// RangeNode представляет range оператор (left .. right).
type RangeNode struct {
	left   Node
	right  Node
	source string
}

// NewRangeNode создаёт новый RangeNode.
func NewRangeNode(left, right Node, source string) *RangeNode {
	return &RangeNode{
		left:   left,
		right:  right,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (r *RangeNode) Children() []Node {
	return []Node{r.left, r.right}
}

// String возвращает строковое представление.
func (r *RangeNode) String() string {
	return r.source
}

// SourceText возвращает исходный текст.
func (r *RangeNode) SourceText() string {
	return r.source
}

// Left возвращает левый операнд.
func (r *RangeNode) Left() Node {
	return r.left
}

// Right возвращает правый операнд.
func (r *RangeNode) Right() Node {
	return r.right
}

// ElvisNode представляет Elvis оператор (expr ?: defaultExpr).
type ElvisNode struct {
	expr        Node
	defaultExpr Node
	source      string
}

// NewElvisNode создаёт новый ElvisNode.
func NewElvisNode(expr, defaultExpr Node, source string) *ElvisNode {
	return &ElvisNode{
		expr:        expr,
		defaultExpr: defaultExpr,
		source:      source,
	}
}

// Children возвращает дочерние узлы.
func (e *ElvisNode) Children() []Node {
	return []Node{e.expr, e.defaultExpr}
}

// String возвращает строковое представление.
func (e *ElvisNode) String() string {
	return e.source
}

// SourceText возвращает исходный текст.
func (e *ElvisNode) SourceText() string {
	return e.source
}

// Expr возвращает выражение.
func (e *ElvisNode) Expr() Node {
	return e.expr
}

// DefaultExpr возвращает выражение по умолчанию.
func (e *ElvisNode) DefaultExpr() Node {
	return e.defaultExpr
}

// ArrayLiteralNode представляет литерал массива [1, 2, 3].
type ArrayLiteralNode struct {
	elements []Node
	source   string
}

// NewArrayLiteralNode создаёт новый ArrayLiteralNode.
func NewArrayLiteralNode(elements []Node, source string) *ArrayLiteralNode {
	return &ArrayLiteralNode{
		elements: elements,
		source:   source,
	}
}

// Children возвращает дочерние узлы.
func (a *ArrayLiteralNode) Children() []Node {
	return a.elements
}

// String возвращает строковое представление.
func (a *ArrayLiteralNode) String() string {
	return a.source
}

// SourceText возвращает исходный текст.
func (a *ArrayLiteralNode) SourceText() string {
	return a.source
}

// Elements возвращает элементы массива.
func (a *ArrayLiteralNode) Elements() []Node {
	return a.elements
}

// MapEntry представляет пару ключ-значение в мапе.
type MapEntry struct {
	Key   Node
	Value Node
}

// MapLiteralNode представляет литерал мапы {key: value, ...}.
type MapLiteralNode struct {
	entries []MapEntry
	source  string
}

// NewMapLiteralNode создаёт новый MapLiteralNode.
func NewMapLiteralNode(entries []MapEntry, source string) *MapLiteralNode {
	return &MapLiteralNode{
		entries: entries,
		source:  source,
	}
}

// Children возвращает дочерние узлы.
func (m *MapLiteralNode) Children() []Node {
	children := make([]Node, 0, len(m.entries)*2)
	for _, entry := range m.entries {
		children = append(children, entry.Key, entry.Value)
	}
	return children
}

// String возвращает строковое представление.
func (m *MapLiteralNode) String() string {
	return m.source
}

// SourceText возвращает исходный текст.
func (m *MapLiteralNode) SourceText() string {
	return m.source
}

// Entries возвращает записи мапы.
func (m *MapLiteralNode) Entries() []MapEntry {
	return m.entries
}

// SetLiteralNode представляет литерал множества {1, 2, 3}.
type SetLiteralNode struct {
	elements []Node
	source   string
}

// NewSetLiteralNode создаёт новый SetLiteralNode.
func NewSetLiteralNode(elements []Node, source string) *SetLiteralNode {
	return &SetLiteralNode{
		elements: elements,
		source:   source,
	}
}

// Children возвращает дочерние узлы.
func (s *SetLiteralNode) Children() []Node {
	return s.elements
}

// String возвращает строковое представление.
func (s *SetLiteralNode) String() string {
	return s.source
}

// SourceText возвращает исходный текст.
func (s *SetLiteralNode) SourceText() string {
	return s.source
}

// Elements возвращает элементы множества.
func (s *SetLiteralNode) Elements() []Node {
	return s.elements
}

// IfNode представляет условный оператор if/else.
type IfNode struct {
	condition Node
	thenBranch Node
	elseBranch Node
	source    string
}

// NewIfNode создаёт новый IfNode.
func NewIfNode(condition, thenBranch, elseBranch Node, source string) *IfNode {
	return &IfNode{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
		source:     source,
	}
}

// Children возвращает дочерние узлы.
func (i *IfNode) Children() []Node {
	children := []Node{i.condition, i.thenBranch}
	if i.elseBranch != nil {
		children = append(children, i.elseBranch)
	}
	return children
}

// String возвращает строковое представление.
func (i *IfNode) String() string {
	return i.source
}

// SourceText возвращает исходный текст.
func (i *IfNode) SourceText() string {
	return i.source
}

// Condition возвращает условие.
func (i *IfNode) Condition() Node {
	return i.condition
}

// ThenBranch возвращает ветку then.
func (i *IfNode) ThenBranch() Node {
	return i.thenBranch
}

// ElseBranch возвращает ветку else (может быть nil).
func (i *IfNode) ElseBranch() Node {
	return i.elseBranch
}

// ForNode представляет цикл for (init; condition; step) body.
type ForNode struct {
	init      Node
	condition Node
	step      Node
	body      Node
	source    string
}

// NewForNode создаёт новый ForNode.
func NewForNode(init, condition, step, body Node, source string) *ForNode {
	return &ForNode{
		init:      init,
		condition: condition,
		step:      step,
		body:      body,
		source:    source,
	}
}

// Children возвращает дочерние узлы.
func (f *ForNode) Children() []Node {
	children := []Node{}
	if f.init != nil {
		children = append(children, f.init)
	}
	if f.condition != nil {
		children = append(children, f.condition)
	}
	if f.step != nil {
		children = append(children, f.step)
	}
	children = append(children, f.body)
	return children
}

// String возвращает строковое представление.
func (f *ForNode) String() string {
	return f.source
}

// SourceText возвращает исходный текст.
func (f *ForNode) SourceText() string {
	return f.source
}

// Init возвращает инициализацию.
func (f *ForNode) Init() Node {
	return f.init
}

// Condition возвращает условие.
func (f *ForNode) Condition() Node {
	return f.condition
}

// Step возвращает шаг.
func (f *ForNode) Step() Node {
	return f.step
}

// Body возвращает тело цикла.
func (f *ForNode) Body() Node {
	return f.body
}

// ForeachNode представляет цикл foreach (var x : items) body.
type ForeachNode struct {
	variable Node
	items    Node
	body     Node
	source   string
}

// NewForeachNode создаёт новый ForeachNode.
func NewForeachNode(variable, items, body Node, source string) *ForeachNode {
	return &ForeachNode{
		variable: variable,
		items:    items,
		body:     body,
		source:   source,
	}
}

// Children возвращает дочерние узлы.
func (f *ForeachNode) Children() []Node {
	return []Node{f.variable, f.items, f.body}
}

// String возвращает строковое представление.
func (f *ForeachNode) String() string {
	return f.source
}

// SourceText возвращает исходный текст.
func (f *ForeachNode) SourceText() string {
	return f.source
}

// Variable возвращает переменную цикла.
func (f *ForeachNode) Variable() Node {
	return f.variable
}

// Items возвращает коллекцию для итерации.
func (f *ForeachNode) Items() Node {
	return f.items
}

// Body возвращает тело цикла.
func (f *ForeachNode) Body() Node {
	return f.body
}

// WhileNode представляет цикл while (condition) body.
type WhileNode struct {
	condition Node
	body      Node
	source    string
}

// NewWhileNode создаёт новый WhileNode.
func NewWhileNode(condition, body Node, source string) *WhileNode {
	return &WhileNode{
		condition: condition,
		body:      body,
		source:    source,
	}
}

// Children возвращает дочерние узлы.
func (w *WhileNode) Children() []Node {
	return []Node{w.condition, w.body}
}

// String возвращает строковое представление.
func (w *WhileNode) String() string {
	return w.source
}

// SourceText возвращает исходный текст.
func (w *WhileNode) SourceText() string {
	return w.source
}

// Condition возвращает условие.
func (w *WhileNode) Condition() Node {
	return w.condition
}

// Body возвращает тело цикла.
func (w *WhileNode) Body() Node {
	return w.body
}

// DoWhileNode представляет цикл do body while (condition).
type DoWhileNode struct {
	condition Node
	body      Node
	source    string
}

// NewDoWhileNode создаёт новый DoWhileNode.
func NewDoWhileNode(condition, body Node, source string) *DoWhileNode {
	return &DoWhileNode{
		condition: condition,
		body:      body,
		source:    source,
	}
}

// Children возвращает дочерние узлы.
func (d *DoWhileNode) Children() []Node {
	return []Node{d.body, d.condition}
}

// String возвращает строковое представление.
func (d *DoWhileNode) String() string {
	return d.source
}

// SourceText возвращает исходный текст.
func (d *DoWhileNode) SourceText() string {
	return d.source
}

// Condition возвращает условие.
func (d *DoWhileNode) Condition() Node {
	return d.condition
}

// Body возвращает тело цикла.
func (d *DoWhileNode) Body() Node {
	return d.body
}

// BlockNode представляет блок кода { statements }.
type BlockNode struct {
	statements []Node
	source     string
}

// NewBlockNode создаёт новый BlockNode.
func NewBlockNode(statements []Node, source string) *BlockNode {
	return &BlockNode{
		statements: statements,
		source:     source,
	}
}

// Children возвращает дочерние узлы.
func (b *BlockNode) Children() []Node {
	return b.statements
}

// String возвращает строковое представление.
func (b *BlockNode) String() string {
	return b.source
}

// SourceText возвращает исходный текст.
func (b *BlockNode) SourceText() string {
	return b.source
}

// Statements возвращает statements блока.
func (b *BlockNode) Statements() []Node {
	return b.statements
}

// BreakNode представляет оператор break.
type BreakNode struct {
	source string
}

// NewBreakNode создаёт новый BreakNode.
func NewBreakNode(source string) *BreakNode {
	return &BreakNode{source: source}
}

// Children возвращает пустой список.
func (b *BreakNode) Children() []Node {
	return nil
}

// String возвращает строковое представление.
func (b *BreakNode) String() string {
	return b.source
}

// SourceText возвращает исходный текст.
func (b *BreakNode) SourceText() string {
	return b.source
}

// ContinueNode представляет оператор continue.
type ContinueNode struct {
	source string
}

// NewContinueNode создаёт новый ContinueNode.
func NewContinueNode(source string) *ContinueNode {
	return &ContinueNode{source: source}
}

// Children возвращает пустой список.
func (c *ContinueNode) Children() []Node {
	return nil
}

// String возвращает строковое представление.
func (c *ContinueNode) String() string {
	return c.source
}

// SourceText возвращает исходный текст.
func (c *ContinueNode) SourceText() string {
	return c.source
}

// ReturnNode представляет оператор return.
type ReturnNode struct {
	value  Node
	source string
}

// NewReturnNode создаёт новый ReturnNode.
func NewReturnNode(value Node, source string) *ReturnNode {
	return &ReturnNode{
		value:  value,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (r *ReturnNode) Children() []Node {
	if r.value == nil {
		return nil
	}
	return []Node{r.value}
}

// String возвращает строковое представление.
func (r *ReturnNode) String() string {
	return r.source
}

// SourceText возвращает исходный текст.
func (r *ReturnNode) SourceText() string {
	return r.source
}

// Value возвращает возвращаемое значение (может быть nil).
func (r *ReturnNode) Value() Node {
	return r.value
}

// VarNode представляет объявление переменной var x или var x = value.
type VarNode struct {
	name   *IdentifierNode
	value  Node
	source string
}

// NewVarNode создаёт новый VarNode.
func NewVarNode(name *IdentifierNode, value Node, source string) *VarNode {
	return &VarNode{
		name:   name,
		value:  value,
		source: source,
	}
}

// Children возвращает дочерние узлы.
func (v *VarNode) Children() []Node {
	if v.value != nil {
		return []Node{v.name, v.value}
	}
	return []Node{v.name}
}

// String возвращает строковое представление.
func (v *VarNode) String() string {
	return v.source
}

// SourceText возвращает исходный текст.
func (v *VarNode) SourceText() string {
	return v.source
}

// Name возвращает имя переменной.
func (v *VarNode) Name() *IdentifierNode {
	return v.name
}

// Value возвращает значение переменной (может быть nil).
func (v *VarNode) Value() Node {
	return v.value
}

// LambdaNode представляет lambda функцию (x, y) -> x + y или (x, y) => x + y.
// Аналог org.apache.commons.jexl3.parser.ASTJexlLambda.
type LambdaNode struct {
	parameters []*IdentifierNode
	body       Node
	source     string
}

// NewLambdaNode создаёт новый LambdaNode.
func NewLambdaNode(parameters []*IdentifierNode, body Node, source string) *LambdaNode {
	return &LambdaNode{
		parameters: parameters,
		body:       body,
		source:     source,
	}
}

// Children возвращает дочерние узлы (параметры и тело).
func (l *LambdaNode) Children() []Node {
	children := make([]Node, 0, len(l.parameters)+1)
	for _, param := range l.parameters {
		children = append(children, param)
	}
	if l.body != nil {
		children = append(children, l.body)
	}
	return children
}

// String возвращает строковое представление.
func (l *LambdaNode) String() string {
	return l.source
}

// SourceText возвращает исходный текст.
func (l *LambdaNode) SourceText() string {
	return l.source
}

// Parameters возвращает список параметров.
func (l *LambdaNode) Parameters() []*IdentifierNode {
	return l.parameters
}

// Body возвращает тело lambda функции.
func (l *LambdaNode) Body() Node {
	return l.body
}
