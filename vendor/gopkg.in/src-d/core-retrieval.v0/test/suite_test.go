package test

import (
	"testing"

	"gopkg.in/src-d/core-retrieval.v0/model"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(SuiteSuite))
}

type SuiteSuite struct {
	Suite

	store *model.RepositoryStore
}

func (s *SuiteSuite) SetupTest() {
	s.Setup()

	s.store = model.NewRepositoryStore(s.DB)
}

func (s *SuiteSuite) TearDownTest() {
	s.TearDown()
}

func (s *SuiteSuite) TestSchemaChanges() {
	err := s.store.Insert(model.NewRepository())
	s.NoError(err)

	repo, err := s.store.FindOne(model.NewRepositoryQuery())
	s.NoError(err)
	s.NotNil(repo)
}
