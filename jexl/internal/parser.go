package internal

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mentatxx/jexl-golang/jexl"
)

// Parser представляет парсер JEXL.
type Parser interface {
	ParseExpression(info *jexl.Info, source string, features *jexl.Features) (*jexl.ScriptNode, error)
	ParseScript(info *jexl.Info, source string, features *jexl.Features, names []string) (*jexl.ScriptNode, error)
}

// defaultParser простая реализация парсера.
type defaultParser struct{}

// NewDefaultParser создаёт парсер по умолчанию.
func NewDefaultParser() Parser {
	return &defaultParser{}
}

// ParseExpression парсит выражение.
func (p *defaultParser) ParseExpression(info *jexl.Info, source string, features *jexl.Features) (*jexl.ScriptNode, error) {
	builder := newSimpleParser(info, source, features)
	node, err := builder.parseExpression(0)
	if err != nil {
		return nil, err
	}
	if err := builder.expect(tokenEOF); err != nil {
		return nil, err
	}

	ast := jexl.NewScriptNode(info, strings.TrimSpace(source), features)
	ast.AddChild(node)
	return ast, nil
}

// ParseScript парсит скрипт.
func (p *defaultParser) ParseScript(info *jexl.Info, source string, features *jexl.Features, names []string) (*jexl.ScriptNode, error) {
	builder := newSimpleParser(info, source, features)
	ast := jexl.NewScriptNode(info, strings.TrimSpace(source), features)

	for {
		if builder.peek().typ == tokenEOF {
			break
		}
		// Пробуем распарсить statement, если не получается - expression
		node, err := builder.parseStatement()
		if err != nil {
			return nil, err
		}
		if node == nil {
			// Если statement не распарсился, пробуем expression
			node, err = builder.parseExpression(0)
			if err != nil {
				return nil, err
			}
		}
		ast.AddChild(node)
		if builder.match(tokenSemicolon) {
			continue
		}
		if builder.peek().typ != tokenEOF {
			return nil, builder.errorf("unexpected token %s", builder.peek().literal)
		}
		break
	}

	if err := builder.expect(tokenEOF); err != nil {
		return nil, err
	}

	ast.SetParameters(names)
	return ast, nil
}

// simpleParser реализует примитивный Pratt-парсер для базовой арифметики.
type simpleParser struct {
	info     *jexl.Info
	source   string
	features *jexl.Features
	tokens   []token
	pos      int
}

func newSimpleParser(info *jexl.Info, source string, features *jexl.Features) *simpleParser {
	lexer := newLexer(source)
	tokens := lexer.lex()
	return &simpleParser{
		info:     info,
		source:   source,
		features: features,
		tokens:   tokens,
	}
}

