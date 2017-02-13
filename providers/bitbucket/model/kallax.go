// IMPORTANT! This is auto generated code by https://github.com/src-d/go-kallax
// Please, do not touch the code below, and if you do, do it under your own
// risk. Take into account that all the code you write here will be completely
// erased from earth the next time you generate the kallax models.
package model

import (
	"database/sql"
	"fmt"

	"github.com/src-d/go-kallax"
	"github.com/src-d/go-kallax/types"
)

var _ types.SQLType
var _ fmt.Formatter

// NewRepository returns a new instance of Repository.
func NewRepository() (record *Repository) {
	record = &Repository{}
	if record != nil {
		record.SetID(kallax.NewID())
	}
	return
}

func (r *Repository) ColumnAddress(col string) (interface{}, error) {
	switch col {
	case "id":
		return &r.Model.ID, nil
	case "created_at":
		return &r.Timestamps.CreatedAt, nil
	case "updated_at":
		return &r.Timestamps.UpdatedAt, nil
	case "next":
		return &r.Next, nil
	case "scm":
		return &r.Scm, nil
	case "website":
		return &r.Website, nil
	case "name":
		return &r.Name, nil
	case "links":
		return types.JSON(&r.Links), nil
	case "fork_policy":
		return &r.ForkPolicy, nil
	case "uuid":
		return &r.UUID, nil
	case "language":
		return &r.Language, nil
	case "created_on":
		return &r.CreatedOn, nil
	case "parent":
		if r.Parent == nil {
			r.Parent = new(Parent)
		}
		return types.JSON(r.Parent), nil
	case "full_name":
		return &r.FullName, nil
	case "has_issues":
		return &r.HasIssues, nil
	case "owner":
		return types.JSON(&r.Owner), nil
	case "updated_on":
		return &r.UpdatedOn, nil
	case "size":
		return &r.Size, nil
	case "type":
		return &r.Type, nil
	case "slug":
		return &r.Slug, nil
	case "is_private":
		return &r.IsPrivate, nil
	case "description":
		return &r.Description, nil

	default:
		return nil, fmt.Errorf("kallax: invalid column in Repository: %s", col)
	}
}

func (r *Repository) Value(col string) (interface{}, error) {
	switch col {
	case "id":
		return r.Model.ID, nil
	case "created_at":
		return r.Timestamps.CreatedAt, nil
	case "updated_at":
		return r.Timestamps.UpdatedAt, nil
	case "next":
		return r.Next, nil
	case "scm":
		return r.Scm, nil
	case "website":
		return r.Website, nil
	case "name":
		return r.Name, nil
	case "links":
		return types.JSON(r.Links), nil
	case "fork_policy":
		return r.ForkPolicy, nil
	case "uuid":
		return r.UUID, nil
	case "language":
		return r.Language, nil
	case "created_on":
		return r.CreatedOn, nil
	case "parent":
		return types.JSON(r.Parent), nil
	case "full_name":
		return r.FullName, nil
	case "has_issues":
		return r.HasIssues, nil
	case "owner":
		return types.JSON(r.Owner), nil
	case "updated_on":
		return r.UpdatedOn, nil
	case "size":
		return r.Size, nil
	case "type":
		return r.Type, nil
	case "slug":
		return r.Slug, nil
	case "is_private":
		return r.IsPrivate, nil
	case "description":
		return r.Description, nil

	default:
		return nil, fmt.Errorf("kallax: invalid column in Repository: %s", col)
	}
}

func (r *Repository) NewRelationshipRecord(field string) (kallax.Record, error) {
	return nil, fmt.Errorf("kallax: model Repository has no relationships")
}

func (r *Repository) SetRelationship(field string, rel interface{}) error {
	return fmt.Errorf("kallax: model Repository has no relationships")
}

// RepositoryStore is the entity to access the records of the type Repository
// in the database.
type RepositoryStore struct {
	*kallax.Store
}

// NewRepositoryStore creates a new instance of RepositoryStore
// using a SQL database.
func NewRepositoryStore(db *sql.DB) *RepositoryStore {
	return &RepositoryStore{kallax.NewStore(db)}
}

// Insert inserts a Repository in the database. A non-persisted object is
// required for this operation.
func (s *RepositoryStore) Insert(record *Repository) error {

	if err := record.BeforeSave(); err != nil {
		return err
	}

	return s.Store.Insert(Schema.Repository.BaseSchema, record)

}

// Update updates the given record on the database. If the columns are given,
// only these columns will be updated. Otherwise all of them will be.
// Be very careful with this, as you will have a potentially different object
// in memory but not on the database.
// Only writable records can be updated. Writable objects are those that have
// been just inserted or retrieved using a query with no custom select fields.
func (s *RepositoryStore) Update(record *Repository, cols ...kallax.SchemaField) (updated int64, err error) {

	if err := record.BeforeSave(); err != nil {
		return 0, err
	}

	return s.Store.Update(Schema.Repository.BaseSchema, record, cols...)

}

