package creator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"golang-cass-crud-gen/helpers"
	"golang-cass-crud-gen/structs"
	"reflect"
)

func CreateCassCrud(keyspace string, pk []string, model interface{}) {
	f := NewFilePathName("generated", "generated")
	f.ImportName("github.com/gocql/gocql", "gocql")

	createModelSortedParams(f, model)
	createDalInterface(f, model, pk)
	createDalStruct(f, model)
	generateCassTableCreateQuery(model, f, keyspace, pk)
	generateAddOperation(f, model, keyspace)
	generateGetByIdOperation(f, model, keyspace, pk)
	generateGetAllOperation(f, model, keyspace)
	generateUpdateByIdOperation(f, model, keyspace, pk)
	generateDeleteByIdOperation(f, model, keyspace, pk)
	helpers.CreateFileOnDirectory(fmt.Sprintf("%#v", f), "Cass"+structs.Name(model)+"Dal.go", "generated")
}

func createModelSortedParams(f *File, model interface{}) {

	f.Type().Id(structs.Name(model)).Struct(
		getModelFiledCodes(model)...,
	)
}

func getModelFiledCodes(model interface{}) []Code {
	var codes []Code
	var code Code
	x := structs.Fields(model)
	kinds := make([]reflect.Kind, 0)
	for _, v := range x {
		kinds = append(kinds, v.Kind())
	}

	for i, v := range structs.Names(model) {
		switch kinds[i].String() {
		case "int64":
			code = Id(v).Int64()
		case "bool":
			code = Id(v).Bool()
		case "float32":
			code = Id(v).Float32()
		case "float64":
			code = Id(v).Float64()
		case "int32":
			code = Id(v).Int32()
		case "int":
			code = Id(v).Int64()
		case "int16":
			code = Id(v).Int16()
		case "int8":
			code = Id(v).Int8()
		case "string":
			code = Id(v).String()
		}
		codes = append(codes, code)
	}
	return codes
}

func createDalInterface(f *File, model interface{}, pk []string) {
	f.Type().Id("I"+structs.Name(model)+"Dal").Interface(
		Id("Add").Params(
			Id("m").Op("*").Id(structs.Name(model)),
		).Error(),
		Id("GetById").Params(
			getPKQueryWithType(pk)...,
		).Params(Id("m").Op("*").Id(structs.Name(model)), Id("e").Error()),
		Id("GetAll").Params().Params(Id("m").Op("*").Index().Id(structs.Name(model)), Id("e").Error()),
		Id("UpdateById").Params(
			getPKQueryWithTypeAndModel(pk, model)...,
		).Error(),
		Id("DeleteById").Params(
			getPKQueryWithType(pk)...,
		).Error(),
	)
}

func createDalStruct(f *File, model interface{}) {
	f.Type().Id("cass"+structs.Name(model)+"Dal").Struct(
		Id("Client").Op("*").Id("gocql.Session"),
		Id("Table").String(),
	)
}

func generateCassTableCreateQuery(model interface{}, f *File, keyspaceName string, pk []string) {
	tableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s( %s PRIMARY KEY((%s)))",
		keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)),
		helpers.StructToFieldListWithType(model), helpers.GetListToStr(pk))
	f.Var().Id("Create" + structs.Name(model) + "TableQuery").Op("=").Lit(tableQuery)
}

func generateCassCreateModelQuery(model interface{}, keyspaceName string) string {
	tableQuery := fmt.Sprintf("INSERT INTO %s.%s(%s)VALUES(%s)",
		keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)),
		helpers.StructToFieldListWithoutType(model), helpers.GetQuestionMarkByStructFieldCount(model))
	return tableQuery
}

func generateCassGetByIdModelQuery(model interface{}, keyspaceName string, pk []string) string {
	tableQuery := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s LIMIT 1 ",
		helpers.StructToFieldListWithoutType(model), keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)), helpers.PkToStrQuery(pk))
	return tableQuery
}

func generateCassGetAllModelQuery(model interface{}, keyspaceName string) string {
	tableQuery := fmt.Sprintf("SELECT %s FROM %s.%s",
		helpers.StructToFieldListWithoutType(model), keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)))
	return tableQuery
}

func generateCassUpdateByIdModelQuery(model interface{}, keyspaceName string, pk []string) string {
	tableQuery := fmt.Sprintf("UPDATE %s.%s SET %s WHERE %s",
		keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)),
		helpers.StructToFieldListWithoutTypeWithQuestionMark(model), helpers.PkToStrQuery(pk))
	return tableQuery
}