func (p *simpleParser) parseExpression(precedence int) (jexl.Node, error) {
	tok := p.next()
	var left jexl.Node

	switch tok.typ {
	case tokenPlus, tokenMinus, tokenBang, tokenTilde:
		operand, err := p.parseExpression(prefixPrecedence)
		if err != nil {
			return nil, err
		}
		left = jexl.NewUnaryOpNode(tok.literal, operand, tok.literal+operand.SourceText())
	case tokenPlusPlus:
		// Префиксный инкремент: ++x
		operand, err := p.parseExpression(prefixPrecedence)
		if err != nil {
			return nil, err
		}
		if !isAssignableTarget(operand) {
			return nil, p.errorf("operand of increment must be assignable")
		}
		// ++x становится x = x + 1
		one := jexl.NewLiteralNode(int64(1), "1")
		addNode := jexl.NewBinaryOpNode("+", operand, one, fmt.Sprintf("%s + 1", operand.SourceText()))
		left = jexl.NewAssignmentNode(operand, addNode, fmt.Sprintf("++%s", operand.SourceText()))
	case tokenMinusMinus:
		// Префиксный декремент: --x
		operand, err := p.parseExpression(prefixPrecedence)
		if err != nil {
			return nil, err
		}
		if !isAssignableTarget(operand) {
			return nil, p.errorf("operand of decrement must be assignable")
		}
		// --x становится x = x - 1
		one := jexl.NewLiteralNode(int64(1), "1")
		subNode := jexl.NewBinaryOpNode("-", operand, one, fmt.Sprintf("%s - 1", operand.SourceText()))
		left = jexl.NewAssignmentNode(operand, subNode, fmt.Sprintf("--%s", operand.SourceText()))
	case tokenEmpty:
		// Оператор empty: empty x
		operand, err := p.parseExpression(prefixPrecedence)
		if err != nil {
			return nil, err
		}
		// Создаём вызов функции empty(operand)
		emptyIdent := jexl.NewIdentifierNode("empty", "empty")
		args := []jexl.Node{operand}
		left = jexl.NewMethodCallNode(nil, emptyIdent, args, fmt.Sprintf("empty(%s)", operand.SourceText()))
	case tokenNot:
		// Оператор not: not empty x или not x
		next := p.peek()
		if next.typ == tokenEmpty {
			// not empty x
			p.next() // consume 'empty'
			operand, err := p.parseExpression(prefixPrecedence)
			if err != nil {
				return nil, err
			}
			// Создаём вызов функции empty(operand) и инвертируем результат
			emptyIdent := jexl.NewIdentifierNode("empty", "empty")
			args := []jexl.Node{operand}
			emptyCall := jexl.NewMethodCallNode(nil, emptyIdent, args, fmt.Sprintf("empty(%s)", operand.SourceText()))
			left = jexl.NewUnaryOpNode("!", emptyCall, fmt.Sprintf("not empty %s", operand.SourceText()))
		} else {
			// not x - обычное логическое отрицание
			operand, err := p.parseExpression(prefixPrecedence)
			if err != nil {
				return nil, err
			}
			left = jexl.NewUnaryOpNode("!", operand, fmt.Sprintf("not %s", operand.SourceText()))
		}
	case tokenNumber:
		value, err := parseNumberLiteral(tok.literal)
		if err != nil {
			return nil, p.errorf("invalid number %s", tok.literal)
		}
		left = jexl.NewLiteralNode(value, tok.literal)
	case tokenString:
		val, err := parseStringLiteral(tok.literal)
		if err != nil {
			return nil, err
		}
		left = jexl.NewLiteralNode(val, tok.literal)
	case tokenIdent:
		left = jexl.NewIdentifierNode(tok.literal, tok.literal)
	case tokenBool:
		left = jexl.NewLiteralNode(tok.value, tok.literal)
	case tokenNull:
		left = jexl.NewLiteralNode(nil, tok.literal)
	case tokenLBracket:
		// Массив: [1, 2, 3]
		return p.parseArrayLiteral()
	case tokenLBrace:
		// Мапа или множество: {key: value} или {1, 2, 3}
		return p.parseMapOrSetLiteral()
	case tokenLParen:
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		if err := p.expect(tokenRParen); err != nil {
			return nil, err
		}
		left = expr
	default:
		return nil, p.errorf("unexpected token %s", tok.literal)
	}

	// Обрабатываем постфиксные операции (вызовы методов, доступ к свойствам, индексация)
	for {
		next := p.peek()
		if next.typ == tokenEOF || next.typ == tokenRParen || next.typ == tokenSemicolon {
			break
		}

		// Вызов метода или функции: expr(...)
		if next.typ == tokenLParen {
			var err error
			left, err = p.parseCall(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Доступ к свойству: expr.prop
		if next.typ == tokenDot {
			p.next() // consume '.'
			if p.peek().typ != tokenIdent {
				return nil, p.errorf("expected identifier after '.'")
			}
			prop := p.next()
			propNode := jexl.NewIdentifierNode(prop.literal, prop.literal)
			left = jexl.NewPropertyAccessNode(left, propNode, fmt.Sprintf("%s.%s", left.SourceText(), prop.literal))
			continue
		}

		// Side-effect операторы: expr += value, expr -= value, и т.д.
		if next.typ == tokenPlusEqual || next.typ == tokenMinusEqual ||
			next.typ == tokenStarEqual || next.typ == tokenSlashEqual ||
			next.typ == tokenPercentEqual {
			if !isAssignableTarget(left) {
				return nil, p.errorf("left-hand side of assignment is not assignable")
			}
			op := p.next() // consume оператор
			right, err := p.parseExpression(infixPrecedence(next.typ))
			if err != nil {
				return nil, err
			}
			// Преобразуем side-effect оператор в обычное присваивание с операцией
			// x += 3 становится x = x + 3
			var opSymbol string
			switch op.typ {
			case tokenPlusEqual:
				opSymbol = "+"
			case tokenMinusEqual:
				opSymbol = "-"
			case tokenStarEqual:
				opSymbol = "*"
			case tokenSlashEqual:
				opSymbol = "/"
			case tokenPercentEqual:
				opSymbol = "%"
			}
			// Создаём бинарную операцию left op right
			opNode := jexl.NewBinaryOpNode(opSymbol, left, right, fmt.Sprintf("%s %s %s", left.SourceText(), opSymbol, right.SourceText()))
			// Создаём присваивание left = (left op right)
			source := fmt.Sprintf("%s %s %s", left.SourceText(), op.literal, right.SourceText())
			left = jexl.NewAssignmentNode(left, opNode, source)
			continue
		}

		// Присваивание: expr = value
		if next.typ == tokenEqual {
			if !isAssignableTarget(left) {
				return nil, p.errorf("left-hand side of assignment is not assignable")
			}
			p.next() // consume '='
			right, err := p.parseExpression(infixPrecedence(tokenEqual))
			if err != nil {
				return nil, err
			}
			source := fmt.Sprintf("%s = %s", left.SourceText(), right.SourceText())
			left = jexl.NewAssignmentNode(left, right, source)
			continue
		}

		// Индексация: expr[index]
		if next.typ == tokenLBracket {
			p.next() // consume '['
			index, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if err := p.expect(tokenRBracket); err != nil {
				return nil, err
			}
			left = jexl.NewIndexAccessNode(left, index, fmt.Sprintf("%s[%s]", left.SourceText(), index.SourceText()))
			continue
		}

		// Тернарный оператор: condition ? trueExpr : falseExpr
		if next.typ == tokenQuestion && precedence < ternaryPrecedence {
			p.next() // consume '?'
			trueExpr, err := p.parseExpression(ternaryPrecedence)
			if err != nil {
				return nil, err
			}
			if err := p.expect(tokenColon); err != nil {
				return nil, err
			}
			falseExpr, err := p.parseExpression(ternaryPrecedence)
			if err != nil {
				return nil, err
			}
			source := fmt.Sprintf("%s ? %s : %s", left.SourceText(), trueExpr.SourceText(), falseExpr.SourceText())
			left = jexl.NewTernaryNode(left, trueExpr, falseExpr, source)
			continue
		}

		// Elvis оператор: expr ?: defaultExpr
		if next.typ == tokenQuestionQuestion {
			p.next() // consume '??'
			defaultExpr, err := p.parseExpression(infixPrecedence(tokenQuestionQuestion) + 1)
			if err != nil {
				return nil, err
			}
			source := fmt.Sprintf("%s ?? %s", left.SourceText(), defaultExpr.SourceText())
			left = jexl.NewElvisNode(left, defaultExpr, source)
			continue
		}

		// Бинарные операции
		nextPrec := infixPrecedence(next.typ)
		if nextPrec < 0 || nextPrec < precedence {
			break
		}

		op := p.next()
		right, err := p.parseExpression(nextPrec + 1)
		if err != nil {
			return nil, err
		}
		left = jexl.NewBinaryOpNode(op.literal, left, right, fmt.Sprintf("%s %s %s", left.SourceText(), op.literal, right.SourceText()))
	}

	return left, nil
}

// parseArrayLiteral парсит литерал массива [1, 2, 3]
func (p *simpleParser) parseArrayLiteral() (jexl.Node, error) {
	var elements []jexl.Node

	if p.peek().typ != tokenRBracket {
		for {
			element, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)

			if p.peek().typ == tokenRBracket {
				break
			}
			if err := p.expect(tokenComma); err != nil {
				return nil, err
			}
			if p.peek().typ == tokenRBracket {
				break
			}
		}
	}

	if err := p.expect(tokenRBracket); err != nil {
		return nil, err
	}

	source := "["
	for i, elem := range elements {
		if i > 0 {
			source += ", "
		}
		source += elem.SourceText()
	}
	source += "]"

	return jexl.NewArrayLiteralNode(elements, source), nil
}

// parseMapOrSetLiteral парсит литерал мапы {key: value} или множества {1, 2, 3}
func (p *simpleParser) parseMapOrSetLiteral() (jexl.Node, error) {
	if p.peek().typ == tokenRBrace {
		p.next() // consume '}'
		return jexl.NewSetLiteralNode(nil, "{}"), nil
	}

	// Пробуем определить, это мапа или множество
	// Если следующий токен после первого выражения - двоеточие, это мапа
	firstExpr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	if p.peek().typ == tokenColon {
		// Это мапа: {key: value, ...}
		p.next() // consume ':'
		value, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}

		entries := []jexl.MapEntry{
			{Key: firstExpr, Value: value},
		}

		for p.peek().typ == tokenComma {
			p.next() // consume ','
			if p.peek().typ == tokenRBrace {
				break
			}
			key, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if err := p.expect(tokenColon); err != nil {
				return nil, err
			}
			val, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			entries = append(entries, jexl.MapEntry{Key: key, Value: val})
		}

		if err := p.expect(tokenRBrace); err != nil {
			return nil, err
		}

		source := "{"
		for i, entry := range entries {
			if i > 0 {
				source += ", "
			}
			source += fmt.Sprintf("%s: %s", entry.Key.SourceText(), entry.Value.SourceText())
		}
		source += "}"

		return jexl.NewMapLiteralNode(entries, source), nil
	}

	// Это множество: {1, 2, 3}
	elements := []jexl.Node{firstExpr}

	for p.peek().typ == tokenComma {
		p.next() // consume ','
		if p.peek().typ == tokenRBrace {
			break
		}
		element, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	if err := p.expect(tokenRBrace); err != nil {
		return nil, err
	}

	source := "{"
	for i, elem := range elements {
		if i > 0 {
			source += ", "
		}
		source += elem.SourceText()
	}
	source += "}"

	return jexl.NewSetLiteralNode(elements, source), nil
}

// parseStatement парсит statement (if, for, while, etc.)
func (p *simpleParser) parseStatement() (jexl.Node, error) {
	next := p.peek()
	switch next.typ {
	case tokenIf:
		return p.parseIfStatement()
	case tokenFor:
		return p.parseForStatement()
	case tokenWhile:
		return p.parseWhileStatement()
	case tokenDo:
		return p.parseDoWhileStatement()
	case tokenBreak:
		p.next()
		return jexl.NewBreakNode("break"), nil
	case tokenContinue:
		p.next()
		return jexl.NewContinueNode("continue"), nil
	case tokenReturn:
		return p.parseReturnStatement()
	case tokenLBrace:
		return p.parseBlock()
	default:
		return nil, nil // Не statement, вернём nil
	}
}

// parseIfStatement парсит if/else statement
func (p *simpleParser) parseIfStatement() (jexl.Node, error) {
	p.next() // consume 'if'
	if err := p.expect(tokenLParen); err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	if err := p.expect(tokenRParen); err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatementOrBlock()
	if err != nil {
		return nil, err
	}

	var elseBranch jexl.Node
	if p.match(tokenElse) {
		// Проверяем, является ли это else if
		if p.peek().typ == tokenIf {
			// Это else if - парсим как вложенный if
			elseBranch, err = p.parseIfStatement()
			if err != nil {
				return nil, err
			}
		} else {
			// Обычный else
			elseBranch, err = p.parseStatementOrBlock()
			if err != nil {
				return nil, err
			}
		}
	}

	source := fmt.Sprintf("if (%s) %s", condition.SourceText(), thenBranch.SourceText())
	if elseBranch != nil {
		source += fmt.Sprintf(" else %s", elseBranch.SourceText())
	}

	return jexl.NewIfNode(condition, thenBranch, elseBranch, source), nil
}

// parseForStatement парсит for statement (for (init; condition; step) body или for (var x : items) body)
func (p *simpleParser) parseForStatement() (jexl.Node, error) {
	p.next() // consume 'for'
	if err := p.expect(tokenLParen); err != nil {
		return nil, err
	}

	// Пробуем определить тип цикла: foreach (var x : items) или классический for
	peek := p.peek()
	if peek.typ == tokenVar {
		// Это foreach: for (var x : items)
		p.next() // consume 'var'
		varName := p.next()
		if varName.typ != tokenIdent {
			return nil, p.errorf("expected identifier after 'var'")
		}
		if err := p.expect(tokenColon); err != nil {
			return nil, err
		}
		items, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		if err := p.expect(tokenRParen); err != nil {
			return nil, err
		}
		body, err := p.parseStatementOrBlock()
		if err != nil {
			return nil, err
		}
		source := fmt.Sprintf("for (var %s : %s) %s", varName.literal, items.SourceText(), body.SourceText())
		return jexl.NewForeachNode(jexl.NewIdentifierNode(varName.literal, varName.literal), items, body, source), nil
	} else if peek.typ == tokenIdent {
		// Может быть foreach без var: for (x : items)
		varName := p.next()
		if p.match(tokenColon) {
			// Это foreach: for (x : items)
			items, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if err := p.expect(tokenRParen); err != nil {
				return nil, err
			}
			body, err := p.parseStatementOrBlock()
			if err != nil {
				return nil, err
			}
			source := fmt.Sprintf("for (%s : %s) %s", varName.literal, items.SourceText(), body.SourceText())
			return jexl.NewForeachNode(jexl.NewIdentifierNode(varName.literal, varName.literal), items, body, source), nil
		}
		// Не foreach, возвращаемся назад
		p.pos--
	}

	// Классический for: for (init; condition; step) body
	var init, condition, step jexl.Node
	var err error

	if p.peek().typ != tokenSemicolon {
		init, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}
	if err := p.expect(tokenSemicolon); err != nil {
		return nil, err
	}

	if p.peek().typ != tokenSemicolon {
		condition, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}
	if err := p.expect(tokenSemicolon); err != nil {
		return nil, err
	}

	if p.peek().typ != tokenRParen {
		step, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}
	if err := p.expect(tokenRParen); err != nil {
		return nil, err
	}

	body, err := p.parseStatementOrBlock()
	if err != nil {
		return nil, err
	}

	source := fmt.Sprintf("for (")
	if init != nil {
		source += init.SourceText()
	}
	source += "; "
	if condition != nil {
		source += condition.SourceText()
	}
	source += "; "
	if step != nil {
		source += step.SourceText()
	}
	source += fmt.Sprintf(") %s", body.SourceText())

	return jexl.NewForNode(init, condition, step, body, source), nil
}

// parseWhileStatement парсит while statement
func (p *simpleParser) parseWhileStatement() (jexl.Node, error) {
	p.next() // consume 'while'
	if err := p.expect(tokenLParen); err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	if err := p.expect(tokenRParen); err != nil {
		return nil, err
	}
	body, err := p.parseStatementOrBlock()
	if err != nil {
		return nil, err
	}
	source := fmt.Sprintf("while (%s) %s", condition.SourceText(), body.SourceText())
	return jexl.NewWhileNode(condition, body, source), nil
}

// parseDoWhileStatement парсит do-while statement
func (p *simpleParser) parseDoWhileStatement() (jexl.Node, error) {
	p.next() // consume 'do'
	body, err := p.parseStatementOrBlock()
	if err != nil {
		return nil, err
	}
	if err := p.expect(tokenWhile); err != nil {
		return nil, err
	}
	if err := p.expect(tokenLParen); err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	if err := p.expect(tokenRParen); err != nil {
		return nil, err
	}
	source := fmt.Sprintf("do %s while (%s)", body.SourceText(), condition.SourceText())
	return jexl.NewDoWhileNode(condition, body, source), nil
}

// parseReturnStatement парсит return statement
func (p *simpleParser) parseReturnStatement() (jexl.Node, error) {
	p.next() // consume 'return'
	var value jexl.Node
	var err error
	if p.peek().typ != tokenSemicolon && p.peek().typ != tokenEOF && p.peek().typ != tokenRBrace {
		value, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}
	source := "return"
	if value != nil {
		source += " " + value.SourceText()
	}
	return jexl.NewReturnNode(value, source), nil
}

// parseBlock парсит блок { statements }
func (p *simpleParser) parseBlock() (jexl.Node, error) {
	p.next() // consume '{'
	var statements []jexl.Node

	for p.peek().typ != tokenRBrace && p.peek().typ != tokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			// Пробуем expression
			expr, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			statements = append(statements, expr)
		} else {
			statements = append(statements, stmt)
		}
		if p.match(tokenSemicolon) {
			continue
		}
		if p.peek().typ == tokenRBrace {
			break
		}
	}

	if err := p.expect(tokenRBrace); err != nil {
		return nil, err
	}

	source := "{"
	for i, stmt := range statements {
		if i > 0 {
			source += "; "
		}
		source += stmt.SourceText()
	}
	source += "}"

	return jexl.NewBlockNode(statements, source), nil
}

