package store

import (
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/usace/goquery"
)

var domainTable = goquery.TableDataSet{
	Name:   "domain",
	Schema: "nsi",
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert":     `insert into domain (name, version, notes) values ($1, $2, $3) returning id`,
	},
	Fields: model.Domain{},
}

var fieldTable = goquery.TableDataSet{
	Name:   "field",
	Schema: "nsi",
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert":     `insert into field (name, type, description, is_domain) values ($1, $2, $3, $4) returning id`,
	},
	Fields: model.Field{},
}

var schemaTable = goquery.TableDataSet{
	Name:   "schema",
	Schema: "nsi",
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectId":   `select id from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert":     `insert into field (name, version, notes) values ($1, $2, $3) returning id`,
	},
	Fields: model.Schema{},
}

var datasetTable = goquery.TableDataSet{
	Name:   "dataset",
	Schema: "nsi",
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert":     `insert into field (name, version, notes) values ($1, $2, $3) returning id`,
	},
	Fields: model.Schema{},
}
