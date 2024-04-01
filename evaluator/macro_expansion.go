package evaluator

import (
	"banana/ast"
	"banana/object"
)

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("We only support returning AST ndoes from macros.")
		}
		return quote.Node
	})
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	id, ok := exp.Fun.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(id.Val)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, a := range exp.Args {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)
	for paramIndex, param := range macro.Parameters {
		extended.Set(param.Val, args[paramIndex])
	}
	return extended
}

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}
	for i,stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i)
		}
	}
	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex + 1:]...,)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStmt.Val.(*ast.MacroLiteral)
	if !ok {
		return false
	}
	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLit, _ := letStmt.Val.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLit.Parameters,
		Env: env,
		Body: macroLit.Body,
	}
	env.Set(letStmt.Name.Val, macro)
}