// parseStatementOrBlock парсит statement или block
func (p *simpleParser) parseStatementOrBlock() (jexl.Node, error) {
	if p.peek().typ == tokenLBrace {
		return p.parseBlock()
	}
	return p.parseStatement()
}

// parseCall парсит вызов метода или функции.
func (p *simpleParser) parseCall(target jexl.Node) (jexl.Node, error) {
	p.next() // consume '('

	var args []jexl.Node
	if p.peek().typ != tokenRParen {
		for {
			arg, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.peek().typ == tokenRParen {
				break
			}
			if err := p.expect(tokenComma); err != nil {
				return nil, err
			}
		}
	}

	if err := p.expect(tokenRParen); err != nil {
		return nil, err
	}

	// Если target - это идентификатор, это вызов функции верхнего уровня
	var method jexl.Node
	if target == nil {
		return nil, p.errorf("invalid call target")
	}

	if ident, ok := target.(*jexl.IdentifierNode); ok {
		// Функция верхнего уровня: func(args)
		method = target
		target = nil
		source := fmt.Sprintf("%s(", ident.Name())
		for i, arg := range args {
			if i > 0 {
				source += ", "
			}
			source += arg.SourceText()
		}
		source += ")"
		return jexl.NewMethodCallNode(nil, method, args, source), nil
	} else if prop, ok := target.(*jexl.PropertyAccessNode); ok {
		// Вызов метода: obj.method(args)
		method = prop.Property()
		target = prop.Object()
		source := fmt.Sprintf("%s.", target.SourceText())
		if ident, ok := method.(*jexl.IdentifierNode); ok {
			source += ident.Name()
		} else {
			source += method.SourceText()
		}
		source += "("
		for i, arg := range args {
			if i > 0 {
				source += ", "
			}
			source += arg.SourceText()
		}
		source += ")"
		return jexl.NewMethodCallNode(target, method, args, source), nil
	} else {
		// Вызов как функции: expr(args) - expr вычисляется и должен быть callable
		method = target
		target = nil
		source := method.SourceText() + "("
		for i, arg := range args {
			if i > 0 {
				source += ", "
			}
			source += arg.SourceText()
		}
		source += ")"
		return jexl.NewMethodCallNode(nil, method, args, source), nil
	}
}

