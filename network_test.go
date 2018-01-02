package main

import (
	"testing"
)

func TestNetwork_Complex1(t *testing.T) {
	defStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}`)
	defStr := helperJson2Map(`{"type":"stream","stream":{"type":"number"}}`)
	def := helperJson2Map(`{"type":"number"}`)

	dummy := func(in, out *Port) {}

	o2, _ := MakeOperator("O2", nil, defStr, def, nil)
	o3, _ := MakeOperator("O3", dummy, def, def, o2)
	o4, _ := MakeOperator("O4", dummy, defStr, def, o2)

	err := o2.InPort().Stream().Connect(o3.InPort())
	assertNoError(t, err)
	err = o3.OutPort().Connect(o4.InPort().Stream())
	assertNoError(t, err)
	err = o4.OutPort().Connect(o2.OutPort())
	assertNoError(t, err)

	if !o2.InPort().Stream().Connected(o3.InPort()) {
		t.Error("should be connected")
	}

	if !o3.OutPort().Connected(o4.InPort().Stream()) {
		t.Error("should be connected")
	}

	if !o4.OutPort().Connected(o2.OutPort()) {
		t.Error("should be connected")
	}

	if o3.BasePort() != o2.InPort() {
		t.Error("wrong base port")
	}

	if !o2.InPort().Connected(o4.InPort()) {
		t.Error("should be connected via base port")
	}

	//

	o1, _ := MakeOperator("O1", nil, defStrStr, defStr, nil)
	o2.parent = o1

	err = o1.InPort().Stream().Stream().Connect(o2.InPort().Stream())
	assertNoError(t, err)
	err = o2.OutPort().Connect(o1.OutPort().Stream())
	assertNoError(t, err)

	if !o1.InPort().Stream().Stream().Connected(o2.InPort().Stream()) {
		t.Error("should be connected")
	}

	if !o1.InPort().Stream().Connected(o2.InPort()) {
		t.Error("should be connected")
	}

	if !o2.OutPort().Connected(o1.OutPort().Stream()) {
		t.Error("should be connected")
	}

	if o2.BasePort() != o1.InPort() {
		t.Error("wrong base port")
	}

	if !o1.InPort().Connected(o1.OutPort()) {
		t.Error("should be connected via base port")
	}

	//

	o2.Compile()

	if !o1.InPort().Connected(o1.OutPort()) {
		t.Error("should be connected")
	}

	if !o1.InPort().Stream().Connected(o4.InPort()) {
		t.Error("should be connected")
	}

	if !o1.InPort().Stream().Stream().Connected(o3.InPort()) {
		t.Error("should be connected")
	}

	if !o3.OutPort().Connected(o4.InPort().Stream()) {
		t.Error("should be connected")
	}
}

func TestNetwork_Complex2(t *testing.T) {
	defStrStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}}`)
	defStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}`)
	defStr := helperJson2Map(`{"type":"stream","stream":{"type":"number"}}`)
	def := helperJson2Map(`{"type":"number"}`)

	dummy := func(in, out *Port) {}

	o4, _ := MakeOperator("O4", nil, defStrStrStr, defStrStr, nil)
	o5, _ := MakeOperator("O5", dummy, defStr, def, o4)

	o4.InPort().Stream().Stream().Stream().Connect(o5.InPort().Stream())
	o5.OutPort().Connect(o4.OutPort().Stream().Stream())

	if !o4.InPort().Stream().Stream().Stream().Connected(o5.InPort().Stream()) {
		t.Error("should be connected")
	}

	if !o4.InPort().Stream().Stream().Connected(o5.InPort()) {
		t.Error("should be connected")
	}

	if !o5.OutPort().Connected(o4.OutPort().Stream().Stream()) {
		t.Error("should be connected")
	}

	if !o4.InPort().Stream().Connected(o4.OutPort().Stream()) {
		t.Error("should be connected via base port")
	}

	if !o4.InPort().Connected(o4.OutPort()) {
		t.Error("should be connected via base port")
	}

	//

	o1, _ := MakeOperator("O1", nil, defStr, defStrStr, nil)
	o2, _ := MakeOperator("O2", dummy, def, defStr, o1)
	o3, _ := MakeOperator("O3", dummy, def, defStr, o1)
	o4.parent = o1

	o1.InPort().Stream().Connect(o2.InPort())
	o2.OutPort().Stream().Connect(o3.InPort())
	o3.OutPort().Stream().Connect(o4.InPort().Stream().Stream().Stream())
	o4.OutPort().Stream().Stream().Connect(o1.OutPort().Stream().Stream())

	if !o1.InPort().Stream().Connected(o2.InPort()) {
		t.Error("should be connected")
	}

	if !o2.OutPort().Stream().Connected(o3.InPort()) {
		t.Error("should be connected")
	}

	if !o3.OutPort().Stream().Connected(o4.InPort().Stream().Stream().Stream()) {
		t.Error("should be connected")
	}

	if !o3.OutPort().Connected(o4.InPort().Stream().Stream()) {
		t.Error("should be connected")
	}

	if !o4.OutPort().Stream().Stream().Connected(o1.OutPort().Stream().Stream()) {
		t.Error("should be connected")
	}

	if !o4.OutPort().Stream().Connected(o1.OutPort().Stream()) {
		t.Error("should be connected")
	}

	if !o4.OutPort().Connected(o1.OutPort()) {
		t.Error("should be connected")
	}

	if o2.BasePort() != o1.InPort() {
		t.Error("wrong base port")
	}

	if o3.BasePort() != o2.OutPort() {
		t.Error("wrong base port")
	}

	if !o1.InPort().Connected(o4.InPort()) {
		t.Error("should be connected via base port")
	}

	if !o2.OutPort().Connected(o4.InPort().Stream()) {
		t.Error("should be connected via base port")
	}

	//

	o4.Compile()

	if !o1.InPort().Connected(o1.OutPort()) {
		t.Error("should be connected after merge")
	}

	if !o2.OutPort().Connected(o1.OutPort().Stream()) {
		t.Error("should be connected after merge")
	}

	if !o3.OutPort().Stream().Connected(o5.InPort().Stream()) {
		t.Error("should be connected after merge")
	}

	if !o3.OutPort().Connected(o5.InPort()) {
		t.Error("should be connected after merge")
	}
}

func TestNetwork_Complex1_PushPull(t *testing.T) {
	defStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}`)
	defStr := helperJson2Map(`{"type":"stream","stream":{"type":"number"}}`)
	def := helperJson2Map(`{"type":"number"}`)

	double := func(in, out *Port) {
		for true {
			i := in.Pull()
			if n, ok := i.(float64); ok {
				out.Push(2 * n)
			} else {
				out.Push(i)
			}
		}
	}

	sum := func(in, out *Port) {
		for true {
			i := in.Pull()
			if ns, ok := i.([]interface{}); ok {
				sum := 0.0
				for _, n := range ns {
					sum += n.(float64)
				}
				out.Push(sum)
			} else {
				out.Push(i)
			}
		}
	}

	o1, _ := MakeOperator("O1", nil, defStrStr, defStr, nil)
	o2, _ := MakeOperator("O2", nil, defStr, def, o1)
	o3, _ := MakeOperator("O3", double, def, def, o2)
	o4, _ := MakeOperator("O4", sum, defStr, def, o2)

	o1.InPort().Stream().Stream().Connect(o2.InPort().Stream())
	o2.InPort().Stream().Connect(o3.InPort())
	o3.OutPort().Connect(o4.InPort().Stream())
	o4.OutPort().Connect(o2.OutPort())
	o2.OutPort().Connect(o1.OutPort().Stream())

	o2.Compile()

	o1.OutPort().Stream().buf = make(chan interface{}, 100)

	go o3.Start()
	go o4.Start()

	o1.InPort().Push(helperJson2I(`[[1,2,3],[4,5]]`))
	o1.InPort().Push(helperJson2I(`[[],[2]]`))
	o1.InPort().Push(helperJson2I(`[]`))
	assertPortItems(t, helperJson2I(`[[12,18],[0,4],[]]`).([]interface{}), o1.OutPort())
}

