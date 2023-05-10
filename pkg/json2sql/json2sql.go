package json2sql

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type JsonObject map[string]interface{}

type TypeMap map[string]string

func DefaultValueTypeToSQLType(value interface{}) string {
	switch value.(type) {
	case string:
		return "VARCHAR"
	case float64:
		_, frac := math.Modf(value.(float64))
		if frac == 0.0 {
			return "INTEGER"
		}
		return "DOUBLE"
	case bool:
		return "BOOLEAN"
	case map[string]interface{}:
		return "ROW"
	default:
		return "VARCHAR"
	}
}

func EscapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func (o TypeMap) ValueToSQLValue(field string, value interface{}) string {
	switch o[field] {
	case "VARCHAR":
		return fmt.Sprintf("'%s'", EscapeSingleQuotes(value.(string)))
	case "DOUBLE":
		return fmt.Sprintf("%f", value)
	case "INTEGER":
		return fmt.Sprintf("%d", int(value.(float64)))
	case "BOOLEAN":
		return fmt.Sprintf("%t", value)
	default:
		return fmt.Sprintf("'%s'", value)
	}
}

func (o JsonObject) SortedKeys() []string {
	keys := make([]string, 0)
	for key, _ := range o {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (o TypeMap) SortedKeys() []string {
	keys := make([]string, 0)
	for key, _ := range o {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (o TypeMap) CreateCreateStatement(table string) (statement string, err error) {
	// SQL DDL to create a table from a JSON object
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", table))

	declarations := make([]string, 0)
	for _, key := range o.SortedKeys() {
		declarations = append(declarations, fmt.Sprintf("%s %s", key, o[key]))
	}

	sql.WriteString(strings.Join(declarations, ", "))
	sql.WriteString(");")
	return sql.String(), nil
}

func (o JsonObject) CreateTypeMap() (typeMap TypeMap) {
	typeMap = make(TypeMap)
	for _, key := range o.SortedKeys() {
		value := o[key]
		valueType := DefaultValueTypeToSQLType(value)
		typeMap[key] = valueType
	}
	return
}

func (o JsonObject) CreateInsertStatementHeader(table string) (statement string, err error) {
	// INSERT INTO hive.scratch_platform_health.npm_spam_campaign (name, y, description, readme, keywords, day) VALUES
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("INSERT INTO %s (", table))

	sql.WriteString(strings.Join(o.SortedKeys(), ", "))
	sql.WriteString(") VALUES")
	return sql.String(), nil
}

func (o JsonObject) CreateValueStatement(typeMap TypeMap) (statement string, err error) {

	var sql strings.Builder
	sql.WriteString("(")

	values := make([]string, 0)
	for _, key := range o.SortedKeys() {
		value := o[key]
		values = append(values, typeMap.ValueToSQLValue(key, value))
	}
	sql.WriteString(strings.Join(values, ", "))
	sql.WriteString(")")
	return sql.String(), nil
}

func MainLoop(table string, create bool) {
	var typeMap TypeMap
	// Create a scanner to read from stdin
	onFirstLine := true
	// Create a scanner to read the file line by line
	reader := bufio.NewReader(os.Stdin)

	for {
		// Parse the JSON object from the line
		var object JsonObject
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		err = json.Unmarshal(line, &object)
		if err != nil {
			panic(err)
		}

		// Create the SQL create statement if on the first line
		// and generate the SQL insert statement
		if onFirstLine {
			onFirstLine = false
			typeMap = object.CreateTypeMap()
			if create {
				sql, err := typeMap.CreateCreateStatement(table)
				if err != nil {
					panic(err)
				}
				fmt.Println(sql)
			}

			sql, err := object.CreateInsertStatementHeader(table)
			if err != nil {
				panic(err)
			}

			// Output the SQL insert statement to the console
			fmt.Println(sql)
		}
		// now, get a values tuple
		sql, err := object.CreateValueStatement(typeMap)
		if err != nil {
			panic(err)
		}
		fmt.Print(sql)

		_, err = reader.Peek(1)
		if err != io.EOF {
			fmt.Println(",")
		} else {
			fmt.Println("")
		}
	}

	fmt.Println(";")

}