func (p *simpleParser) match(tt tokenType) bool {
	if p.peek().typ == tt {
		p.next()
		return true
	}
	return false
}

func (p *simpleParser) next() token {
	if p.pos >= len(p.tokens) {
		return token{typ: tokenEOF, literal: ""}
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *simpleParser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{typ: tokenEOF, literal: ""}
	}
	return p.tokens[p.pos]
}

func (p *simpleParser) expect(tt tokenType) error {
	if p.peek().typ != tt {
		return p.errorf("expected %v, got %v", tt, p.peek().typ)
	}
	p.next()
	return nil
}

func (p *simpleParser) errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// tokenType представляет тип токена.
type tokenType int

const (
	tokenEOF tokenType = iota
	tokenIdent
	tokenNumber
	tokenString
	tokenBool
	tokenNull
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenPercent
	tokenLParen
	tokenRParen
	tokenLBracket
	tokenRBracket
	tokenLBrace
	tokenRBrace
	tokenDot
	tokenComma
	tokenSemicolon
	tokenColon
	tokenBang
	tokenEqual
	tokenEqualEqual
	tokenBangEqual
	tokenLess
	tokenLessEqual
	tokenGreater
	tokenGreaterEqual
	tokenAnd
	tokenOr
	tokenQuestion
	tokenQuestionQuestion
	// Битовые операторы
	tokenAmpersand      // &
	tokenPipe           // |
	tokenCaret          // ^
	tokenTilde          // ~
	tokenShiftLeft      // <<
	tokenShiftRight     // >>
	tokenShiftRightU    // >>>
	// Строковые операторы
	tokenContains       // =~
	tokenStartsWith     // =^
	tokenEndsWith       // =$
	tokenNotContains    // !~
	tokenNotStartsWith  // !^
	tokenNotEndsWith    // !$
	// Side-effect операторы
	tokenPlusEqual      // +=
	tokenMinusEqual     // -=
	tokenStarEqual      // *=
	tokenSlashEqual     // /=
	tokenPercentEqual   // %=
	tokenPlusPlus       // ++
	tokenMinusMinus     // --
	// Ключевые слова
	tokenIf
	tokenElse
	tokenFor
	tokenWhile
	tokenDo
	tokenBreak
	tokenContinue
	tokenReturn
	tokenVar
	tokenEmpty
	tokenSize
	tokenNot
)