// Save inserts the object if the record is not persisted, otherwise it updates
// it. Same rules of Update and Insert apply depending on the case.
func (s *RepositoryStore) Save(record *Repository) (updated bool, err error) {
	if !record.IsPersisted() {
		return false, s.Insert(record)
	}

	rowsUpdated, err := s.Update(record)
	if err != nil {
		return false, err
	}

	return rowsUpdated > 0, nil
}

// Delete removes the given record from the database.
func (s *RepositoryStore) Delete(record *Repository) error {

	return s.Store.Delete(Schema.Repository.BaseSchema, record)

}

// Find returns the set of results for the given query.
func (s *RepositoryStore) Find(q *RepositoryQuery) (*RepositoryResultSet, error) {
	rs, err := s.Store.Find(q)
	if err != nil {
		return nil, err
	}

	return NewRepositoryResultSet(rs), nil
}

// MustFind returns the set of results for the given query, but panics if there
// is any error.
func (s *RepositoryStore) MustFind(q *RepositoryQuery) *RepositoryResultSet {
	return NewRepositoryResultSet(s.Store.MustFind(q))
}

// Count returns the number of rows that would be retrieved with the given
// query.
func (s *RepositoryStore) Count(q *RepositoryQuery) (int64, error) {
	return s.Store.Count(q)
}

// MustCount returns the number of rows that would be retrieved with the given
// query, but panics if there is an error.
func (s *RepositoryStore) MustCount(q *RepositoryQuery) int64 {
	return s.Store.MustCount(q)
}

// FindOne returns the first row returned by the given query.
// `ErrNotFound` is returned if there are no results.
func (s *RepositoryStore) FindOne(q *RepositoryQuery) (*Repository, error) {
	q.Limit(1)
	q.Offset(0)
	rs, err := s.Find(q)
	if err != nil {
		return nil, err
	}

	if !rs.Next() {
		return nil, kallax.ErrNotFound
	}

	record, err := rs.Get()
	if err != nil {
		return nil, err
	}

	if err := rs.Close(); err != nil {
		return nil, err
	}

	return record, nil
}

// MustFindOne returns the first row retrieved by the given query. It panics
// if there is an error or if there are no rows.
func (s *RepositoryStore) MustFindOne(q *RepositoryQuery) *Repository {
	record, err := s.FindOne(q)
	if err != nil {
		panic(err)
	}
	return record
}

// Reload refreshes the Repository with the data in the database and
// makes it writable.
func (s *RepositoryStore) Reload(record *Repository) error {
	return s.Store.Reload(Schema.Repository.BaseSchema, record)
}

// Transaction executes the given callback in a transaction and rollbacks if
// an error is returned.
// The transaction is only open in the store passed as a parameter to the
// callback.
func (s *RepositoryStore) Transaction(callback func(*RepositoryStore) error) error {
	if callback == nil {
		return kallax.ErrInvalidTxCallback
	}

	return s.Store.Transaction(func(store *kallax.Store) error {
		return callback(&RepositoryStore{store})
	})
}

// RepositoryQuery is the object used to create queries for the Repository
// entity.
type RepositoryQuery struct {
	*kallax.BaseQuery
}

// NewRepositoryQuery returns a new instance of RepositoryQuery.
func NewRepositoryQuery() *RepositoryQuery {
	return &RepositoryQuery{
		BaseQuery: kallax.NewBaseQuery(Schema.Repository.BaseSchema),
	}
}

// Select adds columns to select in the query.
func (q *RepositoryQuery) Select(columns ...kallax.SchemaField) *RepositoryQuery {
	if len(columns) == 0 {
		return q
	}
	q.BaseQuery.Select(columns...)
	return q
}

// SelectNot excludes columns from being selected in the query.
func (q *RepositoryQuery) SelectNot(columns ...kallax.SchemaField) *RepositoryQuery {
	q.BaseQuery.SelectNot(columns...)
	return q
}

// Copy returns a new identical copy of the query. Remember queries are mutable
// so make a copy any time you need to reuse them.
func (q *RepositoryQuery) Copy() *RepositoryQuery {
	return &RepositoryQuery{
		BaseQuery: q.BaseQuery.Copy(),
	}
}

// Order adds order clauses to the query for the given columns.
func (q *RepositoryQuery) Order(cols ...kallax.ColumnOrder) *RepositoryQuery {
	q.BaseQuery.Order(cols...)
	return q
}

// BatchSize sets the number of items to fetch per batch when there are 1:N
// relationships selected in the query.
func (q *RepositoryQuery) BatchSize(size uint64) *RepositoryQuery {
	q.BaseQuery.BatchSize(size)
	return q
}

