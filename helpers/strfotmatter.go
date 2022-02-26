package helpers

import (
	"golang-cass-crud-gen/structs"
	"regexp"
	"strings"
)

func formatStrFromUppercaseToLowercase(text string) string {
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

func structToFieldListWithType(model interface{}) string {

	fields := structs.Names(model)
	fieldList := ""
	for i, v := range fields {
		t := structs.Fields(model)[i].Kind()
		fieldList = fieldList + formatStrFromUppercaseToLowercase(v) + " " + golangTypeToCassType(t.String()) + ", "
	}
	return fieldList
}

func structToFieldListWithoutType(model interface{}) string {

	fields := structs.Names(model)
	fieldList := ""
	for i, v := range fields {
		if i != len(fields)-1 {
			fieldList = fieldList + formatStrFromUppercaseToLowercase(v) + ", "
			continue
		}
		fieldList = fieldList + formatStrFromUppercaseToLowercase(v)
	}
	return fieldList
}

func getQuestionMarkByStructFieldCount(model interface{}) string {
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
