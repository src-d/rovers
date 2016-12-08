package core

import (
	"gopkg.in/mgo.v2"
	"srcd.works/domain.v6/container"
)

type Client struct {
	dbName  string
	session mgo.Session
}

func NewClient(dbName string) *Client {
	return &Client{
		dbName:  dbName,
		session: *container.GetAnalysisMgoSession(),
	}
}

func (c *Client) Collection(collection string) *mgo.Collection {
	return c.session.DB(c.dbName).C(collection)
}

func (c *Client) DropDatabase() {
	c.session.DB(c.dbName).DropDatabase()
}