// token представляет токен.
type token struct {
	typ     tokenType
	literal string
	value   any
}

// lexer разбирает исходный текст на токены.
type lexer struct {
	source string
	pos    int
	start  int
}

func newLexer(source string) *lexer {
	return &lexer{source: source}
}

func (l *lexer) lex() []token {
	var tokens []token
	for !l.isAtEnd() {
		l.start = l.pos
		tok := l.nextToken()
		if tok.typ != tokenEOF || len(tokens) == 0 {
			tokens = append(tokens, tok)
		}
		if tok.typ == tokenEOF {
			break
		}
	}
	return tokens
}

func (l *lexer) nextToken() token {
	l.skipWhitespace()
	if l.isAtEnd() {
		return token{typ: tokenEOF}
	}

	l.start = l.pos
	c := l.advance()

	switch c {
	case '+':
		if l.match('=') {
			return token{typ: tokenPlusEqual, literal: "+="}
		}
		if l.match('+') {
			return token{typ: tokenPlusPlus, literal: "++"}
		}
		return token{typ: tokenPlus, literal: "+"}
	case '-':
		if l.match('=') {
			return token{typ: tokenMinusEqual, literal: "-="}
		}
		if l.match('-') {
			return token{typ: tokenMinusMinus, literal: "--"}
		}
		return token{typ: tokenMinus, literal: "-"}
	case '*':
		if l.match('=') {
			return token{typ: tokenStarEqual, literal: "*="}
		}
		return token{typ: tokenStar, literal: "*"}
	case '/':
		if l.match('=') {
			return token{typ: tokenSlashEqual, literal: "/="}
		}
		return token{typ: tokenSlash, literal: "/"}
	case '%':
		if l.match('=') {
			return token{typ: tokenPercentEqual, literal: "%="}
		}
		return token{typ: tokenPercent, literal: "%"}
	case '(':
		return token{typ: tokenLParen, literal: "("}
	case ')':
		return token{typ: tokenRParen, literal: ")"}
	case '[':
		return token{typ: tokenLBracket, literal: "["}
	case ']':
		return token{typ: tokenRBracket, literal: "]"}
	case '{':
		return token{typ: tokenLBrace, literal: "{"}
	case '}':
		return token{typ: tokenRBrace, literal: "}"}
	case '.':
		return token{typ: tokenDot, literal: "."}
	case ',':
		return token{typ: tokenComma, literal: ","}
	case ';':
		return token{typ: tokenSemicolon, literal: ";"}
	case ':':
		return token{typ: tokenColon, literal: ":"}
	case '!':
		if l.match('=') {
			return token{typ: tokenBangEqual, literal: "!="}
		}
		if l.match('~') {
			return token{typ: tokenNotContains, literal: "!~"}
		}
		if l.match('^') {
			return token{typ: tokenNotStartsWith, literal: "!^"}
		}
		if l.match('$') {
			return token{typ: tokenNotEndsWith, literal: "!$"}
		}
		return token{typ: tokenBang, literal: "!"}
	case '=':
		if l.match('=') {
			return token{typ: tokenEqualEqual, literal: "=="}
		}
		if l.match('~') {
			return token{typ: tokenContains, literal: "=~"}
		}
		if l.match('^') {
			return token{typ: tokenStartsWith, literal: "=^"}
		}
		if l.match('$') {
			return token{typ: tokenEndsWith, literal: "=$"}
		}
		return token{typ: tokenEqual, literal: "="}
	case '<':
		if l.match('<') {
			if l.match('=') {
				return token{typ: tokenEOF, literal: "<<="} // TODO: поддержка <<=
			}
			return token{typ: tokenShiftLeft, literal: "<<"}
		}
		if l.match('=') {
			return token{typ: tokenLessEqual, literal: "<="}
		}
		return token{typ: tokenLess, literal: "<"}
	case '>':
		if l.match('>') {
			if l.match('>') {
				return token{typ: tokenShiftRightU, literal: ">>>"}
			}
			if l.match('=') {
				return token{typ: tokenEOF, literal: ">>="} // TODO: поддержка >>=
			}
			return token{typ: tokenShiftRight, literal: ">>"}
		}
		if l.match('=') {
			return token{typ: tokenGreaterEqual, literal: ">="}
		}
		return token{typ: tokenGreater, literal: ">"}
	case '&':
		if l.match('&') {
			return token{typ: tokenAnd, literal: "&&"}
		}
		return token{typ: tokenAmpersand, literal: "&"}
	case '|':
		if l.match('|') {
			return token{typ: tokenOr, literal: "||"}
		}
		return token{typ: tokenPipe, literal: "|"}
	case '^':
		return token{typ: tokenCaret, literal: "^"}
	case '~':
		return token{typ: tokenTilde, literal: "~"}
	case '?':
		if l.match('?') {
			return token{typ: tokenQuestionQuestion, literal: "??"}
		}
		return token{typ: tokenQuestion, literal: "?"}
	case '"', '\'':
		return l.string(c)
	default:
		if unicode.IsDigit(c) {
			return l.number()
		}
		if unicode.IsLetter(c) || c == '_' || c == '$' {
			return l.identifier()
		}
		return l.errorToken("unexpected character")
	}
}

