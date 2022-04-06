package store

import (
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/usace/goquery"
)

// TODO maybe this should be a config field
const DbSchema string = "nsi"

var accessTable = goquery.TableDataSet{
	Name:   "access",
	Schema: DbSchema,
	Statements: map[string]string{
		"selectId": `select id from access where dataset_id=$1 and value=$2`,
		"insert":   `insert into access (dataset_id, access_group, role, permission) values ($1, $2, $3, $4) returning id`,
	},
	Fields: model.Domain{},
}

var datasetTable = goquery.TableDataSet{
	Name:   "dataset",
	Schema: DbSchema,
	Statements: map[string]string{
		"selectId":   `select id from dataset where name=$1 and version=$2 and purpose=$3 and quality_id=$4`,
		"select":     `select * from dataset where name=$1 and version=$2 and purpose=$3 and quality_id=$4`,
		"selectById": `select * from dataset where id=$1`,
		"insertNullShape": `insert into dataset (
            name,
            version,
            nsi_schema_id,
            table_name,
            shape,
            description,
            purpose,
            created_by,
            quality_id
        ) values ($1, $2, $3, $4, ST_Envelope('POLYGON((0 0, 0 0, 0 0, 0 0))'::geometry), $5, $6, $7, $8) returning id`,
		"updateBBox": `update dataset set shape=(select ST_Envelope(ST_Collect(shape)) from {table_name}) where id=$1`,
	},
}

var domainTable = goquery.TableDataSet{
	Name:   "domain",
	Schema: DbSchema,
	Statements: map[string]string{
		"selectId": `select id from domain where field_id=$1 and value=$2`,
		"insert":   `insert into domain (field_id, value) values ($1, $2) returning id`,
	},
	Fields: model.Domain{},
}

var fieldTable = goquery.TableDataSet{
	Name:   "field",
	Schema: DbSchema,
	Statements: map[string]string{
		"select":     `select id from field where name=$1`,
		"selectById": `select * from field where id=$1`,
		"insert":     `insert into field (name, type, description, is_domain) values ($1, $2, $3, $4) returning id`,
	},
	Fields: model.Field{},
}

var qualityTable = goquery.TableDataSet{
	Name:   "quality",
	Schema: DbSchema,
	Statements: map[string]string{
		"selectId": `select id from quality where value=$1`,
		"insert":   `insert into quality (value, description) values ($1, $2) returning id`,
	},
	Fields: model.Quality{},
}

var schemaFieldTable = goquery.TableDataSet{
	Name:   "schema_field",
	Schema: DbSchema,
	Statements: map[string]string{
		"selectId": `select id from schema_field where id=$1 and field_id=$2`,
		"insert":   `insert into schema_field (id, field_id, is_private) values ($1, $2, $3) returning id`,
	},
	Fields: model.Field{},
}

var schemaTable = goquery.TableDataSet{
	Name:   "schema",
	Schema: DbSchema,
	Statements: map[string]string{
		"select":     `select * from nsi_schema where name=$1 and version=$2`,
		"selectId":   `select id from nsi_schema where name=$1 and version=$2`,
		"selectById": `select * from nsi_schema where id=$1`,
		"insert":     `insert into nsi_schema (name, version, notes) values ($1, $2, $3) returning id`,
	},
	Fields: model.Schema{},
}