// Limit sets the max number of items to retrieve.
func (q *RepositoryQuery) Limit(n uint64) *RepositoryQuery {
	q.BaseQuery.Limit(n)
	return q
}

// Offset sets the number of items to skip from the result set of items.
func (q *RepositoryQuery) Offset(n uint64) *RepositoryQuery {
	q.BaseQuery.Offset(n)
	return q
}

// Where adds a condition to the query. All conditions added are concatenated
// using a logical AND.
func (q *RepositoryQuery) Where(cond kallax.Condition) *RepositoryQuery {
	q.BaseQuery.Where(cond)
	return q
}

// RepositoryResultSet is the set of results returned by a query to the
// database.
type RepositoryResultSet struct {
	ResultSet kallax.ResultSet
	last      *Repository
	lastErr   error
}

// NewRepositoryResultSet creates a new result set for rows of the type
// Repository.
func NewRepositoryResultSet(rs kallax.ResultSet) *RepositoryResultSet {
	return &RepositoryResultSet{ResultSet: rs}
}

// Next fetches the next item in the result set and returns true if there is
// a next item.
// The result set is closed automatically when there are no more items.
func (rs *RepositoryResultSet) Next() bool {
	if !rs.ResultSet.Next() {
		rs.lastErr = rs.ResultSet.Close()
		rs.last = nil
		return false
	}

	var record kallax.Record
	record, rs.lastErr = rs.ResultSet.Get(Schema.Repository.BaseSchema)
	if rs.lastErr != nil {
		rs.last = nil
	} else {
		var ok bool
		rs.last, ok = record.(*Repository)
		if !ok {
			rs.lastErr = fmt.Errorf("kallax: unable to convert record to *Repository")
			rs.last = nil
		}
	}

	return true
}

// Get retrieves the last fetched item from the result set and the last error.
func (rs *RepositoryResultSet) Get() (*Repository, error) {
	return rs.last, rs.lastErr
}

// ForEach iterates over the complete result set passing every record found to
// the given callback. It is possible to stop the iteration by returning
// `kallax.ErrStop` in the callback.
// Result set is always closed at the end.
func (rs *RepositoryResultSet) ForEach(fn func(*Repository) error) error {
	for rs.Next() {
		record, err := rs.Get()
		if err != nil {
			return err
		}

		if err := fn(record); err != nil {
			if err == kallax.ErrStop {
				return rs.Close()
			}

			return err
		}
	}
	return nil
}