func generateCassDeleteByIdModelQuery(model interface{}, keyspaceName string, pk []string) string {
	tableQuery := fmt.Sprintf("DELETE FROM %s.%s WHERE %s",
		keyspaceName, helpers.FormatStrFromUppercaseToLowercase(structs.Name(model)), helpers.PkToStrQuery(pk))
	return tableQuery
}

func generateAddOperation(f *File, model interface{}, keyspaceName string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("Add").Params(
		Id("m").Op("*").Id(structs.Name(model)),
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call(getAddQueryCallParams(keyspaceName, model)...).Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}

func getAddQueryCallParams(keyspaceName string, model interface{}) []Code {
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

func generateGetByIdOperation(f *File, model interface{}, keySpace string, pk []string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("GetById").Params(
		getPKQueryWithType(pk)...,
	).Params(Id("data").Op("*").Id(structs.Name(model)), Id("e").Error()).Block(
		Id("m").Op(":=").Op("&").Id(structs.Name(model)).Id("{}"),
		If(
			Err().Op(":=").Id("c.Client.Query").Call(getQueryParams(model, keySpace, pk)...).Id(".Scan").Call(
				getScanQueryCallParams(model)...),
			Err().Op("!=").Nil(),
		).Block(
			Return(Nil(), Err()),
		),
		Return(Id("m"), Nil()),
	)
}

func getQueryParams(model interface{}, keySpace string, pk []string) []Code {
	var code []Code
	for i, v := range pk {
		if i == 0 {
			code = append(code, Lit(generateCassGetByIdModelQuery(model, keySpace, pk)))
		}
		code = append(code, Id(v))
	}
	return code
}

func getScanQueryCallParams(model interface{}) []Code {
	var code []Code
	fields := structs.Names(model)
	for _, v := range fields {
		code = append(code, Op("&").Id("m.").Id(v))
	}
	return code
}
func generateGetAllOperation(f *File, model interface{}, keyspace string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("GetAll").Params().Params(Id("data").Op("*").Index().Id(structs.Name(model)), Id("e").Error()).Block(
		Id("m").Op(":=").Id(structs.Name(model)).Id("{}"),
		Var().Id("models").Index().Id(structs.Name(model)),
		Id("iter").Op(":=").Id("c.Client.Query").Call(
			Lit(generateCassGetAllModelQuery(model, keyspace))).Id(".Iter()"),
		For(
			Id("iter").Id(".Scan").Call(getScanQueryCallParams(model)...),
		).Block(
			Id("models").Op("=").Append(Id("models"), Id("m")),
		),
		If(
			Err().Op(":=").Id("iter").Id(".Close").Call(),
			Err().Op("!=").Nil(),
		).Block(
			Return(Nil(), Err()),
		),
		Return(Op("&").Id("models"), Nil()),
	)
}

func generateUpdateByIdOperation(f *File, model interface{}, keyspace string, pk []string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("UpdateById").Params(
		getPKQueryWithTypeAndModel(pk, model)...,
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call(
				getUpdateQueryCallParams(keyspace, pk, model)...,
			).Id(".Exec").Call(),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}
func getUpdateQueryCallParams(keyspace string, pk []string, model interface{}) []Code {
	var code []Code
	fields := structs.Names(model)
	for i, v := range fields {
		if i == 0 {
			code = append(code, Lit(generateCassUpdateByIdModelQuery(model, keyspace, pk)))
		}
		code = append(code, Id("m.").Id(v))
	}
	for _, v := range pk {
		code = append(code, Id(v))
	}
	return code
}

func generateDeleteByIdOperation(f *File, model interface{}, keyspace string, pk []string) {
	f.Func().Params(
		Id("c").Op("*").Id("cass"+structs.Name(model)+"Dal"),
	).Id("DeleteById").Params(
		getPKQueryWithType(pk)...,
	).Error().Block(
		If(
			Err().Op(":=").Id("c.Client.Query").Call(getDeleteQueryCallParams(keyspace, pk, model)...).Id(".Exec()"),
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Return(Id("nil")),
	)
}
func getDeleteQueryCallParams(keyspace string, pk []string, model interface{}) []Code {
	var code []Code
	for i, v := range pk {
		if i == 0 {
			code = append(code, Lit(generateCassDeleteByIdModelQuery(model, keyspace, pk)))
		}
		code = append(code, Id(v))
	}
	return code
}

func getPKQueryWithType(pk []string) []Code {
	var code []Code
	for _, v := range pk {
		code = append(code, Id(v).Int64())
	}
	return code
}
func getPKQueryWithTypeAndModel(pk []string, model interface{}) []Code {
	var code []Code
	for i, v := range pk {
		code = append(code, Id(v).Int64())
		if i == len(pk)-1 {
			code = append(code, Id("m").Op("*").Id(structs.Name(model)))
		}
	}
	return code
}
