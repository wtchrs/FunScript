package parser

import (
	"fmt"
	"funscript/ast"
	"funscript/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not 'return', got %q", returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	expected := "hello world"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != expected {
		t.Errorf("literal.Value not %q. got=%q", expected, literal.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	booleanTests := []struct {
		input    string
		expected bool
	}{
		{
			"true;",
			true,
		},
		{
			"false;",
			false,
		},
	}

	for _, tt := range booleanTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testBoolean(t, stmt.Expression, tt.expected)
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(literal.Elements) != 3 {
		t.Fatalf("literal.Elements does not contain 3 elements. got=%d", len(literal.Elements))
	}

	testIntegerLiteral(t, literal.Elements[0], 1)
	testInfixExpression(t, literal.Elements[1], "*", 2, 2)
	testInfixExpression(t, literal.Elements[2], "+", 3, 3)
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Pairs is not empty, but has %d items.", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs does not contain 3 elements. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiterals(t *testing.T) {
	input := `
{
	"one": 1,
	"two": 2,
	"three": 3,
	1: "one",
	2: "two",
	3: "three",
	true: "TRUE",
	false: "FALSE"
}
`
	expected := map[interface{}]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
		1:       "one",
		2:       "two",
		3:       "three",
		true:    "TRUE",
		false:   "FALSE",
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != len(expected) {
		t.Fatalf("hash.Pairs does not contain %d elements. got=%d", len(expected), len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		switch key := key.(type) {
		case *ast.StringLiteral:
			expectedValue := expected[key.Value]
			testValueLiteralExpression(t, value, expectedValue)
		case *ast.IntegerLiteral:
			expectedValue := expected[int(key.Value)]
			testValueLiteralExpression(t, value, expectedValue)
		case *ast.Boolean:
			expectedValue := expected[key.Value]
			testValueLiteralExpression(t, value, expectedValue)
		default:
			t.Errorf("wrong key type. got=%T (%+v)", key, key)
		}
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs does not contain 3 elements. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(expression ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, "+", 0, 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, "-", 10, 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, "/", 15, 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		testFunc := tests[literal.Value]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
		}

		testFunc(value)
	}
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not  *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, exp.Index, "+", 1, 1) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{
			"!5;",
			"!",
			5,
		},
		{
			"-15",
			"-",
			15,
		},
		{
			"!true;",
			"!",
			true,
		},
		{
			"!false;",
			"!",
			false,
		},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{
			"5 + 5;",
			5,
			"+",
			5,
		},
		{
			"5 - 5;",
			5,
			"-",
			5,
		},
		{
			"5 * 5;",
			5,
			"*",
			5,
		},
		{
			"5 / 5;",
			5,
			"/",
			5,
		},
		{
			"5 > 5;",
			5,
			">",
			5,
		},
		{
			"5 < 5;",
			5,
			"<",
			5,
		},
		{
			"5 >= 5;",
			5,
			">=",
			5,
		},
		{
			"5 <= 5;",
			5,
			"<=",
			5,
		},
		{
			"5 == 5;",
			5,
			"==",
			5,
		},
		{
			"5 != 5;",
			5,
			"!=",
			5,
		},
		{
			"true == true",
			true,
			"==",
			true,
		},
		{
			"true != false",
			true,
			"!=",
			false,
		},
		{
			"false == false",
			false,
			"==",
			false,
		},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.operator, tt.leftValue, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"true == 3 < 5",
			"(true == (3 < 5))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / ( 5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "<", "x", "y") {
		return
	}

	conseq, ok := exp.Consequence.(*ast.BlockStatement)
	if !ok {
		t.Errorf("consequence is not ast.BlockStatement. got=%T", exp.Consequence)
	}

	if len(conseq.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(conseq.Statements))
	}

	consequence, ok := conseq.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", conseq.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestParsingIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "<", "x", "y") {
		return
	}

	conseq := exp.Consequence.(*ast.BlockStatement)
	if !ok {
		t.Errorf("consequence is not ast.BlockStatement. got=%T", exp.Consequence)
	}

	if len(conseq.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(conseq.Statements))
	}

	consequence, ok := conseq.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", conseq.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil.")
	}

	alter := exp.Alternative.(*ast.BlockStatement)
	if !ok {
		t.Errorf("alternative is not ast.BlockStatement. got=%T", exp.Alternative)
	}

	if len(alter.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d", len(alter.Statements))
	}

	alternative, ok := alter.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", alter.Statements[0])
	}

	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Errorf("Function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Errorf("function.Body.Statements has no 1 statements. got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("function.Body.Statements[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "+", "x", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			"fn() {};",
			[]string{},
		},
		{
			"fn(x) {};",
			[]string{"x"},
		},
		{
			"fn(x, y, z) {};",
			[]string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. expected=%d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Errorf("length parameters wrong. expected=%d, got=%d\n", 3, len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], "*", 2, 3)
	testInfixExpression(t, exp.Arguments[2], "+", 4, 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input             string
		expectedIdent     string
		expectedArguments []string
	}{
		{
			"function();",
			"function",
			[]string{},
		},
		{
			"print(1)",
			"print",
			[]string{"1"},
		},
		{
			"sum(1, 2 * 3, 4 + 5)",
			"sum",
			[]string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArguments) {
			t.Errorf("length parameters wrong. expected=%d, got=%d\n", len(tt.expectedArguments), len(exp.Arguments))
		}

		for i, expected := range tt.expectedArguments {
			argument := exp.Arguments[i].String()
			if argument != expected {
				t.Errorf("argument %d wrong. expected=%s, got=%s", i, expected, argument)
			}
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testInfixExpression(t *testing.T, exp ast.Expression, operator string, left interface{}, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T", exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func testValueLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
	}

	if b.Value != value {
		t.Errorf("b.Value not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() not %t. got=%s", value, b.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
	}

	if str.Value != value {
		t.Errorf("str.Value not %q. got=%q", value, str.Value)
		return false
	}

	if str.TokenLiteral() != value {
		t.Errorf("str.TokenLiteral() not %q. got=%q", value, str.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value interface{}) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	if !testLiteralExpression(t, letStmt.Value, value) {
		return false
	}

	return true
}
