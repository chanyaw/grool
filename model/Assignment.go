package model

import (
	"fmt"
	"github.com/newm4n/grool/context"
)

type Assignment struct {
	Variable         string
	Expression       *Expression
	knowledgeContext *context.KnowledgeContext
	ruleCtx          *context.RuleContext
	dataCtx          *context.DataContext
}

func (ins *Assignment) Initialize(knowledgeContext *context.KnowledgeContext, ruleCtx *context.RuleContext, dataCtx *context.DataContext) {
	ins.knowledgeContext = knowledgeContext
	ins.ruleCtx = ruleCtx
	ins.dataCtx = dataCtx

	if ins.Expression != nil {
		ins.Expression.Initialize(knowledgeContext, ruleCtx, dataCtx)
	}
}

func (assign *Assignment) AcceptExpression(expression *Expression) error {
	if assign.Expression != nil {
		return fmt.Errorf("expression were set twice in assignment")
	}
	assign.Expression = expression
	return nil
}

func (assign *Assignment) AcceptVariable(name string) error {
	if assign.Variable == "" {
		assign.Variable = name
		return nil
	} else {
		return fmt.Errorf("variable already defined")
	}
}