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
		// Пропускаем точку с запятой в начале
		if builder.peek().typ == tokenSemicolon {
			builder.next()
			continue
		}
		// Проверяем, не является ли следующий токен else или else if
		// Если да, и последний добавленный узел - это IfNode, то это часть предыдущего if statement
		if builder.peek().typ == tokenElse {
			// Проверяем, есть ли последний узел и является ли он IfNode
			children := ast.Children()
			if len(children) > 0 {
				if ifNode, ok := children[len(children)-1].(*jexl.IfNode); ok {
					// Это else для предыдущего if - нужно добавить else branch
					builder.next() // consume 'else'
					if builder.peek().typ == tokenIf {
						// Это else if - парсим как вложенный if
						elseBranch, err := builder.parseIfStatement()
						if err != nil {
							return nil, err
						}
						// Обновляем IfNode с else branch
						// Но IfNode неизменяем, нужно создать новый
						// Проще всего - заменить последний узел новым IfNode с else branch
						newIfNode := jexl.NewIfNode(ifNode.Condition(), ifNode.ThenBranch(), elseBranch, 
							fmt.Sprintf("%s else %s", ifNode.SourceText(), elseBranch.SourceText()))
						ast.SetChild(len(children)-1, newIfNode)
						continue
					} else {
						// Обычный else
						elseBranch, err := builder.parseStatementOrBlock()
						if err != nil {
							return nil, err
						}
						newIfNode := jexl.NewIfNode(ifNode.Condition(), ifNode.ThenBranch(), elseBranch,
							fmt.Sprintf("%s else %s", ifNode.SourceText(), elseBranch.SourceText()))
						ast.SetChild(len(children)-1, newIfNode)
						continue
					}
				}
			}
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
		// Пропускаем точку с запятой после statement/expression
		if builder.match(tokenSemicolon) {
			// После точки с запятой продолжаем цикл
			continue
		}
		// Если нет точки с запятой, проверяем, не конец ли это
		if builder.peek().typ == tokenEOF {
			break
		}
		// Если следующий токен - не EOF и не точка с запятой, это может быть допустимо
		// для последнего выражения в скрипте, но лучше попробовать продолжить парсинг
		// Если это не точка с запятой и не EOF, пробуем продолжить парсинг следующего statement/expression
		// Это позволяет обрабатывать скрипты без точек с запятой между statements
		continue
	}

	if err := builder.expect(tokenEOF); err != nil {
		return nil, err
	}

	// Устанавливаем параметры скрипта, если они указаны
	if len(names) > 0 {
		ast.SetParameters(names)
	} else {
		// Если параметры не указаны, проверяем, является ли скрипт только lambda
		// Если да, устанавливаем параметры lambda в ScriptNode
		children := ast.Children()
		if len(children) == 1 {
			if lambdaNode, ok := children[0].(*jexl.LambdaNode); ok {
				// Скрипт содержит только lambda - устанавливаем параметры lambda
				params := make([]string, 0, len(lambdaNode.Parameters()))
				for _, param := range lambdaNode.Parameters() {
					params = append(params, param.Name())
				}
				ast.SetParameters(params)
			}
		}
	}

	return ast, nil
}

