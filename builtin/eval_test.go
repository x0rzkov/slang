package builtin

import (
	"slang/core"
	"slang/tests/assertions"
	"testing"

	"github.com/Knetic/govaluate"
)

func TestEvaluableExpression__Translates_Variables(t *testing.T) {
	a := assertions.New(t)

	evalExpr, _ := NewEvaluableExpression("a + 100")

	a.Equal([]string{"a"}, evalExpr.Vars())
}

func TestEvaluableExpression__Translates_AccessingFields(t *testing.T) {
	a := assertions.New(t)

	evalExpr, _ := NewEvaluableExpression("vec.x + 100")

	a.Equal([]string{"vec__x"}, evalExpr.Vars())
}

func TestEvaluableExpression__Translates_AccessingNestedFields(t *testing.T) {
	a := assertions.New(t)

	evalExpr, _ := NewEvaluableExpression(`person.address.zipcode == "ABCD"`)

	a.Equal([]string{"person__address__zipcode"}, evalExpr.Vars())
}

func TestEvaluableExpression__Evaluates_ArgumentsAcessingFields(t *testing.T) {
	a := assertions.New(t)

	evalExpr, _ := NewEvaluableExpression("(vec.0.x + vec.1.x) * (vec.0.y + vec.1.y)")
	params := NewFlatMapParameters(map[string]interface{}{
		"vec": []interface{}{
			map[string]interface{}{
				"x": 1,
				"y": 2,
			},
			map[string]interface{}{
				"x": 3,
				"y": 4,
			},
		},
	})
	result, _ := evalExpr.Evaluate(params)
	a.Equal(24.0, result)
}

func TestEvaluableExpression__Translates_Combined(t *testing.T) {
	a := assertions.New(t)

	evalExpr, _ := NewEvaluableExpression(`a * vec.x + vec.y`)

	a.Equal([]string{"a", "vec__x", "vec__y"}, evalExpr.Vars())
}

func TestFlatMapParameters__NestingLevel0(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
	}), params)
}

func TestFlatMapParameters__NestingLevel0_Arrays(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
		"l":   []interface{}{1, 2},
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo":  100,
		"bar":  200,
		"l__0": 1,
		"l__1": 2,
	}), params)
}

func TestFlatMapParameters__NestingLevel1(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
		"vec": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo":    100,
		"bar":    200,
		"vec__x": 1,
		"vec__y": 2,
	}), params)
}

func TestFlatMapParameters__NestingLevel2(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
		"vec": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
		"person": map[string]interface{}{
			"phone": map[string]interface{}{
				"mobile": "0123/4567890",
			},
			"name": map[string]interface{}{
				"first": "Paul J.",
				"last":  "Morrison",
			},
		},
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo":                   100,
		"bar":                   200,
		"vec__x":                1,
		"vec__y":                2,
		"person__phone__mobile": "0123/4567890",
		"person__name__first":   "Paul J.",
		"person__name__last":    "Morrison",
	}), params)
}

func TestFlatMapParameters__NestingLevel3(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
		"vec": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
		"root": map[string]interface{}{
			"left": map[string]interface{}{
				"left": map[string]interface{}{
					"i": 100,
				},
				"right": map[string]interface{}{

					"i": 101,
				},
			},
			"right": map[string]interface{}{
				"left": map[string]interface{}{
					"i": 110,
				},
				"right": map[string]interface{}{

					"i": 111,
				},
			},
		},
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo":                   100,
		"bar":                   200,
		"vec__x":                1,
		"vec__y":                2,
		"root__left__left__i":   100,
		"root__left__right__i":  101,
		"root__right__left__i":  110,
		"root__right__right__i": 111,
	}), params)
}

func TestFlatMapParameters__ComplexMixed(t *testing.T) {
	a := assertions.New(t)

	params := NewFlatMapParameters(map[string]interface{}{
		"foo": 100,
		"bar": 200,
		"vec": []interface{}{
			map[string]interface{}{
				"x": 1,
				"y": 2,
			},
			map[string]interface{}{
				"x": 3,
				"y": 4,
			},
		},
		"person": []interface{}{
			map[string]interface{}{
				"phone": map[string]interface{}{
					"mobile": "0123/4567890",
				},
				"name": map[string]interface{}{
					"first": "Taleh",
					"last":  "Didover",
				},
			},
			map[string]interface{}{
				"phone": map[string]interface{}{
					"mobile": "0123/4567890",
				},
				"name": map[string]interface{}{
					"first": "Julian",
					"last":  "Matschinske",
				},
			},
		},
	})

	a.Equal(govaluate.MapParameters(map[string]interface{}{
		"foo":                      100,
		"bar":                      200,
		"vec__0__x":                1,
		"vec__0__y":                2,
		"vec__1__x":                3,
		"vec__1__y":                4,
		"person__0__phone__mobile": "0123/4567890",
		"person__0__name__first":   "Taleh",
		"person__0__name__last":    "Didover",
		"person__1__phone__mobile": "0123/4567890",
		"person__1__name__first":   "Julian",
		"person__1__name__last":    "Matschinske",
	}), params)
}

