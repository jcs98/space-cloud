package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	goqu "github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *SQL) generator(find map[string]interface{}) (goqu.Expression, []string) {
	var regxarr []string
	array := []goqu.Expression{}
	for k, v := range find {
		if k == "$or" {
			orArray := v.([]interface{})
			orFinalArray := []goqu.Expression{}
			for _, item := range orArray {
				exp, a := s.generator(item.(map[string]interface{}))
				orFinalArray = append(orFinalArray, exp)
				regxarr = append(regxarr, a...)
			}

			array = append(array, goqu.Or(orFinalArray...))
			continue
		}
		val, isObj := v.(map[string]interface{})
		if isObj {
			for k2, v2 := range val {
				switch k2 {
				case "$regex":
					switch s.dbType {
					case "postgres":
						regxarr = append(regxarr, fmt.Sprintf("%s = $", k))
					case "mysql":
						regxarr = append(regxarr, fmt.Sprintf("%s = ?", k))
					}
					array = append(array, goqu.I(k).Eq(v2))
				case "$eq":
					array = append(array, goqu.I(k).Eq(v2))
				case "$ne":
					array = append(array, goqu.I(k).Neq(v2))
				case "$contains":
					data, err := json.Marshal(v2)
					if err != nil {
						logrus.Errorf("error marshalling data $contains data (%s)", err.Error())
						break
					}
					switch s.dbType {
					case string(utils.MySQL):
						array = append(array, goqu.L(fmt.Sprintf("json_contains(%s,?)", k), string(data)))
					case string(utils.Postgres):
						array = append(array, goqu.L(fmt.Sprintf("%s @> ?", k), string(data)))
					default:
						logrus.Errorf("_contains not supported for database (%s)", s.dbType)
					}
				case "$gt":
					array = append(array, goqu.I(k).Gt(v2))

				case "$gte":
					array = append(array, goqu.I(k).Gte(v2))

				case "$lt":
					array = append(array, goqu.I(k).Lt(v2))

				case "$lte":
					array = append(array, goqu.I(k).Lte(v2))

				case "$in":
					array = append(array, goqu.I(k).In(v2))

				case "$nin":
					array = append(array, goqu.I(k).NotIn(v2))
				}
			}
		} else {
			array = append(array, goqu.I(k).Eq(v))
		}
	}
	return goqu.And(array...), regxarr
}

func (s *SQL) generateWhereClause(q *goqu.SelectDataset, find map[string]interface{}) (query *goqu.SelectDataset, arr []string) {
	query = q
	if len(find) == 0 {
		return
	}
	exp, arr := s.generator(find)
	query = query.Where(exp)
	return query, arr
}

func generateRecord(temp interface{}) (goqu.Record, error) {
	insertObj, ok := temp.(map[string]interface{})
	if !ok {
		return nil, errors.New("incorrect insert object provided")
	}

	record := make(goqu.Record, len(insertObj))
	for k, v := range insertObj {
		record[k] = v
	}
	return record, nil
}

func (s *SQL) getDBName(col string) string {
	switch utils.DBType(s.dbType) {
	case utils.Postgres, utils.SQLServer:
		return fmt.Sprintf("%s.%s", s.name, col)
	}
	return col
}

func (s *SQL) generateQuerySQLServer(query string) string {
	return strings.Replace(query, "$", "@p", -1)
}

func mysqlTypeCheck(dbType utils.DBType, types []*sql.ColumnType, mapping map[string]interface{}) {
	var err error
	for _, colType := range types {
		typeName := colType.DatabaseTypeName()
		switch v := mapping[colType.Name()].(type) {
		case []byte:
			switch typeName {
			case "VARCHAR", "TEXT", "JSON", "JSONB":
				val, ok := mapping[colType.Name()].([]byte)
				if ok {
					mapping[colType.Name()] = string(val)
				}
			case "TINYINT":
				mapping[colType.Name()], err = strconv.ParseBool(string(v))
				if err != nil {
					log.Println("Error:", err)
				}
			case "BIGINT", "INT", "SMALLINT":
				mapping[colType.Name()], err = strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					log.Println("Error:", err)
				}
			case "DECIMAL", "NUMERIC", "FLOAT":
				mapping[colType.Name()], err = strconv.ParseFloat(string(v), 64)
				if err != nil {
					log.Println("Error:", err)
				}
			case "DATE", "DATETIME":
				if dbType == utils.MySQL {
					d, _ := time.Parse("2006-01-02 15:04:05", string(v))
					mapping[colType.Name()] = d.Format(time.RFC3339)
					continue
				}
				mapping[colType.Name()] = string(v)
			}
		case int64:
			if typeName == "TINYINT" {
				// this case occurs for mysql database with column type tinyint during the upsert operation
				if v == int64(1) {
					mapping[colType.Name()] = true
				} else {
					mapping[colType.Name()] = false
				}
			}
		case time.Time:
			mapping[colType.Name()] = v.UTC().Format(time.RFC3339)

		case primitive.DateTime:
			mapping[colType.Name()] = v.Time().UTC().Format(time.RFC3339)
		}
	}
}