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
	Rule4 = `
rule UserTestRule4 "test 3"  salience 10{
	when
	  User.Auth.GetEmail() == "watson@test.com"
	then
	  User.Name = "FromRuleScope4";
	  Retract("UserTestRule4");
}
`
)

type UserWithAuth struct {
	Auth *UserAuth
	Name string
}

func (user *UserWithAuth) GetName() string {
	return user.Name
}

type UserAuth struct {
	Email string
}

func (auth *UserAuth) GetEmail() string {
	return auth.Email
}

func TestMethodCall_Issue4(t *testing.T) {
	user := &UserWithAuth{
		Auth: &UserAuth{Email: "watson@test.com"},
	}

	if user.GetName() != "" {
		t.Fatal("User name not empty")
	}

	dataContext := context.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		t.Fatal(err)
	}

	knowledgeBase := model.NewKnowledgeBase()
	ruleBuilder := builder.NewRuleBuilder(knowledgeBase)

	err = ruleBuilder.BuildRuleFromResource(pkg.NewBytesResource([]byte(Rule4)))
	if err != nil {
		t.Log(err)
	} else {
		eng1 := &engine.Grool{MaxCycle: 3}
		err := eng1.Execute(dataContext, knowledgeBase)
		if err != nil {
			t.Fatal(err)
		}
		if user.GetName() != "FromRuleScope4" {
			t.Errorf("User should be FromRuleScope4 but %s", user.GetName())
		}
	}
}
