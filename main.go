package main

import "golang-cass-crud-gen/creator"

func main() {
	type ServerModelDto struct {
		Name                     string
		Id                       int
		Enabled                  bool
		ClientId                 int64
		AverageDailySessionCount int16
	}

	creator.CreateCassCrud("test_db",
		[]string{
			"Id",
			"ClientId",
		}, &ServerModelDto{})
}