func TestNetwork_Complex2_PushPull(t *testing.T) {
	defStrStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}}`)
	defStrStr := helperJson2Map(`{"type":"stream","stream":{"type":"stream","stream":{"type":"number"}}}`)
	defStr := helperJson2Map(`{"type":"stream","stream":{"type":"number"}}`)
	def := helperJson2Map(`{"type":"number"}`)

	numgen := func(in, out *Port) {
		for true {
			i := in.Pull()
			if n, ok := i.(float64); ok {
				ns := []interface{}{}
				for i := 1; i <= int(n); i++ {
					ns = append(ns, float64(i))
				}
				out.Push(ns)
			} else {
				out.Push(i)
			}
		}
	}

	sum := func(in, out *Port) {
		for true {
			i := in.Pull()
			if ns, ok := i.([]interface{}); ok {
				sum := 0.0
				for _, n := range ns {
					sum += n.(float64)
				}
				out.Push(sum)
			} else {
				out.Push(i)
			}
		}
	}

	o1, _ := MakeOperator("O1", nil, defStr, defStrStr, nil)
	o2, _ := MakeOperator("O2", numgen, def, defStr, o1)
	o3, _ := MakeOperator("O3", numgen, def, defStr, o1)
	o4, _ := MakeOperator("O4", nil, defStrStrStr, defStrStr, o1)
	o5, _ := MakeOperator("O5", sum, defStr, def, o4)

	o1.InPort().Stream().Connect(o2.InPort())
	o2.OutPort().Stream().Connect(o3.InPort())
	o3.OutPort().Stream().Connect(o4.InPort().Stream().Stream().Stream())
	o4.InPort().Stream().Stream().Stream().Connect(o5.InPort().Stream())
	o5.OutPort().Connect(o4.OutPort().Stream().Stream())
	o4.OutPort().Stream().Stream().Connect(o1.OutPort().Stream().Stream())

	//

	o4.Compile()

	o1.OutPort().Stream().Stream().buf = make(chan interface{}, 100)

	go o2.Start()
	go o3.Start()
	go o5.Start()

	o1.InPort().Push(helperJson2I(`[1,2,3]`))
	o1.InPort().Push(helperJson2I(`[]`))
	o1.InPort().Push(helperJson2I(`[4]`))
	assertPortItems(t, helperJson2I(`[[[1],[1,3],[1,3,6]],[],[[1,3,6,10]]]`).([]interface{}), o1.OutPort())
}
