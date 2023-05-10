package json2sql_test

import (
	"reflect"
	"testing"

	"github.com/willf/json2sql/pkg/json2sql"
)

func TestDefaultValueTypeToSQLType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			value: "hello",
			want:  "VARCHAR",
		},
		{
			name:  "integer",
			value: float64(42),
			want:  "INTEGER",
		},
		{
			name:  "float",
			value: 3.14,
			want:  "DOUBLE",
		},
		{
			name:  "boolean",
			value: true,
			want:  "BOOLEAN",
		},
		{
			name:  "map",
			value: map[string]interface{}{"name": "John", "age": 30},
			want:  "ROW",
		},
		{
			name:  "unknown",
			value: struct{}{},
			want:  "VARCHAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json2sql.DefaultValueTypeToSQLType(tt.value); got != tt.want {
				t.Errorf("DefaultValueTypeToSQLType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeSingleQuotes(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "no quotes",
			value: "hello world",
			want:  "hello world",
		},
		{
			name:  "single quote",
			value: "it's raining",
			want:  "it''s raining",
		},
		{
			name:  "multiple quotes",
			value: "I said 'hello'",
			want:  "I said ''hello''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json2sql.EscapeSingleQuotes(tt.value); got != tt.want {
				t.Errorf("EscapeSingleQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMap_ValueToSQLValue(t *testing.T) {
	typeMap := json2sql.TypeMap{
		"field1": "VARCHAR",
		"field2": "DOUBLE",
		"field3": "INTEGER",
		"field4": "BOOLEAN",
		"field5": "UNKNOWN",
	}

	tests := []struct {
		name  string
		field string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			field: "field1",
			value: "hello",
			want:  "'hello'",
		},
		{
			name:  "float",
			field: "field2",
			value: 3.14,
			want:  "3.140000",
		},
		{
			name:  "integer",
			field: "field3",
			value: 42.0,
			want:  "42",
		},
		{
			name:  "boolean",
			field: "field4",
			value: true,
			want:  "true",
		},
		{
			name:  "unknown",
			field: "field5",
			value: "unknown",
			want:  "'unknown'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeMap.ValueToSQLValue(tt.field, tt.value); got != tt.want {
				t.Errorf("ValueToSQLValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonObject_SortedKeys(t *testing.T) {
	tests := []struct {
		name string
		obj  json2sql.JsonObject
		want []string
	}{
		{
			name: "empty object",
			obj:  json2sql.JsonObject{},
			want: []string{},
		},
		{
			name: "single key",
			obj:  json2sql.JsonObject{"foo": "bar"},
			want: []string{"foo"},
		},
		{
			name: "multiple keys",
			obj:  json2sql.JsonObject{"foo": "bar", "baz": "qux"},
			want: []string{"baz", "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.SortedKeys()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMap_SortedKeys(t *testing.T) {
	typeMap := json2sql.TypeMap{
		"field1": "VARCHAR",
		"field2": "DOUBLE",
		"field3": "INTEGER",
		"field4": "BOOLEAN",
		"field5": "UNKNOWN",
	}

	tests := []struct {
		name string
		obj  json2sql.TypeMap
		want []string
	}{
		{
			name: "empty map",
			obj:  json2sql.TypeMap{},
			want: []string{},
		},
		{
			name: "single key",
			obj:  json2sql.TypeMap{"foo": "bar"},
			want: []string{"foo"},
		},
		{
			name: "multiple keys",
			obj:  typeMap,
			want: []string{"field1", "field2", "field3", "field4", "field5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.SortedKeys()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMap_CreateCreateStatement(t *testing.T) {
	typeMap := json2sql.TypeMap{
		"field2": "DOUBLE",
		"field1": "VARCHAR",
		"field3": "INTEGER",
		"field4": "BOOLEAN",
		"field5": "UNKNOWN",
	}

	tests := []struct {
		name    string
		table   string
		want    string
		wantErr bool
	}{

		{
			name:    "valid table name",
			table:   "mytable",
			want:    "CREATE TABLE IF NOT EXISTS mytable (field1 VARCHAR, field2 DOUBLE, field3 INTEGER, field4 BOOLEAN, field5 UNKNOWN);",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := typeMap.CreateCreateStatement(tt.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCreateStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateCreateStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonObject_CreateTypeMap(t *testing.T) {
	tests := []struct {
		name string
		obj  json2sql.JsonObject
		want json2sql.TypeMap
	}{
		{
			name: "single key",
			obj:  json2sql.JsonObject{"foo": "bar"},
			want: json2sql.TypeMap{"foo": "VARCHAR"},
		},
		{
			name: "multiple keys",
			obj:  json2sql.JsonObject{"foo": "bar", "baz": float64(42)},
			want: json2sql.TypeMap{"foo": "VARCHAR", "baz": "INTEGER"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.CreateTypeMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTypeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