func (l *lexer) string(quote rune) token {
	for l.peek() != quote && !l.isAtEnd() {
		if l.peek() == '\\' {
			l.advance()
		}
		l.advance()
	}

	if l.isAtEnd() {
		return l.errorToken("unterminated string")
	}

	l.advance()                            // закрывающая кавычка
	value := l.source[l.start+1 : l.pos-1] // убираем кавычки
	return token{typ: tokenString, literal: l.source[l.start:l.pos], value: value}
}

func (l *lexer) number() token {
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && unicode.IsDigit(l.peekNext()) {
		l.advance()
		for unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}

	return token{typ: tokenNumber, literal: l.source[l.start:l.pos]}
}

func (l *lexer) identifier() token {
	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' || l.peek() == '$' {
		l.advance()
	}

	text := l.source[l.start:l.pos]
	tokType := tokenIdent
	var value any

	switch text {
	case "true", "false":
		tokType = tokenBool
		value = text == "true"
	case "null", "nil":
		tokType = tokenNull
		value = nil
	case "if":
		tokType = tokenIf
	case "else":
		tokType = tokenElse
	case "for":
		tokType = tokenFor
	case "while":
		tokType = tokenWhile
	case "do":
		tokType = tokenDo
	case "break":
		tokType = tokenBreak
	case "continue":
		tokType = tokenContinue
	case "return":
		tokType = tokenReturn
	case "var":
		tokType = tokenVar
	case "empty":
		tokType = tokenEmpty
	case "size":
		tokType = tokenSize
	case "not":
		tokType = tokenNot
	case "eq":
		tokType = tokenEqualEqual
	case "ne":
		tokType = tokenBangEqual
	case "lt":
		tokType = tokenLess
	case "le":
		tokType = tokenLessEqual
	case "gt":
		tokType = tokenGreater
	case "ge":
		tokType = tokenGreaterEqual
	case "and":
		tokType = tokenAnd
	case "or":
		tokType = tokenOr
	}

	return token{typ: tokType, literal: text, value: value}
}