func TestBuiltin_Eval__CreatorFuncIsRegistered(t *testing.T) {
	a := assertions.New(t)

	ocEval := getCreatorFunc("eval")
	a.NotNil(ocEval)
}

func TestBuiltin_Eval__NilProperties(t *testing.T) {
	a := assertions.New(t)
	_, err := createOpEval(core.InstanceDef{Operator: "eval"})
	a.Error(err)
}

func TestBuiltin_Eval__EmptyExpression(t *testing.T) {
	a := assertions.New(t)
	_, err := createOpEval(core.InstanceDef{Operator: "eval", Properties: map[string]interface{}{"expression": ""}})
	a.Error(err)
}

func TestBuiltin_Eval__InvalidExpression(t *testing.T) {
	a := assertions.New(t)
	_, err := createOpEval(core.InstanceDef{Operator: "eval", Properties: map[string]interface{}{"expression": "+"}})
	a.Error(err)
}

func TestBuiltin_Eval__NilIn(t *testing.T) {
	a := assertions.New(t)
	_, err := createOpEval(core.InstanceDef{Operator: "eval", Properties: map[string]interface{}{"expression": "100"}})
	a.Error(err)
}

func TestBuiltin_Eval__NilOut(t *testing.T) {
	a := assertions.New(t)
	_, err := createOpEval(core.InstanceDef{Operator: "eval", Properties: map[string]interface{}{"expression": "100"}})
	a.Error(err)
}

func TestBuiltin_Eval__Add(t *testing.T) {
	a := assertions.New(t)
	fo, err := createOpEval(core.InstanceDef{
		Operator:   "eval",
		Properties: map[string]interface{}{"expression": "a+b"},
		In: &core.PortDef{
			Type: "map",
			Map: map[string]core.PortDef{
				"a": {Type: "number"},
				"b": {Type: "number"},
			},
		},
		Out: &core.PortDef{Type: "number"},
	})
	a.NoError(err)
	a.NotNil(fo)
	fo.Out().Bufferize()

	go fo.Start()

	fo.In().Push(map[string]interface{}{"a": 1.0, "b": 2.0})
	fo.In().Push(map[string]interface{}{"a": -5.0, "b": 2.5})
	fo.In().Push(map[string]interface{}{"a": 0.0, "b": 333.0})

	a.PortPushes([]interface{}{3.0, -2.5, 333.0}, fo.Out())
}

func TestBuiltin_Eval__BoolArith(t *testing.T) {
	a := assertions.New(t)
	fo, err := createOpEval(core.InstanceDef{
		Operator:   "eval",
		Properties: map[string]interface{}{"expression": "a && (b != c)"},
		In: &core.PortDef{
			Type: "map",
			Map: map[string]core.PortDef{
				"a": {Type: "boolean"},
				"b": {Type: "boolean"},
				"c": {Type: "boolean"},
			},
		},
		Out: &core.PortDef{Type: "boolean"},
	})
	a.NoError(err)
	a.NotNil(fo)
	fo.Out().Bufferize()

	go fo.Start()

	fo.In().Push(map[string]interface{}{"a": true, "b": true, "c": false})
	fo.In().Push(map[string]interface{}{"a": false, "b": false, "c": false})
	fo.In().Push(map[string]interface{}{"a": false, "b": false, "c": true})
	fo.In().Push(map[string]interface{}{"a": true, "b": false, "c": true})
	fo.In().Push(map[string]interface{}{"a": true, "b": false, "c": false})

	a.PortPushes([]interface{}{true, false, false, true, false}, fo.Out())
}

func TestBuiltin_Eval_VectorArith(t *testing.T) {
	a := assertions.New(t)
	fo, err := createOpEval(core.InstanceDef{
		Operator:   "eval",
		Properties: map[string]interface{}{"expression": "vec0.x*vec1.x+vec0.y*vec1.y"},
		In: &core.PortDef{
			Type: "map",
			Map: map[string]core.PortDef{
				"vec0": {
					Type: "map",
					Map: map[string]core.PortDef{
						"x": {Type: "number"},
						"y": {Type: "number"},
					},
				},
				"vec1": {
					Type: "map",
					Map: map[string]core.PortDef{
						"x": {Type: "number"},
						"y": {Type: "number"},
					},
				},
			},
		},
		Out: &core.PortDef{Type: "boolean"},
	})
	a.NoError(err)
	a.NotNil(fo)
	fo.Out().Bufferize()

	go fo.Start()

	fo.In().Push(map[string]interface{}{
		"vec0": map[string]interface{}{
			"x": 2,
			"y": 4,
		},
		"vec1": map[string]interface{}{
			"x": 3,
			"y": 5,
		},
	})
	fo.In().Push(map[string]interface{}{
		"vec0": map[string]interface{}{
			"x": 10,
			"y": 0,
		},
		"vec1": map[string]interface{}{
			"x": 0,
			"y": 10,
		},
	})
	fo.In().Push(map[string]interface{}{
		"vec0": map[string]interface{}{
			"x": 1,
			"y": 1,
		},
		"vec1": map[string]interface{}{
			"x": 1,
			"y": 1,
		},
	})

	a.PortPushes([]interface{}{26., 0., 2.}, fo.Out())
}