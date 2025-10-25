package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/senither/zen-lang/tokens"
)

type Node interface {
	GetToken() tokens.Token
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) GetToken() tokens.Token { return tokens.Token{} }
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type EmptyStatement struct {
	Token tokens.Token
}

func (es *EmptyStatement) statementNode()         {}
func (es *EmptyStatement) GetToken() tokens.Token { return es.Token }
func (es *EmptyStatement) TokenLiteral() string   { return es.Token.Literal }
func (es *EmptyStatement) String() string         { return "" }

type VariableStatement struct {
	Token   tokens.Token
	Name    *Identifier
	Value   Expression
	Mutable bool
}

func (vs *VariableStatement) statementNode()         {}
func (vs *VariableStatement) GetToken() tokens.Token { return vs.Token }
func (vs *VariableStatement) TokenLiteral() string   { return vs.Token.Literal }
func (ls *VariableStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")

	if ls.Mutable {
		out.WriteString("mut ")
	}

	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode()        {}
func (i *Identifier) GetToken() tokens.Token { return i.Token }
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string { return i.Value }

type ReturnStatement struct {
	Token       tokens.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()         {}
func (rs *ReturnStatement) GetToken() tokens.Token { return rs.Token }
func (rs *ReturnStatement) TokenLiteral() string   { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()         {}
func (es *ExpressionStatement) GetToken() tokens.Token { return es.Token }
func (es *ExpressionStatement) TokenLiteral() string   { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()        {}
func (il *IntegerLiteral) GetToken() tokens.Token { return il.Token }
func (il *IntegerLiteral) TokenLiteral() string   { return il.Token.Literal }
func (il *IntegerLiteral) String() string         { return il.Token.Literal }

type FloatLiteral struct {
	Token tokens.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()        {}
func (fl *FloatLiteral) GetToken() tokens.Token { return fl.Token }
func (fl *FloatLiteral) TokenLiteral() string   { return fl.Token.Literal }
func (fl *FloatLiteral) String() string         { return fl.Token.Literal }

type StringLiteral struct {
	Token tokens.Token
	Value string
}

func (sl *StringLiteral) expressionNode()        {}
func (sl *StringLiteral) GetToken() tokens.Token { return sl.Token }
func (sl *StringLiteral) TokenLiteral() string   { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("%q", sl.Value)
}

type NullLiteral struct {
	Token tokens.Token
}

func (nl *NullLiteral) expressionNode()        {}
func (nl *NullLiteral) GetToken() tokens.Token { return nl.Token }
func (nl *NullLiteral) TokenLiteral() string   { return nl.Token.Literal }
func (nl *NullLiteral) String() string         { return "null" }

type BooleanLiteral struct {
	Token tokens.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()        {}
func (b *BooleanLiteral) GetToken() tokens.Token { return b.Token }
func (b *BooleanLiteral) TokenLiteral() string   { return b.Token.Literal }
func (b *BooleanLiteral) String() string         { return b.Token.Literal }

type ArrayLiteral struct {
	Token    tokens.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()        {}
func (al *ArrayLiteral) GetToken() tokens.Token { return al.Token }
func (al *ArrayLiteral) TokenLiteral() string   { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct {
	Token tokens.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()        {}
func (hl *HashLiteral) GetToken() tokens.Token { return hl.Token }
func (hl *HashLiteral) TokenLiteral() string   { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, fmt.Sprintf("%q: %s", key.String(), value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type IndexExpression struct {
	Token tokens.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()        {}
func (ie *IndexExpression) GetToken() tokens.Token { return ie.Token }
func (ie *IndexExpression) TokenLiteral() string   { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()        {}
func (pe *PrefixExpression) GetToken() tokens.Token { return pe.Token }
func (pe *PrefixExpression) TokenLiteral() string   { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()        {}
func (ie *InfixExpression) GetToken() tokens.Token { return ie.Token }
func (ie *InfixExpression) TokenLiteral() string   { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type SuffixExpression struct {
	Token    tokens.Token
	Operator string
	Left     Expression
}

func (se *SuffixExpression) expressionNode()        {}
func (se *SuffixExpression) GetToken() tokens.Token { return se.Token }
func (se *SuffixExpression) TokenLiteral() string   { return se.Token.Literal }
func (se *SuffixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString(se.Operator)
	out.WriteString(")")

	return out.String()
}

type ChainExpression struct {
	Token tokens.Token
	Left  Expression
	Right Expression
}

func (ce *ChainExpression) expressionNode()        {}
func (ce *ChainExpression) GetToken() tokens.Token { return ce.Token }
func (ce *ChainExpression) TokenLiteral() string   { return ce.Token.Literal }
func (ce *ChainExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ce.Left.String())
	out.WriteString(".")
	out.WriteString(ce.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token        tokens.Token
	Condition    Expression
	Consequence  *BlockStatement
	Intermediary *IfExpression
	Alternative  *BlockStatement
}

func (ie *IfExpression) expressionNode()        {}
func (ie *IfExpression) GetToken() tokens.Token { return ie.Token }
func (ie *IfExpression) TokenLiteral() string   { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Condition.String())
	out.WriteString(") { ")
	out.WriteString(ie.Consequence.String())
	out.WriteString(" }")

	if ie.Intermediary != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Intermediary.String())
	}

	if ie.Alternative != nil {
		out.WriteString(" else { ")
		out.WriteString(ie.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

type WhileExpression struct {
	Token     tokens.Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileExpression) expressionNode()        {}
func (w *WhileExpression) GetToken() tokens.Token { return w.Token }
func (w *WhileExpression) TokenLiteral() string   { return w.Token.Literal }
func (w *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while (")
	out.WriteString(w.Condition.String())
	out.WriteString(") { ")
	out.WriteString(w.Body.String())
	out.WriteString(" }")

	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()         {}
func (bs *BlockStatement) GetToken() tokens.Token { return bs.Token }
func (bs *BlockStatement) TokenLiteral() string   { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      tokens.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()        {}
func (fl *FunctionLiteral) GetToken() tokens.Token { return fl.Token }
func (fl *FunctionLiteral) TokenLiteral() string   { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(" ")

	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")

	return out.String()
}

type CallExpression struct {
	Token     tokens.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()        {}
func (ce *CallExpression) GetToken() tokens.Token { return ce.Token }
func (ce *CallExpression) TokenLiteral() string   { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type AssignmentExpression struct {
	Token tokens.Token
	Left  Expression
	Right Expression
}

func (ae *AssignmentExpression) expressionNode()        {}
func (ae *AssignmentExpression) GetToken() tokens.Token { return ae.Token }
func (ae *AssignmentExpression) TokenLiteral() string   { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Right.String())

	return out.String()
}

type ImportStatement struct {
	Token   tokens.Token
	Path    string
	Aliased *Identifier
}

func (is *ImportStatement) statementNode()         {}
func (is *ImportStatement) GetToken() tokens.Token { return is.Token }
func (is *ImportStatement) TokenLiteral() string   { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer

	out.WriteString("import '")
	out.WriteString(is.Path)
	out.WriteString("'")

	if is.Aliased != nil {
		out.WriteString(" as ")
		out.WriteString(is.Aliased.String())
	}

	return out.String()
}

type ExportStatement struct {
	Token tokens.Token
	Value Expression
}

func (es *ExportStatement) statementNode()         {}
func (es *ExportStatement) GetToken() tokens.Token { return es.Token }
func (es *ExportStatement) TokenLiteral() string   { return es.Token.Literal }
func (es *ExportStatement) String() string {
	var out bytes.Buffer

	out.WriteString("export ")
	out.WriteString(es.Value.String())
	out.WriteString(";")

	return out.String()
}

type BreakStatement struct {
	Token tokens.Token
}

func (bs *BreakStatement) statementNode()         {}
func (bs *BreakStatement) GetToken() tokens.Token { return bs.Token }
func (bs *BreakStatement) TokenLiteral() string   { return bs.Token.Literal }
func (bs *BreakStatement) String() string         { return "break;" }

type ContinueStatement struct {
	Token tokens.Token
}

func (cs *ContinueStatement) statementNode()         {}
func (cs *ContinueStatement) GetToken() tokens.Token { return cs.Token }
func (cs *ContinueStatement) TokenLiteral() string   { return cs.Token.Literal }
func (cs *ContinueStatement) String() string         { return "continue;" }
