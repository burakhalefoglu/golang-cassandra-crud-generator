package helpers

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"golang-cass-crud-gen/structs"
)

func CreateCassCrud(keyspace string, pk string, model interface{}) {
	f := NewFilePathName("generated", "generated")
	f.ImportName("github.com/gocql/gocql", "gocql")
	//create model
	createInterface(f, model)
	createStruct(f, model)
	generateCassTableCreateQuery(model, f, keyspace, pk)
	createAddOperation(f, model, keyspace)
	createGetByIdOperation(f, model, keyspace)
	createFileOnDirectory(fmt.Sprintf("%#v", f), "Cass"+structs.Name(model)+"Dal.go", "generated")
}

func createInterface(f *File, model interface{}) {
	f.Type().Id("I"+structs.Name(model)+"Dal").Interface(
		Id("Add").Params(
			Id("m").Op("*").Id(structs.Name(model)),
		).Error(),
		Id("GetById").Params(
			Id("ClientId").Int64(),
			Id("ProjectId").Int64(),
		).Params(Id("m").Op("*").Id(structs.Name(model)), Id("e").Error()),
		Id("GetAll").Params().Params(Id("m").Op("*").Index().Id(structs.Name(model)), Id("e").Error()),
		Id("UpdateById").Params(
			Id("ClientId").Int64(),
			Id("ProjectId").Int64(),
			Id("m").Op("*").Id(structs.Name(model)),
		).Error(),
		Id("DeleteById").Params(
			Id("ClientId").Int64(),
			Id("ProjectId").Int64(),
		).Error(),
	)
}

func createStruct(f *File, model interface{}) {
	f.Type().Id("cass"+structs.Name(model)+"Dal").Struct(
		Id("Client").Op("*").Id("gocql.Session"),
		Id("Table").String(),
	)
}

func generateCassTableCreateQuery(model interface{}, f *File, keyspaceName string, primaryKeys string) {
	tableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s( %s PRIMARY KEY((%s)))",
		keyspaceName, formatStrFromUppercaseToLowercase(structs.Name(model)), structToFieldListWithType(model), primaryKeys)
	f.Var().Id("Create" + structs.Name(model) + "TableQuery").Op("=").Lit(tableQuery)
}

func generateCassCreateModelQuery(model interface{}, keyspaceName string) string {
	tableQuery := fmt.Sprintf("INSERT INTO %s.%s(%s)VALUES(%s)",
		keyspaceName, formatStrFromUppercaseToLowercase(structs.Name(model)), structToFieldListWithoutType(model), getQuestionMarkByStructFieldCount(model))
	return tableQuery
}

func generateCassGetByIdModelQuery(model interface{}, keyspaceName string) string {
	tableQuery := fmt.Sprintf("SELECT %s FROM %s.%s WHERE ClientId  = ? AND ProjectId = ? LIMIT 1 ",
		structToFieldListWithoutType(model), keyspaceName, formatStrFromUppercaseToLowercase(structs.Name(model)))
	return tableQuery
}

func createAddOperation(f *File, model interface{}, keyspaceName string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("Add").Params(
		Id("m").Op("*").Id(structs.Name(model)),
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call(getCreateCQueryCallParams(keyspaceName, model)...).Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}

func getCreateCQueryCallParams(keyspaceName string, model interface{}) []Code {
	var code []Code
	fields := structs.Names(model)
	for i, v := range fields {
		if i == 0 {
			code = append(code, Lit(generateCassCreateModelQuery(model, keyspaceName)))
		}
		code = append(code, Id("m.").Id(v))
	}
	return code
}

func createGetByIdOperation(f *File, model interface{}, keySpace string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("GetById").Params(
		Id("ClientId").Int64(),
		Id("ProjectId").Int64(),
	).Params(Id("data").Op("*").Id(structs.Name(model)), Id("e").Error()).Block(
		Id("m").Op(":=").Op("&").Id(structs.Name(model)).Id("{}"),
		If(
			Err().Op(":=").Id("c.Client.Query").Call(
				Lit(generateCassGetByIdModelQuery(model, keySpace)), Id("ClientId"), Id("ProjectId")).Id(".Scan").Call(
				getCreateRQueryCallParams(model)...),
			Err().Op("!=").Nil(),
		).Block(
			Return(Nil(), Err()),
		),
		Return(Id("m"), Nil()),
	)
}

func getCreateRQueryCallParams(model interface{}) []Code {
	var code []Code
	fields := structs.Names(model)
	for _, v := range fields {
		code = append(code, Id("m.").Id(v))
	}
	return code
}
func createGetAllOperation(f *File, model interface{}) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("Add").Params(
		Id("m").Op("*").Id(structs.Name(model)),
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call().Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}

func createUpdateByIdOperation(f *File, model interface{}) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("Add").Params(
		Id("m").Op("*").Id(structs.Name(model)),
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call().Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}

func createDeleteByIdOperation(f *File, model interface{}) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("Add").Params(
		Id("m").Op("*").Id(structs.Name(model)),
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call().Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}
