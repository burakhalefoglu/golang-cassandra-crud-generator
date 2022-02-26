package generated

import "github.com/gocql/gocql"

type IServerModelDal interface {
	Add(m *ServerModel) error
	GetById(ClientId int64, ProjectId int64) (m *ServerModel, e error)
	GetAll() (m *[]ServerModel, e error)
	UpdateById(ClientId int64, ProjectId int64, m *ServerModel) error
	DeleteById(ClientId int64, ProjectId int64) error
}
type cassServerModelDal struct {
	Client *gocql.Session
	Table  string
}

var CreateServerModelTableQuery = "CREATE TABLE IF NOT EXISTS test_db.server_model( name text, id int, enabled boolean, client_id bigint,  PRIMARY KEY((name, id)))"

func (c *cassServerModelDal) Add(m *ServerModel) error {
	if err := c.Client.Query("INSERT INTO test_db.server_model(name, id, enabled, client_id)VALUES(?, ?, ?, ?)", m.Name, m.Id, m.Enabled, m.ClientId).Exec(); err != nil {
		return err
	}
	return nil
}
func (c *cassServerModelDal) GetById(ClientId int64, ProjectId int64) (data *ServerModel, e error) {
	m := &ServerModel{}
	if err := c.Client.Query("SELECT name, id, enabled, client_id FROM test_db.server_model WHERE ClientId  = ? AND ProjectId = ? LIMIT 1 ", ClientId, ProjectId).Scan(m.Name, m.Id, m.Enabled, m.ClientId); err != nil {
		return nil, err
	}
	return m, nil
}
