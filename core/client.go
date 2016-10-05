package core

import (
	"gop.kg/src-d/domain@v6/container"
	"gopkg.in/mgo.v2"
)

const (
	DatabaseName = "sources"
)

type Client struct {
	session mgo.Session
}

func NewClient() *Client {
	return &Client{*container.GetAnalysisMgoSession()}
}

func (c *Client) Collection(collection string) *mgo.Collection {
	return c.session.DB(DatabaseName).C(collection)
}

func (c *Client) DropDatabase() {
	c.session.DB(DatabaseName).DropDatabase()
}
