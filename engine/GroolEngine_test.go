package engine

import (
	"github.com/newm4n/grool/builder"
	"github.com/newm4n/grool/context"
	"github.com/newm4n/grool/model"
	"github.com/newm4n/grool/pkg"
	"reflect"
	"sort"
	"testing"
	"time"
)

type Sorting struct {
	Val int
}

func TestGroolSorting(t *testing.T) {
	arr := make([]*Sorting, 0)
	arr = append(arr, &Sorting{Val: 4})
	arr = append(arr, &Sorting{Val: 7})
	arr = append(arr, &Sorting{Val: 3})
	arr = append(arr, &Sorting{Val: 6})
	arr = append(arr, &Sorting{Val: 9})
	arr = append(arr, &Sorting{Val: 8})
	arr = append(arr, &Sorting{Val: 1})
	arr = append(arr, &Sorting{Val: 2})

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Val > arr[j].Val
	})

	if arr[0].Val != 9 {
		t.FailNow()
	}
}

type TestCar struct {
	SpeedUp        bool
	Speed          int
	MaxSpeed       int
	SpeedIncrement int
}

type DistanceRecorder struct {
	TotalDistance int
	TestTime      time.Time
}

const (
	rules = `
rule SpeedUp "When testcar is speeding up we keep increase the speed." salience 10 {
    when
        TestCar.SpeedUp == true && TestCar.Speed < TestCar.MaxSpeed
    then
        TestCar.Speed = TestCar.Speed + TestCar.SpeedIncrement;
		DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}

rule StartSpeedDown "When testcar is speeding up and over max speed we change to speed down." salience 10  {
    when
        TestCar.SpeedUp == true && TestCar.Speed >= TestCar.MaxSpeed
    then
        TestCar.SpeedUp = false;
		Log("Now we slow down");
}

rule SlowDown "When testcar is slowing down we keep decreasing the speed." salience 10  {
    when
        TestCar.SpeedUp == false && TestCar.Speed > 0
    then
        TestCar.Speed = TestCar.Speed - TestCar.SpeedIncrement;
		DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}

rule SetTime "When Distance Recorder time not set, set it." {
	when
		IsZero(DistanceRecord.TestTime)
	then
		Log("Set the test time");
		DistanceRecord.TestTime = Now();
		Log(TimeFormat(DistanceRecord.TestTime,"Mon Jan _2 15:04:05 2006"));
}
`
)

func TestGrool_Execute(t *testing.T) {
	tc := &TestCar{
		SpeedUp:        true,
		Speed:          0,
		MaxSpeed:       100,
		SpeedIncrement: 2,
	}
	dr := &DistanceRecorder{
		TotalDistance: 0,
	}
	dctx := context.NewDataContext()
	dctx.Add("TestCar", tc)
	dctx.Add("DistanceRecord", dr)

	kb := model.NewKnowledgeBase()
	rb := builder.NewRuleBuilder(kb)
	err := rb.BuildRuleFromResource(pkg.NewBytesResource([]byte(rules)))
	if err != nil {
		t.Errorf("Got error : %v", err)
		t.FailNow()
	} else {
		engine := NewGroolEngine()
		start := time.Now()
		err = engine.Execute(dctx, kb)
		if err != nil {
			t.Errorf("Got error : %v", err)
			t.FailNow()
		} else {
			dur := time.Since(start)
			t.Log(dr.TotalDistance)
			t.Logf("Duration %f ms", float64(dur)/float64(time.Millisecond))
		}
	}
}

func TestEmptyValueEquality(t *testing.T) {
	t1 := time.Time{}
	tv1 := reflect.ValueOf(t1)
	tv2 := reflect.Zero(tv1.Type())

	if tv1.Type() != tv2.Type() {
		t.Logf("%s vs %s", tv1.Type().String(), tv2.Type().String())
		t.FailNow()
	}

	if pkg.ValueToInterface(tv1) != pkg.ValueToInterface(tv2) {
		t.Logf("%s vs %s", tv1.Kind().String(), tv2.Kind().String())
		t.Logf("%s vs %s", tv1.Type().String(), tv2.Type().String())
		t.Logf("%v vs %v", tv1.IsValid(), tv2.IsValid())

		t.FailNow()
	}
}
