package linkedin

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"gop.kg/src-d/domain@v2.4/models"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type linkedInSuite struct {
	session *mgo.Session
	db      *mgo.Database
	store   *models.CompanyStore
}

var _ = Suite(&linkedInSuite{})

func (s *linkedInSuite) SetUpSuite(c *C) {
	var err error

	s.session, err = mgo.Dial("127.0.0.1:27017")
	c.Assert(err, IsNil, Commentf("A local MongoDB instance is required for tests"))

	s.db = s.session.DB("unittest")
	s.store = models.NewCompanyStore(s.db)

	foo := s.store.New("Foo Inc.", "foo", bson.NewObjectId())
	_, err = s.store.Save(foo)
	c.Assert(err, IsNil)
}

func (s *linkedInSuite) TearDownSuite(c *C) {
	err := s.db.DropDatabase()
	c.Assert(err, IsNil)

	s.session.Close()
}
