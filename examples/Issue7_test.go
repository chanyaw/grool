package examples

import (
	"github.com/newm4n/grool/builder"
	"github.com/newm4n/grool/context"
	"github.com/newm4n/grool/engine"
	"github.com/newm4n/grool/model"
	"github.com/newm4n/grool/pkg"
	"testing"
)

const (
	Rule7 = `
rule UserTestRule7 "test 7"  salience 10{
	when
	  User.Age > 1
	then
	  User.SetName("FromRule");
	  Retract("UserTestRule7");
}
`
)

type AUserIssue7 struct {
	Name string
	Age  int
}

func (u *AUserIssue7) GetName() string {
	return u.Name
}

func (u *AUserIssue7) SetName(name interface{}) {
	u.Name = name.(string)
}

func TestMethodCall_Issue7(t *testing.T) {
	user := &AUserIssue7{
		Name: "Watson",
		Age:  7,
	}

	dataContext := context.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		t.Fatal(err)
	}

	knowledgeBase := model.NewKnowledgeBase()
	ruleBuilder := builder.NewRuleBuilder(knowledgeBase)

	err = ruleBuilder.BuildRuleFromResource(pkg.NewBytesResource([]byte(Rule7)))
	if err != nil {
		t.Log(err)
	} else {
		eng1 := &engine.Grool{MaxCycle: 5}
		err := eng1.Execute(dataContext, knowledgeBase)
		if err != nil {
			t.Fatal(err)
		}
		if user.GetName() != "FromRule" {
			t.Errorf("User should be FromRule but %s", user.GetName())
		}
	}
}
