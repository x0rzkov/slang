package elem

import (
	"github.com/Bitspark/slang/pkg/core"
)

var streamStreamToMapCfg = &builtinConfig{
	opDef: core.OperatorDef{
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "stream",
					Stream: &core.TypeDef{
						Type: "map",
						Map: map[string]*core.TypeDef{
							"{key}": {
								Type:    "generic",
								Generic: "keyType",
							},
							"{value}": {
								Type:    "generic",
								Generic: "valueType",
							},
						},
					},
				},
				Out: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"{entries}": {
							Type:    "generic",
							Generic: "valueType",
						},
					},
				},
			},
		},
		DelegateDefs: map[string]*core.DelegateDef{},
		PropertyDefs: core.TypeDefMap{
			"key": {
				Type: "string",
			},
			"value": {
				Type: "string",
			},
			"entries": {
				Type: "stream",
				Stream: &core.TypeDef{
					Type: "string",
				},
			},
		},
	},
	opFunc: func(op *core.Operator) {
		in := op.Main().In()
		out := op.Main().Out()
		entries := []string{}
		for _, entry := range op.Property("entries").([]interface{}) {
			entries = append(entries, entry.(string))
		}
		keyStr := op.Property("key").(string)
		valueStr := op.Property("value").(string)
		for !op.CheckStop() {
			i := in.Pull()
			if core.IsMarker(i) {
				out.Push(i)
				continue
			}

			is := i.([]interface{})

			mapOut := make(map[string]interface{})
			for _, entry := range entries {
				for _, value := range is {
					valueMap := value.(map[string]interface{})
					key := valueMap[keyStr].(string)
					value := valueMap[valueStr]
					if key == entry {
						mapOut[entry] = value
					}
				}
				if _, ok := mapOut[entry]; !ok {
					mapOut[entry] = nil
				}
			}
			out.Push(mapOut)
		}
	},
}
