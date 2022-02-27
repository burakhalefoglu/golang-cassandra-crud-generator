package helpers

import (
	"golang-cass-crud-gen/structs"
	"regexp"
	"strings"
)

func FormatStrFromUppercaseToLowercase(text string) string {
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	submatchall := re.FindAllString(text, -1)
	name := ""
	for _, element := range submatchall {
		if name == "" {
			name = element
			continue
		}
		name = name + "_" + element
	}
	return strings.ToLower(name)
}

func StructToFieldListWithType(model interface{}) string {

	fields := structs.Names(model)
	fieldList := ""
	for i, v := range fields {
		t := structs.Fields(model)[i].Kind()
		fieldList = fieldList + FormatStrFromUppercaseToLowercase(v) + " " + golangTypeToCassType(t.String()) + ", "
	}
	return fieldList
}

func StructToFieldListWithoutType(model interface{}) string {

	fields := structs.Names(model)
	fieldList := ""
	for i, v := range fields {
		if i != len(fields)-1 {
			fieldList = fieldList + FormatStrFromUppercaseToLowercase(v) + ", "
			continue
		}
		fieldList = fieldList + FormatStrFromUppercaseToLowercase(v)
	}
	return fieldList
}

func StructToFieldListWithoutTypeWithQuestionMark(model interface{}) string {

	fields := structs.Names(model)
	fieldList := ""
	for i, v := range fields {
		if i != len(fields)-1 {
			fieldList = fieldList + FormatStrFromUppercaseToLowercase(v) + "=?, "
			continue
		}
		fieldList = fieldList + FormatStrFromUppercaseToLowercase(v) + "=?"
	}
	return fieldList
}

func GetQuestionMarkByStructFieldCount(model interface{}) string {
	fields := structs.Names(model)
	questionMarks := ""

	for i := 0; i < len(fields); i++ {
		if i != len(fields)-1 {
			questionMarks = questionMarks + "?, "
			continue
		}
		questionMarks = questionMarks + "?"
	}
	return questionMarks
}

func GetListToStr(t []string) string {

	texts := ""
	for i, v := range t {
		if i != len(t)-1 {
			texts = texts + FormatStrFromUppercaseToLowercase(v) + ", "
			continue
		}
		texts = texts + FormatStrFromUppercaseToLowercase(v)
	}
	return texts
}

func PkToStrQuery(pk []string) string {
	texts := ""
	for i, v := range pk {
		if i != len(pk)-1 {
			texts = texts + FormatStrFromUppercaseToLowercase(v) + "= ? AND "
			continue
		}
		texts = texts + FormatStrFromUppercaseToLowercase(v) + " = ?"
	}
	return texts

}

func golangTypeToCassType(t string) string {
	fieldName := ""
	switch t {
	case "int64":
		fieldName = "bigint"
	case "bool":
		fieldName = "boolean"
	case "float32":
		fieldName = "float"
	case "float64":
		fieldName = "double"
	case "int32":
		fieldName = "int"
	case "int":
		fieldName = "int"
	case "int16":
		fieldName = "smallint"
	case "int8":
		fieldName = "tinyint"
	case "string":
		fieldName = "text"
	}
	return fieldName
}
