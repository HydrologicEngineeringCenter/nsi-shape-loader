package store

import (
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/usace/goquery"
)

// TODO maybe this should be a config field
const dbSchema string = "nsi"

var accessTable = goquery.TableDataSet{
	Name:   "access",
	Schema: dbSchema,
	Statements: map[string]string{
		"selectId": `select id from access where dataset_id=$1 and value=$2`,
		"insert":   `insert into access (dataset_id, access_group, role, permission) values ($1, $2, $3, $4) returning id`,
	},
	Fields: model.Domain{},
}

var datasetTable = goquery.TableDataSet{
	Name:   "dataset",
	Schema: dbSchema,
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert": `insert into field (
            name,
            version,
            nsi_schema_id,
            table_name,
            shape,
            description,
            purpose,
            created_by,
            quality_id,
        ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`,
	},
	Fields: model.Dataset{},
}

var domainTable = goquery.TableDataSet{
	Name:   "domain",
	Schema: dbSchema,
	Statements: map[string]string{
		"selectId": `select id from domain where field_id=$1 and value=$2`,
		"insert":   `insert into domain (field_id, value) values ($1, $2) returning id`,
	},
	Fields: model.Domain{},
}

var fieldTable = goquery.TableDataSet{
	Name:   "field",
	Schema: dbSchema,
	Statements: map[string]string{
		"select":     `select id from field where name=$1`,
		"selectById": `select * from field where id=$1`,
		"insert":     `insert into field (name, type, description, is_domain) values ($1, $2, $3, $4) returning id`,
	},
	Fields: model.Field{},
}

var qualityTable = goquery.TableDataSet{
	Name:   "quality",
	Schema: dbSchema,
	Statements: map[string]string{
		"selectId": `select id from quality where value=$1`,
		"insert":   `insert into quality (value, description) values ($1, $2) returning id`,
	},
	Fields: model.Quality{},
}

var schemaFieldTable = goquery.TableDataSet{
	Name:   "schema_field",
	Schema: dbSchema,
	Statements: map[string]string{
		"select": `select id from schema_field`,
		"insert": `insert into schemafield (id, field_id) values ($1, $2) returning id`,
	},
	Fields: model.Field{},
}

var schemaTable = goquery.TableDataSet{
	Name:   "schema",
	Schema: dbSchema,
	Statements: map[string]string{
		"select":     `select * from domain where name=$1 and version=$2`,
		"selectId":   `select id from domain where name=$1 and version=$2`,
		"selectById": `select * from domain where id=$1`,
		"insert":     `insert into field (name, version, notes) values ($1, $2, $3) returning id`,
	},
	Fields: model.Schema{},
}
