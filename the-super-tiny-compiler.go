package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Token struct {
	Type  string
	Value string
}

func Tokenizer(input string) []Token {
	input += "\n"

	current := 0
	var tokens []Token

	for current < len([]rune(input)) {
		c := string([]rune(input)[current])

		if c == "(" {
			tokens = append(tokens, Token{
				Type:  "paren",
				Value: "(",
			})

			current++
			continue
		}

		if c == ")" {
			tokens = append(tokens, Token{
				Type:  "paren",
				Value: ")",
			})

			current++
			continue
		}

		if c == " " {
			current++
			continue
		}

		if isNumber(c) {
			fullNumber := ""
			for isNumber(c) {
				fullNumber += c
				current++
				c = string([]rune(input)[current])
			}

			tokens = append(tokens, Token{
				Type:  "number",
				Value: fullNumber,
			})
			continue
		}

		if isLetter(c) {
			fullWord := ""
			for isLetter(c) {
				fullWord += c
				current++
				c = string([]rune(input)[current])
			}

			tokens = append(tokens, Token{
				Type:  "name",
				Value: fullWord,
			})

			continue
		}
		break
	}

	return tokens
}

func isLetter(c string) bool {
	if c == "" {
		return false
	}
	n := []rune(c)[0]
	if n >= 'a' && n <= 'z' {
		return true
	}
	return false
}

func isNumber(c string) bool {
	if _, err := strconv.Atoi(c); err != nil {
		return false
	}
	return true
}

type Node struct {
	Type       string
	Value      string
	Body       []Node
	Name       string
	Params     []Node
	Context    *[]Node
	Callee     *Node
	Arguments  *[]Node
	Expression *Node
}

var counter int
var globalTokens []Token

func Parser(tokens []Token) Node {
	counter = 0
	globalTokens = tokens

	ast := Node{
		Type: "Program",
		Body: []Node{},
	}

	for counter < len(tokens) {
		ast.Body = append(ast.Body, Walk())
	}

	return ast
}

func Walk() Node {
	token := globalTokens[counter]

	if token.Type == "number" {
		counter++

		return Node{
			Type:  "NumberLiteral",
			Value: token.Value,
		}
	}

	if token.Type == "paren" && token.Value == "(" {
		counter++
		token = globalTokens[counter]

		node := Node{
			Type:   "CallExpression",
			Name:   token.Value,
			Params: []Node{},
		}

		counter++

		token = globalTokens[counter]

		for token.Type != "paren" || (token.Type == "paren" && token.Value != ")") {
			node.Params = append(node.Params, Walk())
			token = globalTokens[counter]
		}

		counter++

		return node

	}

	log.Fatalf("Unknown token type %s %s", token.Type, token.Value)
	return Node{}
}

type Visitor map[string]func(node *Node, parent Node)

func Traverser(ast Node, visitor Visitor) {
	TraverserNode(ast, Node{}, visitor)
}

func TraverserArray(array []Node, visitor Visitor, parent Node) {
	for _, node := range array {
		TraverserNode(node, parent, visitor)
	}
}

func TraverserNode(node, parent Node, visitor Visitor) {
	for k, fn := range visitor {
		if k == node.Type {
			fn(&node, parent)
		}
	}

	switch node.Type {
	case "Program":
		TraverserArray(node.Body, visitor, node)
	case "CallExpression":
		TraverserArray(node.Params, visitor, node)
	case "NumberLiteral":
		break
	default:
		panic("Unknown node type")
	}
}

func Transformer(ast Node) Node {

	newAst := Node{
		Type: "Program",
		Body: []Node{},
	}

	ast.Context = &newAst.Body

	Traverser(ast, Visitor{
		"NumberLiteral": func(node *Node, parent Node) {
			*parent.Context = append(*parent.Context, Node{
				Type:  "NumberLiteral",
				Value: node.Value,
			})
		},
		"CallExpression": func(node *Node, parent Node) {
			expr := Node{
				Type: "CallExpression",
				Callee: &Node{
					Type: "Identifier",
					Name: node.Name,
				},
				Arguments: new([]Node),
			}

			node.Context = expr.Arguments

			if parent.Type != "CallExpression" {
				exprs := Node{
					Type:       "ExpressionStatement",
					Expression: &expr,
				}
				*parent.Context = append(*parent.Context, exprs)
			} else {
				*parent.Context = append(*parent.Context, expr)
			}

		},
	})

	return newAst
}

func CodeGenerator(node Node) string {
	switch node.Type {
	case "Program":
		var out []string
		for _, n := range node.Body {
			out = append(out, CodeGenerator(n))
		}

		return strings.Join(out, "\n")
	case "ExpressionStatement":
		return CodeGenerator(*node.Expression) + ";"
	case "CallExpression":
		var args []string
		codeGen := CodeGenerator(*node.Callee)
		for _, arg := range *node.Arguments {
			args = append(args, CodeGenerator(arg))
		}

		return codeGen + "(" + strings.Join(args, ", ") + ")"
	case "Identifier":
		return node.Name
	case "NumberLiteral":
		return node.Value
	default:
		panic("Unknown node type")
		return ""
	}
}

func Compiler(input string) string {
	tokens := Tokenizer(input)
	ast := Parser(tokens)
	newAst := Transformer(ast)
	output := CodeGenerator(newAst)
	return output
}

func main() {
	program := "(add 10 (subtract 10 6))"
	out := Compiler(program)
	fmt.Println(out)
}
