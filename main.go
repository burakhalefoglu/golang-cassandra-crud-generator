package main

import "golang-cass-crud-gen/helpers"

func main() {
	type ServerModel struct {
		Name     string
		Id       int
		Enabled  bool
		ClientId int64
	}

	helpers.CreateCassCrud("test_db", "name, id", &ServerModel{})
}
