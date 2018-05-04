package builtin

import (
	"github.com/Bitspark/slang/pkg/core"
)

var syncForkOpCfg = &builtinConfig{
	oDef: core.OperatorDef{
		Services: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.PortDef{
					Type: "map",
					Map: map[string]*core.PortDef{
						"item": {
							Type:    "generic",
							Generic: "itemType",
						},
						"select": {
							Type: "boolean",
						},
					},
				},
				Out: core.PortDef{
					Type: "map",
					Map: map[string]*core.PortDef{
						"true": {
							Type:    "generic",
							Generic: "itemType",
						},
						"false": {
							Type:    "generic",
							Generic: "itemType",
						},
					},
				},
			},
		},
	},
	oFunc: func(srvs map[string]*core.Service, dels map[string]*core.Delegate, store interface{}) {
		in := srvs[core.MAIN_SERVICE].In()
		out := srvs[core.MAIN_SERVICE].Out()
		for true {
			item := in.Pull()
			m, ok := item.(map[string]interface{})
			if !ok {
				out.Push(item)
				continue
			}

			if m["select"].(bool) {
				out.Map("true").Push(m["item"])
				out.Map("false").Push(nil)
			} else {
				out.Map("true").Push(nil)
				out.Map("false").Push(m["item"])
			}
		}
	},
}