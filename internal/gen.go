package internal

import (
	"errors"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	ErrFileAlreadyExist = errors.New("file already exist")
)

type RepoTemplate struct {
	SchemaName      string
	CamelSchemaName string
	ModuleName      string
	Fields          []*Field
	FieldNameType   map[string]FieldType
}

func GenRepoFile(templateByte []byte, fields []*Field, workingDir string, modulePath, moduleName string, schemaName string) (err error) {

	tmpl, err := template.New("repo").Parse(string(templateByte))
	if err != nil {
		return err
	}

	exist, err := checkRepoFile(path.Join(workingDir, modulePath), schemaName)
	if err != nil || exist {
		return err
	}

	os.Mkdir(path.Join(workingDir, modulePath, schemaName), os.ModePerm)

	file, err := os.Create(path.Join(workingDir, modulePath, schemaName) + "/repo.go")
	if err != nil {
		return err
	}
	defer file.Close()

	camelSchemaName := strings.ToUpper(schemaName[:1]) + schemaName[1:]

	fieldNameType := make(map[string]FieldType, len(fields))
	for _, field := range fields {
		fieldNameType[field.Name] = field.Type
	}

	err = tmpl.Execute(file, RepoTemplate{
		ModuleName:      moduleName,
		SchemaName:      schemaName,
		CamelSchemaName: camelSchemaName,
		Fields:          fields,
		FieldNameType:   fieldNameType,
	})
	if err != nil {
		return err
	}

	return nil

}

func checkRepoFile(schemaPath, schemaName string) (bool, error) {
	if _, err := os.Stat(schemaPath + "/" + schemaName + "/repo.go"); err == nil {
		return true, ErrFileAlreadyExist
	}
	return false, nil

}

// type UsecaseTemplate struct {
// 	SchemaName      string
// 	CamelSchemaName string
// 	ModuleName      string
// 	Fields          []*Field
// 	FieldNameType   map[string]FieldType
// }

// func GenUsecaseFile(templatePath string, fields []*Field, workingDir string, modulePath, moduleName string, schemaName string) (err error) {
// 	tmpl, err := template.ParseFiles(path.Join(templatePath, "/usecase.tmpl"))
// 	if err != nil {
// 		return err
// 	}
// 	exist, err := checkUsecaseFile(path.Join(workingDir, modulePath), schemaName)
// 	if err != nil || exist {
// 		return err
// 	}

// 	os.Mkdir(path.Join(workingDir, modulePath, schemaName), os.ModePerm)

// 	file, err := os.Create(path.Join(workingDir, modulePath, schemaName) + "/usecase.go")
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	camelSchemaName := strings.ToUpper(schemaName[:1]) + schemaName[1:]

// 	fieldNameType := make(map[string]FieldType, len(fields))
// 	for _, field := range fields {
// 		fieldNameType[field.Name] = field.Type
// 	}

// 	err = tmpl.Execute(file, UsecaseTemplate{
// 		ModuleName:      moduleName,
// 		SchemaName:      schemaName,
// 		CamelSchemaName: camelSchemaName,
// 		Fields:          fields,
// 		FieldNameType:   fieldNameType,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

func checkUsecaseFile(schemaPath, schemaName string) (bool, error) {
	if _, err := os.Stat(schemaPath + "/" + schemaName + "/usecase.go"); err == nil {
		return true, ErrFileAlreadyExist
	}
	return false, nil

}
