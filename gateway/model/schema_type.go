package model

type (

	// Type is the data structure for storing the parsed values of schema string
	Type map[string]Collection // key is database name
	// Collection is a data structure for storing fields of schema
	Collection map[string]Fields // key is collection name
	// Fields is a data structure for storing the type of field
	Fields map[string]*FieldType // key is field name

	// FieldType stores information about a particular column in table
	FieldType struct {
		FieldName           string `json:"fieldName"`
		IsFieldTypeRequired bool   `json:"isFieldTypeRequired"`
		IsList              bool   `json:"isList"`
		Kind                string `json:"kind"`
		// Directive           string
		NestedObject Fields `json:"nestedObject"`

		// For directives
		IsPrimary   bool             `json:"isPrimary"`
		IsIndex     bool             `json:"isIndex"`
		IsUnique    bool             `json:"isUnique"`
		IsCreatedAt bool             `json:"isCreatedAt"`
		IsUpdatedAt bool             `json:"isUpdatedAt"`
		IsLinked    bool             `json:"isLinked"`
		IsForeign   bool             `json:"isForeign"`
		IsDefault   bool             `json:"isDefault"`
		IndexInfo   *TableProperties `json:"indexInfo"`
		LinkedTable *TableProperties `json:"linkedTable"`
		JointTable  *TableProperties `json:"jointTable"`
		Default     interface{}      `json:"default"`
	}

	// TableProperties are properties of the table
	TableProperties struct {
		From, To               string
		Table, Field, OnDelete string
		DBType                 string
		Group, Sort            string
		Order                  int
		ConstraintName         string
	}
)

const (
	// TypeInteger is variable used for Variable of type Integer
	TypeInteger string = "Integer"
	// TypeString is variable used for Variable of type String
	TypeString string = "String"
	// TypeFloat is variable used for Variable of type Float
	TypeFloat string = "Float"
	// TypeBoolean is variable used for Variable of type Boolean
	TypeBoolean string = "Boolean"
	// TypeDateTime is variable used for Variable of type DateTime
	TypeDateTime string = "DateTime"
	// TypeID is variable used for Variable of type ID
	TypeID string = "ID"
	// TypeJSON is variable used for Variable of type Jsonb
	TypeJSON string = "JSON"
	// SQLTypeIDSize is variable used for specifing si	ze of sql type ID
	SQLTypeIDSize string = "50"
	// TypeObject is a string with value object
	TypeObject string = "Object"
	// TypeEnum is a variable type enum
	TypeEnum string = "Enum"
	// DirectiveUnique is used in schema module to add unique index
	DirectiveUnique string = "unique"
	// DirectiveIndex is used in schema module to add index
	DirectiveIndex string = "index"
	// DirectiveForeign is used in schema module to add foreign key
	DirectiveForeign string = "foreign"
	// DirectivePrimary is used in schema module to add primary key
	DirectivePrimary string = "primary"
	// DirectiveCreatedAt is used in schema module to specify the created location
	DirectiveCreatedAt string = "createdAt"
	// DirectiveUpdatedAt  is used in schema module to add Updated location
	DirectiveUpdatedAt string = "updatedAt"
	// DirectiveLink is used in schema module to add link
	DirectiveLink string = "link"
	// DirectiveDefault is used to add default key
	DirectiveDefault string = "default"

	// DefaultIndexName  string = ""

	// DefaultIndexSort specifies default order of sorting
	DefaultIndexSort string = "asc"
	// DefaultIndexOrder specifies default order of order
	DefaultIndexOrder int = 1
)

// InspectorFieldType is the type for storing sql inspection information
type InspectorFieldType struct {
	FieldName    string `db:"Field"`
	FieldType    string `db:"Type"`
	FieldNull    string `db:"Null"`
	FieldKey     string `db:"Key"`
	FieldDefault string `db:"Default"`
}

// ForeignKeysType is the type for storing  foreignkeys information of sql inspection
type ForeignKeysType struct {
	TableName      string `db:"TABLE_NAME"`
	ColumnName     string `db:"COLUMN_NAME"`
	ConstraintName string `db:"CONSTRAINT_NAME"`
	DeleteRule     string `db:"DELETE_RULE"`
	RefTableName   string `db:"REFERENCED_TABLE_NAME"`
	RefColumnName  string `db:"REFERENCED_COLUMN_NAME"`
}

// IndexType is the type use to indexkey information of sql inspection
type IndexType struct {
	TableName  string `db:"TABLE_NAME"`
	ColumnName string `db:"COLUMN_NAME"`
	IndexName  string `db:"INDEX_NAME"`
	Order      int    `db:"SEQ_IN_INDEX"`
	Sort       string `db:"SORT"`
	IsUnique   string `db:"IS_UNIQUE"`
}
