package elem

import (
	"github.com/Bitspark/slang/pkg/core"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var databaseExecuteCfg = &builtinConfig{
	opDef: core.OperatorDef{
		Meta: core.OperatorMetaDef{
			Name: "DB execute",
		},
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"trigger": {
							Type: "trigger",
						},
						"{queryParams}": {
							Type: "primitive",
						},
					},
				},
				Out: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"rowsAffected": {
							Type: "number",
						},
						"lastInsertId": {
							Type: "number",
						},
					},
				},
			},
		},
		DelegateDefs: map[string]*core.DelegateDef{},
		PropertyDefs: map[string]*core.TypeDef{
			"query": {
				Type: "string",
			},
			"queryParams": {
				Type: "stream",
				Stream: &core.TypeDef{
					Type: "string",
				},
			},
			"driver": {
				Type: "string",
			},
			"url": {
				Type: "string",
			},
		},
	},
	opFunc: func(op *core.Operator) {
		query := op.Property("query").(string)

		driver := op.Property("driver").(string)
		url := op.Property("url").(string)

		params := []string{}
		for _, param := range op.Property("queryParams").([]interface{}) {
			params = append(params, param.(string))
		}

		db, err := sql.Open(driver, url)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err.Error())
		}

		stmt, err := db.Prepare(query)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()

		in := op.Main().In()
		out := op.Main().Out()
		for !op.CheckStop() {
			i := in.Pull()
			if core.IsMarker(i) {
				out.Push(i)
				continue
			}

			im := i.(map[string]interface{})

			args := []interface{}{}
			for _, param := range params {
				args = append(args, im[param])
			}
			result, err := stmt.Exec(args...)

			if err != nil {
				out.Push(nil)
				continue
			}

			if n, err := result.LastInsertId(); err == nil {
				out.Map("lastInsertId").Push(n)
			} else {
				out.Map("lastInsertId").Push(nil)
			}

			if n, err := result.RowsAffected(); err == nil {
				out.Map("rowsAffected").Push(n)
			} else {
				out.Map("rowsAffected").Push(nil)
			}
		}
	},
}
