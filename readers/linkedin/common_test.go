package linkedin

import (
	"fmt"
	"testing"

	"gop.kg/src-d/domain@v6/models"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type linkedInSuite struct {
	session   *mgo.Session
	db        *mgo.Database
	compStore *models.CompanyStore
	infoStore *models.CompanyInfoStore
}

var _ = Suite(&linkedInSuite{})

func (s *linkedInSuite) SetUpSuite(c *C) {
	var err error

	s.session, err = mgo.Dial("127.0.0.1:27017")
	c.Assert(err, IsNil, Commentf("A local MongoDB instance is required for tests"))

	s.db = s.session.DB("unittest")
	s.compStore = models.NewCompanyStore(s.db)
	s.infoStore = models.NewCompanyInfoStore(s.db)

	foo1 := s.compStore.New("Foo Inc.", "foo", bson.NewObjectId())
	foo1.LinkedInCompanyIds = []int{1}
	foo1.AssociateCompanyIds = []int{2}
	_, err = s.compStore.Save(foo1)
	c.Assert(err, IsNil)

	for i := 1; i < 3; i++ {
		foo2 := s.infoStore.New(i)
		foo2.Name = fmt.Sprintf("Foo %d Inc.", i)
		_, err = s.infoStore.Save(foo2)
		c.Assert(err, IsNil)
	}
}

func (s *linkedInSuite) TearDownSuite(c *C) {
	err := s.db.DropDatabase()
	c.Assert(err, IsNil)

	s.session.Close()
}