// All returns all records on the result set and closes the result set.
func (rs *RepositoryResultSet) All() ([]*Repository, error) {
	var result []*Repository
	for rs.Next() {
		record, err := rs.Get()
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// One returns the first record on the result set and closes the result set.
func (rs *RepositoryResultSet) One() (*Repository, error) {
	if !rs.Next() {
		return nil, kallax.ErrNotFound
	}

	record, err := rs.Get()
	if err != nil {
		return nil, err
	}

	if err := rs.Close(); err != nil {
		return nil, err
	}

	return record, nil
}

// Err returns the last error occurred.
func (rs *RepositoryResultSet) Err() error {
	return rs.lastErr
}

// Close closes the result set.
func (rs *RepositoryResultSet) Close() error {
	return rs.ResultSet.Close()
}

type schema struct {
	Repository *schemaRepository
}

type schemaRepository struct {
	*kallax.BaseSchema
	ID          kallax.SchemaField
	CreatedAt   kallax.SchemaField
	UpdatedAt   kallax.SchemaField
	Next        kallax.SchemaField
	Scm         kallax.SchemaField
	Website     kallax.SchemaField
	Name        kallax.SchemaField
	Links       *schemaRepositoryLinks
	ForkPolicy  kallax.SchemaField
	UUID        kallax.SchemaField
	Language    kallax.SchemaField
	CreatedOn   kallax.SchemaField
	Parent      *schemaRepositoryParent
	FullName    kallax.SchemaField
	HasIssues   kallax.SchemaField
	Owner       *schemaRepositoryOwner
	UpdatedOn   kallax.SchemaField
	Size        kallax.SchemaField
	Type        kallax.SchemaField
	Slug        kallax.SchemaField
	IsPrivate   kallax.SchemaField
	Description kallax.SchemaField
}

type schemaRepositoryLinks struct {
	*kallax.BaseSchemaField
	Clone *schemaRepositoryLinksClone
}

type schemaRepositoryLinksClone struct {
	*kallax.JSONSchemaArray
	Href kallax.SchemaField
	Name kallax.SchemaField
}

func (s *schemaRepositoryLinksClone) At(n int) *schemaRepositoryLinksClone {
	return &schemaRepositoryLinksClone{
		JSONSchemaArray: kallax.NewJSONSchemaArray("links", "clone"),
		Href:            kallax.NewJSONSchemaKey(kallax.JSONText, "links", "clone", fmt.Sprint(n), "href"),
		Name:            kallax.NewJSONSchemaKey(kallax.JSONText, "links", "clone", fmt.Sprint(n), "name"),
	}
}

type schemaRepositoryOwner struct {
	*kallax.BaseSchemaField
	Username    kallax.SchemaField
	DisplayName kallax.SchemaField
	Type        kallax.SchemaField
	UUID        kallax.SchemaField
}

type schemaRepositoryParent struct {
	*kallax.BaseSchemaField
	Type     kallax.SchemaField
	Name     kallax.SchemaField
	FullName kallax.SchemaField
	UUID     kallax.SchemaField
}

var Schema = &schema{
	Repository: &schemaRepository{
		BaseSchema: kallax.NewBaseSchema(
			"bitbucket",
			"__repository",
			kallax.NewSchemaField("id"),
			kallax.ForeignKeys{},
			func() kallax.Record {
				return new(Repository)
			},
			kallax.NewSchemaField("id"),
			kallax.NewSchemaField("created_at"),
			kallax.NewSchemaField("updated_at"),
			kallax.NewSchemaField("next"),
			kallax.NewSchemaField("scm"),
			kallax.NewSchemaField("website"),
			kallax.NewSchemaField("name"),
			kallax.NewSchemaField("links"),
			kallax.NewSchemaField("fork_policy"),
			kallax.NewSchemaField("uuid"),
			kallax.NewSchemaField("language"),
			kallax.NewSchemaField("created_on"),
			kallax.NewSchemaField("parent"),
			kallax.NewSchemaField("full_name"),
			kallax.NewSchemaField("has_issues"),
			kallax.NewSchemaField("owner"),
			kallax.NewSchemaField("updated_on"),
			kallax.NewSchemaField("size"),
			kallax.NewSchemaField("type"),
			kallax.NewSchemaField("slug"),
			kallax.NewSchemaField("is_private"),
			kallax.NewSchemaField("description"),
		),
		ID:        kallax.NewSchemaField("id"),
		CreatedAt: kallax.NewSchemaField("created_at"),
		UpdatedAt: kallax.NewSchemaField("updated_at"),
		Next:      kallax.NewSchemaField("next"),
		Scm:       kallax.NewSchemaField("scm"),
		Website:   kallax.NewSchemaField("website"),
		Name:      kallax.NewSchemaField("name"),
		Links: &schemaRepositoryLinks{
			BaseSchemaField: kallax.NewSchemaField("links").(*kallax.BaseSchemaField),
			Clone: &schemaRepositoryLinksClone{
				JSONSchemaArray: kallax.NewJSONSchemaArray("links", "clone"),
				Href:            kallax.NewJSONSchemaKey(kallax.JSONText, "links", "clone", "href"),
				Name:            kallax.NewJSONSchemaKey(kallax.JSONText, "links", "clone", "name"),
			},
		},
		ForkPolicy: kallax.NewSchemaField("fork_policy"),
		UUID:       kallax.NewSchemaField("uuid"),
		Language:   kallax.NewSchemaField("language"),
		CreatedOn:  kallax.NewSchemaField("created_on"),
		Parent: &schemaRepositoryParent{
			BaseSchemaField: kallax.NewSchemaField("parent").(*kallax.BaseSchemaField),
			Type:            kallax.NewJSONSchemaKey(kallax.JSONText, "parent", "type"),
			Name:            kallax.NewJSONSchemaKey(kallax.JSONText, "parent", "name"),
			FullName:        kallax.NewJSONSchemaKey(kallax.JSONText, "parent", "full_name"),
			UUID:            kallax.NewJSONSchemaKey(kallax.JSONText, "parent", "uuid"),
		},
		FullName:  kallax.NewSchemaField("full_name"),
		HasIssues: kallax.NewSchemaField("has_issues"),
		Owner: &schemaRepositoryOwner{
			BaseSchemaField: kallax.NewSchemaField("owner").(*kallax.BaseSchemaField),
			Username:        kallax.NewJSONSchemaKey(kallax.JSONText, "owner", "username"),
			DisplayName:     kallax.NewJSONSchemaKey(kallax.JSONText, "owner", "display_name"),
			Type:            kallax.NewJSONSchemaKey(kallax.JSONText, "owner", "type"),
			UUID:            kallax.NewJSONSchemaKey(kallax.JSONText, "owner", "uuid"),
		},
		UpdatedOn:   kallax.NewSchemaField("updated_on"),
		Size:        kallax.NewSchemaField("size"),
		Type:        kallax.NewSchemaField("type"),
		Slug:        kallax.NewSchemaField("slug"),
		IsPrivate:   kallax.NewSchemaField("is_private"),
		Description: kallax.NewSchemaField("description"),
	},
}
