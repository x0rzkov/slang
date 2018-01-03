package op

import (
	"errors"
	"fmt"
	"strings"
)

type OFunc func(in, out *Port, store interface{})

type Operator struct {
	name     string
	basePort *Port
	inPort   *Port
	outPort  *Port
	parent   *Operator
	children map[string]*Operator
	function OFunc
	store    interface{}
}

type OperatorDef struct {
	Name        string              `json:"name"`
	In          *PortDef            `json:"in"`
	Out         *PortDef            `json:"out"`
	Operators   []InstanceDef       `json:"operators"`
	Connections map[string][]string `json:"connections"`
	valid       bool
}

type InstanceDef struct {
	Operator   string                 `json:"operator"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	In         *PortDef               `json:"in"`
	Out        *PortDef               `json:"out"`
	valid      bool
}

func MakeOperator(name string, f OFunc, defIn, defOut PortDef, par *Operator) (*Operator, error) {
	o := &Operator{}
	o.function = f
	o.parent = par
	o.name = name
	o.children = make(map[string]*Operator)

	if par != nil {
		par.children[o.name] = o
	}

	var err error

	o.inPort, err = MakePort(o, defIn, DIRECTION_IN)
	if err != nil {
		return nil, err
	}

	o.outPort, err = MakePort(o, defOut, DIRECTION_OUT)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (d OperatorDef) Valid() bool {
	return d.valid
}

func (d InstanceDef) Valid() bool {
	return d.valid
}

func (d *OperatorDef) Validate() error {
	if d.Name == "" {
		return errors.New(`operator name may not be empty`)
	}

	if strings.Contains(d.Name, " ") {
		return fmt.Errorf(`operator name may not contain spaces: "%s"`, d.Name)
	}

	if d.In == nil || d.Out == nil {
		return errors.New(`ports must be defined`)
	}

	if err := d.In.Validate(); err != nil {
		return err
	}

	if err := d.Out.Validate(); err != nil {
		return err
	}

	alreadyUsedInsNames := make(map[string]bool)
	for _, insDef := range d.Operators {
		if err := insDef.Validate(); err != nil {
			return err
		}

		if _, ok := alreadyUsedInsNames[insDef.Name]; ok {
			return fmt.Errorf(`Colliding instance names within same parent operator: "%s"`, insDef.Name)
		}

		alreadyUsedInsNames[insDef.Name] = true

	}

	d.valid = true
	return nil
}

func (d *InstanceDef) Validate() error {
	if d.Name == "" {
		return fmt.Errorf(`instance name may not be empty`)
	}

	if strings.Contains(d.Name, " ") {
		return fmt.Errorf(`operator instance name may not contain spaces: "%s"`, d.Name)
	}

	if d.Operator == "" {
		return errors.New(`operator may not be empty`)
	}

	if strings.Contains(d.Operator, " ") {
		return fmt.Errorf(`operator may not contain spaces: "%s"`, d.Operator)
	}

	if d.In != nil {
		if err := d.In.Validate(); err != nil {
			return err
		}
	}

	if d.Out != nil {
		if err := d.Out.Validate(); err != nil {
			return err
		}
	}

	d.valid = true
	return nil
}

func (o *Operator) InPort() *Port {
	return o.inPort
}

func (o *Operator) OutPort() *Port {
	return o.outPort
}

func (o *Operator) Name() string {
	return o.name
}

func (o *Operator) BasePort() *Port {
	return o.basePort
}

func (o *Operator) Parent() *Operator {
	return o.parent
}

func (o *Operator) Child(name string) *Operator {
	c, _ := o.children[name]
	return c
}

func (o *Operator) Start() {
	o.function(o.inPort, o.outPort, o.store)
}

func (o *Operator) SetStore(store interface{}) {
	o.store = store
}

func (o *Operator) Compile() bool {
	compiled := o.InPort().Merge()
	compiled = o.OutPort().Merge() || compiled
	return compiled
}