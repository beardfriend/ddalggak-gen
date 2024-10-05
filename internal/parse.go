package internal

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	ErrEnitiyNotfound     = errors.New("entity not found")
	ErrEntityNoField      = errors.New("entity has no field")
	ErrModuleFileNotFound = errors.New("go.mod file not found")
	ErrModuleFileAbNormal = errors.New("go.mod file is abnormal")
)

// FieldType represents the type of the field.
type FieldType string

const (
	StringField FieldType = "string"
	UUIDField   FieldType = "uuid.UUID"
	IntField    FieldType = "int"
	Int64Field  FieldType = "int64"
	UintField   FieldType = "uint"
	Uint64Field FieldType = "uint64"
	BoolField   FieldType = "bool"
	FloatField  FieldType = "float"
	AnyField    FieldType = "any"
	CustomField FieldType = "custom"
)

// Field represents the field of the entity.
type Field struct {
	// Name of the field.
	Name string
	// Type of the field.
	Type FieldType
	// Flag for optional field.
	IsOptional bool
	// Flag indicating if the field is automatically generated when creating an entity.
	IsDefaultGeneratedColumn bool
	// Field is a relation id.
	IsRelationID bool
}

// ParseEntity parses the entity file and returns the fields of the entity.
func ParseEntity(workingDir string, schemaName string) (result []*Field, err error) {

	file, err := os.Open(workingDir + "/ent/" + strings.ToLower(schemaName) + ".go")
	defer file.Close()
	if err != nil {
		err = ErrEnitiyNotfound
		return
	}

	scanner := bufio.NewScanner(file)
	camelCase := strings.ToUpper(schemaName[:1]) + schemaName[1:]

	var isStruct bool

	exceptions := []string{"Edges", "selectValues", "config", "//", "loadedTypes"}
	defaultColumns := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"}
	fields := make([]*Field, 0)
	structs := make([][]*Field, 2)

	index := -1 // 0: schema, 1: relation
	for scanner.Scan() {

		if scanner.Text() == "type "+camelCase+" struct {" || scanner.Text() == "type "+camelCase+"Edges"+" struct {" {
			index++
			isStruct = true
			continue
		}

		if !isStruct {
			continue
		}
		if scanner.Text() == "}" {
			structs[index] = fields
			fields = make([]*Field, 0)
			isStruct = false
			continue
		}

		// exception
		isException := false
		for _, exception := range exceptions {
			if strings.Contains(scanner.Text(), exception) {
				isException = true
			}
		}

		if isException {
			continue
		}

		trimed := strings.TrimSpace(scanner.Text())
		spaceNameBetweenTypeIndex := strings.Index(trimed, " ")
		startedFromType := trimed[spaceNameBetweenTypeIndex+1:]
		spaceTypeBetweenTagIndex := strings.Index(startedFromType, " ")

		name := trimed[:spaceNameBetweenTypeIndex]
		typeName := startedFromType[:spaceTypeBetweenTagIndex]

		isDefaultGenreated := false
		for _, defaultColumn := range defaultColumns {
			if name == defaultColumn {
				isDefaultGenreated = true
			}
		}

		field := &Field{
			Name:                     name,
			IsDefaultGeneratedColumn: isDefaultGenreated,
		}

		switch typeName {
		case "string":
			field.Type = StringField
		case "int":
			field.Type = IntField
		case "uuid.UUID":
			if field.Name == "ID" {
				field.IsDefaultGeneratedColumn = false
			}
			field.Type = UUIDField
		case "int64":
			field.Type = Int64Field
		case "uint":
			field.Type = UintField
		case "uint64":
			field.Type = Uint64Field
		case "bool":
			field.Type = BoolField
		case "float":
			field.Type = FloatField
		case "interface{}":
			field.Type = AnyField
		default:
			field.Type = CustomField
		}

		if len(typeName) > 1 {
			if typeName[:1] == "*" {
				field.IsOptional = true
			}
		}

		fields = append(fields, field)
	}

	if len(structs) == 2 {
		for _, field := range structs[0] {
			for _, relation := range structs[1] {
				if field.Name == relation.Name+"ID" {
					field.IsRelationID = true
				}
			}
		}
	}

	if errA := scanner.Err(); errA != nil {
		err = errA
		return
	}
	if len(structs[0]) == 0 {
		err = ErrEntityNoField
		return
	}

	result = structs[0]
	return
}

// ParseProjectModuleName parses the go.mod file and returns the module name of the project.
func ParseProjectModuleName(workingDir string) (moduleName string, err error) {

	file, err := os.Open(workingDir + "/go.mod")
	defer file.Close()
	if err != nil {
		err = ErrModuleFileAbNormal
		return
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "module") {
			moduleName = strings.TrimSpace(strings.Split(scanner.Text(), " ")[1])
			return
		}
	}

	moduleName = strings.ReplaceAll(moduleName, "\"", "")

	if errA := scanner.Err(); errA != nil {
		err = errA
		return
	}

	if moduleName == "" {
		err = ErrModuleFileAbNormal
	}
	return

}