func (l *lexer) match(expected rune) bool {
	if l.isAtEnd() {
		return false
	}
	if rune(l.source[l.pos]) != expected {
		return false
	}
	l.pos++
	return true
}

func (l *lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	return rune(l.source[l.pos])
}

func (l *lexer) peekNext() rune {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return rune(l.source[l.pos+1])
}

func (l *lexer) advance() rune {
	if l.isAtEnd() {
		return 0
	}
	c := rune(l.source[l.pos])
	l.pos++
	return c
}

func (l *lexer) skipWhitespace() {
	for !l.isAtEnd() {
		c := l.peek()
		if c == ' ' || c == '\r' || c == '\n' || c == '\t' {
			l.advance()
		} else if c == '/' && l.peekNext() == '/' {
			// комментарий
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			break
		}
	}
}

func (l *lexer) isAtEnd() bool {
	return l.pos >= len(l.source)
}

func (l *lexer) errorToken(message string) token {
	return token{typ: tokenEOF, literal: message}
}

// precedence константы для операторов
const (
	prefixPrecedence     = 15
	assignmentPrecedence = 1
	ternaryPrecedence    = 2
)

func infixPrecedence(tt tokenType) int {
	switch tt {
	case tokenEqual, tokenPlusEqual, tokenMinusEqual, tokenStarEqual, tokenSlashEqual, tokenPercentEqual:
		return assignmentPrecedence
	case tokenQuestion:
		return ternaryPrecedence
	case tokenQuestionQuestion:
		return 3
	case tokenDot, tokenLBracket:
		return 14
	case tokenLParen:
		return 13
	case tokenStar, tokenSlash, tokenPercent:
		return 12
	case tokenPlus, tokenMinus:
		return 11
	case tokenShiftLeft, tokenShiftRight, tokenShiftRightU:
		return 10
	case tokenLess, tokenLessEqual, tokenGreater, tokenGreaterEqual:
		return 9
	case tokenEqualEqual, tokenBangEqual:
		return 8
	case tokenContains, tokenStartsWith, tokenEndsWith, tokenNotContains, tokenNotStartsWith, tokenNotEndsWith:
		return 7
	case tokenAmpersand:
		return 5
	case tokenCaret:
		return 4
	case tokenPipe:
		return 3
	case tokenAnd:
		return 2
	case tokenOr:
		return 1
	default:
		return -1
	}
}

func parseNumberLiteral(s string) (any, error) {
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func parseStringLiteral(s string) (string, error) {
	// Убираем кавычки и обрабатываем escape-последовательности
	if len(s) < 2 {
		return "", fmt.Errorf("invalid string literal")
	}
	quote := s[0]
	if s[len(s)-1] != quote {
		return "", fmt.Errorf("invalid string literal")
	}
	value := s[1 : len(s)-1]
	// Простая обработка escape-последовательностей
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\t", "\t")
	value = strings.ReplaceAll(value, "\\r", "\r")
	value = strings.ReplaceAll(value, "\\\"", "\"")
	value = strings.ReplaceAll(value, "\\'", "'")
	value = strings.ReplaceAll(value, "\\\\", "\\")
	return value, nil
}

func isAssignableTarget(node jexl.Node) bool {
	switch node.(type) {
	case *jexl.IdentifierNode, *jexl.PropertyAccessNode, *jexl.IndexAccessNode:
		return true
	default:
		return false
	}
}