// simpleParser реализует примитивный Pratt-парсер для базовой арифметики.
type simpleParser struct {
	info      *jexl.Info
	source    string
	features  *jexl.Features
	tokens    []token
	pos       int
	loopCount int // Счетчик вложенных циклов для проверки break/continue
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
		// Используем prefixPrecedence + 1 для правоассоциативности унарных операторов
		operand, err := p.parseExpression(prefixPrecedence + 1)
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
	case tokenSize:
		// Функция size: size(x) или size x
		// Проверяем, есть ли скобки после size
		if p.peek().typ == tokenLParen {
			// size(x) - вызов функции
			sizeIdent := jexl.NewIdentifierNode("size", "size")
			left = sizeIdent
			// parseCall обработает вызов
		} else {
			// size x - оператор без скобок
			operand, err := p.parseExpression(prefixPrecedence)
			if err != nil {
				return nil, err
			}
			// Создаём вызов функции size(operand)
			sizeIdent := jexl.NewIdentifierNode("size", "size")
			args := []jexl.Node{operand}
			left = jexl.NewMethodCallNode(nil, sizeIdent, args, fmt.Sprintf("size(%s)", operand.SourceText()))
		}
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
		// Проверяем, не является ли это lambda функцией с одним параметром без скобок: x -> x + 1
		if p.isLambdaStartAfterIdent() {
			// Это lambda функция - параметр уже прочитан в tok
			param := jexl.NewIdentifierNode(tok.literal, tok.literal)
			parameters := []*jexl.IdentifierNode{param}
			
			// Парсим стрелку: -> или =>
			var arrow string
			if p.peek().typ == tokenLambda {
				p.next() // consume '->'
				arrow = "->"
			} else if p.peek().typ == tokenFatArrow {
				p.next() // consume '=>'
				arrow = "=>"
			} else {
				return nil, p.errorf("expected '->' or '=>' in lambda")
			}
			
			// Парсим тело lambda
			var body jexl.Node
			var err error
			if p.peek().typ == tokenLBrace {
				body, err = p.parseBlock()
			} else {
				body, err = p.parseExpression(0)
			}
			if err != nil {
				return nil, err
			}
			
			source := tok.literal + " " + arrow + " " + body.SourceText()
			left = jexl.NewLambdaNode(parameters, body, source)
		} else {
			left = jexl.NewIdentifierNode(tok.literal, tok.literal)
		}
	case tokenBool:
		left = jexl.NewLiteralNode(tok.value, tok.literal)
	case tokenNull:
		left = jexl.NewLiteralNode(nil, tok.literal)
	case tokenLBracket:
		// Массив: [1, 2, 3]
		arrayNode, err := p.parseArrayLiteral()
		if err != nil {
			return nil, err
		}
		left = arrayNode
	case tokenLBrace:
		// Мапа или множество: {key: value} или {1, 2, 3}
		return p.parseMapOrSetLiteral()
	case tokenLParen:
		// Проверяем, не является ли это lambda функцией
		// Lookahead: если после ( идут идентификаторы, а затем -> или =>, то это lambda
		// Токен ( уже прочитан, поэтому проверяем следующий токен
		savedPos := p.pos
		// isLambdaStartAfterLParen проверяет lambda после уже прочитанной (
		isLambda := p.isLambdaStartAfterLParen()
		// Восстанавливаем позицию
		p.pos = savedPos
		if isLambda {
			// Это lambda функция - ( уже прочитан в tok, парсим параметры
			// parseLambda должен знать, что ( уже прочитан
			lambda, err := p.parseLambda(true) // true = ( уже прочитан
			if err != nil {
				return nil, err
			}
			left = lambda
		} else {
			// Обычное выражение в скобках
			expr, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if err := p.expect(tokenRParen); err != nil {
				return nil, err
			}
			left = expr
		}
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

		// Доступ к свойству: expr.prop или expr.1 (числовое свойство)
		if next.typ == tokenDot {
			p.next() // consume '.'
			prop := p.peek()
			var propName string
			if prop.typ == tokenIdent {
				prop = p.next()
				propName = prop.literal
			} else if prop.typ == tokenNumber {
				prop = p.next()
				propName = prop.literal
			} else {
				return nil, p.errorf("expected identifier or number after '.'")
			}
			propNode := jexl.NewIdentifierNode(propName, propName)
			left = jexl.NewPropertyAccessNode(left, propNode, fmt.Sprintf("%s.%s", left.SourceText(), propName))
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
		// Или Elvis оператор: expr ?: defaultExpr
		if next.typ == tokenQuestion && precedence < ternaryPrecedence {
			p.next() // consume '?'
			// Проверяем, не является ли это Elvis оператором ?:
			if p.peek().typ == tokenColon {
				// Это Elvis оператор ?: (expr ?: defaultExpr)
				p.next() // consume ':'
				defaultExpr, err := p.parseExpression(ternaryPrecedence)
				if err != nil {
					return nil, err
				}
				source := fmt.Sprintf("%s ?: %s", left.SourceText(), defaultExpr.SourceText())
				left = jexl.NewElvisNode(left, defaultExpr, source)
				continue
			}
			// Это тернарный оператор ? : (condition ? trueExpr : falseExpr)
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

		// Elvis оператор: expr ?? defaultExpr
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

		// Range оператор: left .. right
		// Проверяем приоритет перед обработкой
		if next.typ == tokenRange {
			nextPrec := infixPrecedence(tokenRange)
			if nextPrec < precedence {
				break
			}
			p.next() // consume '..'
			right, err := p.parseExpression(infixPrecedence(tokenRange) + 1)
			if err != nil {
				return nil, err
			}
			source := fmt.Sprintf("%s .. %s", left.SourceText(), right.SourceText())
			left = jexl.NewRangeNode(left, right, source)
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
	// Пустая мапа {} должна быть мапой, а не множеством
	if p.peek().typ == tokenRBrace {
		p.next() // consume '}'
		return jexl.NewMapLiteralNode(nil, "{}"), nil
	}

	// Проверяем, не является ли это пустой мапой {:}
	if p.peek().typ == tokenColon {
		p.next() // consume ':'
		if p.peek().typ == tokenRBrace {
			p.next() // consume '}'
			return jexl.NewMapLiteralNode(nil, "{:}"), nil
		}
		// Если после : идет что-то еще, это ошибка
		return nil, p.errorf("unexpected token after ':' in map literal")
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

// parseLambda парсит lambda функцию: (x, y) -> x + y или x -> x + 1 или (x, y) => x + y
// lparenRead указывает, был ли уже прочитан токен ( (true) или нет (false)
func (p *simpleParser) parseLambda(lparenRead bool) (jexl.Node, error) {
	// Проверяем, включена ли поддержка lambda
	if p.features != nil && !p.features.SupportsLambda() {
		return nil, p.errorf("lambda functions are not enabled")
	}

	var parameters []*jexl.IdentifierNode
	var sourceStart string

	// Парсим параметры
	// Проверяем, является ли текущий токен ( или идентификатор
	if !lparenRead && p.peek().typ == tokenLParen {
		// Множественные параметры в скобках: (x, y) - ( еще не прочитан
		p.next() // consume '('
		sourceStart = "("
	} else if p.peek().typ == tokenIdent {
		// Один параметр без скобок: x -> ...
		// Или параметры в скобках: (x) -> ... или (x, y) -> ...
		if lparenRead {
			// ( уже прочитан, начинаем с идентификатора
			sourceStart = "("
		} else {
			// Проверяем, есть ли запятая после идентификатора - если да, то это множественные параметры
			savedPos := p.pos
			p.pos++
			hasComma := p.pos < len(p.tokens) && p.tokens[p.pos].typ == tokenComma
			p.pos = savedPos
			
			if hasComma {
				// Множественные параметры в скобках: (x, y) - но ( еще не прочитан, это ошибка
				return nil, p.errorf("expected '(' before lambda parameters")
			} else {
				// Один параметр без скобок: x -> ...
				paramTok := p.next()
				param := jexl.NewIdentifierNode(paramTok.literal, paramTok.literal)
				parameters = append(parameters, param)
				sourceStart = paramTok.literal
				
				// Парсим стрелку и тело
				var arrow string
				if p.peek().typ == tokenLambda {
					p.next() // consume '->'
					arrow = "->"
				} else if p.peek().typ == tokenFatArrow {
					p.next() // consume '=>'
					arrow = "=>"
				} else {
					return nil, p.errorf("expected '->' or '=>' in lambda")
				}
				
				var body jexl.Node
				var err error
				if p.peek().typ == tokenLBrace {
					body, err = p.parseBlock()
				} else {
					body, err = p.parseExpression(0)
				}
				if err != nil {
					return nil, err
				}
				
				source := sourceStart + " " + arrow + " " + body.SourceText()
				return jexl.NewLambdaNode(parameters, body, source), nil
			}
		}
	} else {
		return nil, p.errorf("expected lambda parameters")
	}

	// Парсим множественные параметры в скобках
	for {
		if p.peek().typ == tokenRParen {
			p.next() // consume ')'
			sourceStart += ")"
			break
		}
		if p.peek().typ != tokenIdent {
			return nil, p.errorf("expected identifier in lambda parameters")
		}
		paramTok := p.next()
		param := jexl.NewIdentifierNode(paramTok.literal, paramTok.literal)
		parameters = append(parameters, param)
		if len(parameters) > 1 {
			sourceStart += ", "
		}
		sourceStart += paramTok.literal

		if p.peek().typ == tokenComma {
			p.next() // consume ','
			sourceStart += ","
			continue
		}
		if p.peek().typ == tokenRParen {
			p.next() // consume ')'
			sourceStart += ")"
			break
		}
	}

	// Парсим стрелку: -> или =>
	var arrow string
	if p.peek().typ == tokenLambda {
		p.next() // consume '->'
		arrow = "->"
	} else if p.peek().typ == tokenFatArrow {
		p.next() // consume '=>'
		arrow = "=>"
	} else {
		return nil, p.errorf("expected '->' or '=>' in lambda")
	}

	// Парсим тело lambda (выражение или блок)
	var body jexl.Node
	var err error
	if p.peek().typ == tokenLBrace {
		// Блок: { return x + y; }
		body, err = p.parseBlock()
		if err != nil {
			return nil, err
		}
	} else {
		// Выражение: x + y
		body, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	source := sourceStart + " " + arrow + " " + body.SourceText()
	return jexl.NewLambdaNode(parameters, body, source), nil
}

// parseStatement парсит statement (if, for, while, etc.)
func (p *simpleParser) parseStatement() (jexl.Node, error) {
	next := p.peek()
	switch next.typ {
	case tokenIf:
		return p.parseIfStatement()
	case tokenFor:
		// Проверяем, включены ли циклы
		if p.features != nil && !p.features.SupportsLoops() {
			return nil, p.errorf("loops are not enabled")
		}
		return p.parseForStatement()
	case tokenWhile:
		return p.parseWhileStatement()
	case tokenDo:
		return p.parseDoWhileStatement()
	case tokenBreak:
		// Проверяем, включены ли циклы (break используется в циклах)
		if p.features != nil && !p.features.SupportsLoops() {
			return nil, p.errorf("loops are not enabled")
		}
		// Проверяем, что break используется внутри цикла
		if p.loopCount == 0 {
			return nil, p.errorf("break statement not within a loop")
		}
		p.next()
		return jexl.NewBreakNode("break"), nil
	case tokenContinue:
		// Проверяем, включены ли циклы (continue используется в циклах)
		if p.features != nil && !p.features.SupportsLoops() {
			return nil, p.errorf("loops are not enabled")
		}
		// Проверяем, что continue используется внутри цикла
		if p.loopCount == 0 {
			return nil, p.errorf("continue statement not within a loop")
		}
		p.next()
		return jexl.NewContinueNode("continue"), nil
	case tokenReturn:
		return p.parseReturnStatement()
	case tokenVar:
		return p.parseVarStatement()
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
	// Проверяем, является ли следующий токен else
	// Это может быть сразу после thenBranch или после точки с запятой (если thenBranch - это statement с точкой с запятой)
	// Но в JEXL else может идти после if statement без точки с запятой между ними
	// Пропускаем точку с запятой, если она есть перед else
	if p.peek().typ == tokenSemicolon {
		// Проверяем, не является ли следующий токен после точки с запятой else
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].typ == tokenElse {
			p.next() // consume ';'
		}
	}
	if p.peek().typ == tokenElse {
		p.next() // consume 'else'
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
	
	// Увеличиваем счетчик циклов
	p.loopCount++
	defer func() { p.loopCount-- }()

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
		if items == nil {
			return nil, p.errorf("expected expression after ':'")
		}
		// body может быть nil для пустого statement (точка с запятой)
		source := fmt.Sprintf("for (var %s : %s)", varName.literal, items.SourceText())
		if body != nil {
			source += " " + body.SourceText()
		} else {
			source += " ;"
		}
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
			if items == nil {
				return nil, p.errorf("expected expression after ':'")
			}
			if err := p.expect(tokenRParen); err != nil {
				return nil, err
			}
			body, err := p.parseStatementOrBlock()
			if err != nil {
				return nil, err
			}
			// body может быть nil для пустого statement (точка с запятой)
			source := fmt.Sprintf("for (%s : %s)", varName.literal, items.SourceText())
			if body != nil {
				source += " " + body.SourceText()
			} else {
				source += " ;"
			}
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
	// Проверяем, включены ли циклы
	if p.features != nil && !p.features.SupportsLoops() {
		return nil, p.errorf("loops are not enabled")
	}
	p.next() // consume 'while'
	if err := p.expect(tokenLParen); err != nil {
		return nil, err
	}
	
	// Увеличиваем счетчик циклов
	p.loopCount++
	defer func() { p.loopCount-- }()
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
	source := fmt.Sprintf("while (%s)", condition.SourceText())
	if body != nil {
		source += " " + body.SourceText()
	} else {
		source += " ;"
	}
	return jexl.NewWhileNode(condition, body, source), nil
}

// parseDoWhileStatement парсит do-while statement
func (p *simpleParser) parseDoWhileStatement() (jexl.Node, error) {
	// Проверяем, включены ли циклы
	if p.features != nil && !p.features.SupportsLoops() {
		return nil, p.errorf("loops are not enabled")
	}
	p.next() // consume 'do'
	
	// Увеличиваем счетчик циклов
	p.loopCount++
	defer func() { p.loopCount-- }()
	body, err := p.parseStatementOrBlock()
	if err != nil {
		return nil, err
	}
	// После body может быть точка с запятой, пропускаем её
	if p.peek().typ == tokenSemicolon {
		p.next() // consume ';'
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
	source := "do"
	if body != nil {
		source += " " + body.SourceText()
	} else {
		source += " ;"
	}
	source += fmt.Sprintf(" while (%s)", condition.SourceText())
	return jexl.NewDoWhileNode(condition, body, source), nil
}

// parseReturnStatement парсит return statement
func (p *simpleParser) parseReturnStatement() (jexl.Node, error) {
	p.next() // consume 'return'
	var value jexl.Node
	var err error
	// Проверяем, не является ли следующий токен концом statement или else
	// else может идти после return в контексте if statement
	next := p.peek()
	if next.typ != tokenSemicolon && next.typ != tokenEOF && next.typ != tokenRBrace && next.typ != tokenElse {
		// Парсим выражение, но нужно быть осторожным - если следующее выражение заканчивается на else,
		// то это не часть return, а начало else clause
		// Парсим выражение с низким приоритетом, чтобы остановиться на else
		value, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		// После парсинга выражения проверяем, не является ли следующий токен else
		// Если да, то это не часть return, а начало else clause - мы уже распарсили return statement
		if p.peek().typ == tokenElse {
			// else идет после return - это нормально, return statement закончен
		}
	}
	source := "return"
	if value != nil {
		source += " " + value.SourceText()
	}
	return jexl.NewReturnNode(value, source), nil
}

// parseVarStatement парсит var statement (var x или var x = value)
func (p *simpleParser) parseVarStatement() (jexl.Node, error) {
	p.next() // consume 'var'
	nameTok := p.next()
	if nameTok.typ != tokenIdent {
		return nil, p.errorf("expected identifier after 'var'")
	}
	name := jexl.NewIdentifierNode(nameTok.literal, nameTok.literal)
	
	var value jexl.Node
	var err error
	source := "var " + nameTok.literal
	
	// Проверяем, есть ли присваивание
	if p.peek().typ == tokenEqual {
		p.next() // consume '='
		value, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		source += " = " + value.SourceText()
	}
	
	return jexl.NewVarNode(name, value, source), nil
}

// parseBlock парсит блок { statements }
func (p *simpleParser) parseBlock() (jexl.Node, error) {
	p.next() // consume '{'
	
	// Проверяем, не является ли это пустым блоком
	if p.peek().typ == tokenRBrace {
		p.next() // consume '}'
		return jexl.NewBlockNode(nil, "{}"), nil
	}
	
	// НЕ пытаемся парсить блок как литерал множества/мапы
	// Блоки всегда должны парситься как BlockNode, даже если они содержат только одно выражение
	// Это важно для lambda функций: (x) -> { x + 1 } должно быть блоком, а не множеством
	
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
	// Проверяем, является ли это пустым statement (;)
	if p.peek().typ == tokenSemicolon {
		p.next() // consume ';'
		// Возвращаем nil для пустого statement
		return nil, nil
	}
	node, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if node != nil {
		// После statement может быть точка с запятой, но мы её не обрабатываем здесь
		// Она будет обработана в ParseScript
		return node, nil
	}
	// Если не statement, пробуем парсить как выражение
	expr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	// После выражения может быть точка с запятой, но мы её не обрабатываем здесь
	// Она будет обработана в ParseScript
	return expr, nil
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
	tokenRange          // ..
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
	// Lambda операторы
	tokenLambda   // ->
	tokenFatArrow // =>
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
		if l.match('>') {
			return token{typ: tokenLambda, literal: "->"}
		}
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
		if l.match('.') {
			return token{typ: tokenRange, literal: ".."}
		}
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
		if l.match('>') {
			return token{typ: tokenFatArrow, literal: "=>"}
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
	case tokenShiftLeft, tokenShiftRight, tokenShiftRightU:
		return 11
	case tokenPlus, tokenMinus:
		return 10
	case tokenLess, tokenLessEqual, tokenGreater, tokenGreaterEqual:
		return 9
	case tokenEqualEqual, tokenBangEqual:
		return 8
	case tokenRange:
		return 7 // Range имеет тот же приоритет, что и строковые операторы
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

// isLambdaStart проверяет, является ли текущая позиция началом lambda функции в скобках: (x, y) -> ...
// Предполагает, что ( уже прочитан
func (p *simpleParser) isLambdaStartAfterLParen() bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	// Сохраняем текущую позицию
	savedPos := p.pos
	defer func() {
		p.pos = savedPos
	}()

	// Проверяем, идут ли после ( идентификаторы
	hasParams := false
	for p.pos < len(p.tokens) {
		tok := p.tokens[p.pos]
		if tok.typ == tokenIdent {
			hasParams = true
			p.pos++
			// Проверяем, есть ли запятая или закрывающая скобка
			if p.pos >= len(p.tokens) {
				return false
			}
			next := p.tokens[p.pos]
			if next.typ == tokenComma {
				p.pos++
				continue
			}
			if next.typ == tokenRParen {
				p.pos++
				break
			}
			// Если после идентификатора идет что-то другое (не запятая и не скобка), это не lambda
			return false
		} else if tok.typ == tokenRParen {
			// Закрывающая скобка - выходим из цикла
			if hasParams {
				p.pos++
				break
			}
			// Если нет параметров, это не lambda
			return false
		} else {
			// Любой другой токен - это не lambda
			return false
		}
	}

	// После ) должна идти стрелка -> или =>
	if p.pos < len(p.tokens) {
		next := p.tokens[p.pos]
		return next.typ == tokenLambda || next.typ == tokenFatArrow
	}
	return false
}

// isLambdaStart проверяет, является ли текущая позиция началом lambda функции в скобках: (x, y) -> ...
func (p *simpleParser) isLambdaStart() bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	// Сохраняем текущую позицию
	savedPos := p.pos
	defer func() {
		p.pos = savedPos
	}()

	// Пропускаем (
	if p.pos >= len(p.tokens) || p.tokens[p.pos].typ != tokenLParen {
		return false
	}
	p.pos++

	// Используем isLambdaStartAfterLParen для проверки после (
	return p.isLambdaStartAfterLParen()
}

// isLambdaStartAfterIdent проверяет, является ли текущая позиция началом lambda функции после идентификатора: x -> ...
func (p *simpleParser) isLambdaStartAfterIdent() bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	// Следующий токен должен быть -> или =>
	next := p.tokens[p.pos]
	return next.typ == tokenLambda || next.typ == tokenFatArrow
}
