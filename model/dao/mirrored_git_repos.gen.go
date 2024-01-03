// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/yankeguo/nas-tools/model"
)

func newMirroredGitRepo(db *gorm.DB, opts ...gen.DOOption) mirroredGitRepo {
	_mirroredGitRepo := mirroredGitRepo{}

	_mirroredGitRepo.mirroredGitRepoDo.UseDB(db, opts...)
	_mirroredGitRepo.mirroredGitRepoDo.UseModel(&model.MirroredGitRepo{})

	tableName := _mirroredGitRepo.mirroredGitRepoDo.TableName()
	_mirroredGitRepo.ALL = field.NewAsterisk(tableName)
	_mirroredGitRepo.Key = field.NewString(tableName, "key")
	_mirroredGitRepo.LastGCAt = field.NewTime(tableName, "last_gc_at")
	_mirroredGitRepo.LastCommitAt = field.NewTime(tableName, "last_commit_at")
	_mirroredGitRepo.LastCommitBy = field.NewString(tableName, "last_commit_by")
	_mirroredGitRepo.LastCommitMessage = field.NewString(tableName, "last_commit_message")
	_mirroredGitRepo.UpdatedAt = field.NewTime(tableName, "updated_at")

	_mirroredGitRepo.fillFieldMap()

	return _mirroredGitRepo
}

type mirroredGitRepo struct {
	mirroredGitRepoDo

	ALL               field.Asterisk
	Key               field.String
	LastGCAt          field.Time
	LastCommitAt      field.Time
	LastCommitBy      field.String
	LastCommitMessage field.String
	UpdatedAt         field.Time

	fieldMap map[string]field.Expr
}

func (m mirroredGitRepo) Table(newTableName string) *mirroredGitRepo {
	m.mirroredGitRepoDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m mirroredGitRepo) As(alias string) *mirroredGitRepo {
	m.mirroredGitRepoDo.DO = *(m.mirroredGitRepoDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *mirroredGitRepo) updateTableName(table string) *mirroredGitRepo {
	m.ALL = field.NewAsterisk(table)
	m.Key = field.NewString(table, "key")
	m.LastGCAt = field.NewTime(table, "last_gc_at")
	m.LastCommitAt = field.NewTime(table, "last_commit_at")
	m.LastCommitBy = field.NewString(table, "last_commit_by")
	m.LastCommitMessage = field.NewString(table, "last_commit_message")
	m.UpdatedAt = field.NewTime(table, "updated_at")

	m.fillFieldMap()

	return m
}

func (m *mirroredGitRepo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *mirroredGitRepo) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 6)
	m.fieldMap["key"] = m.Key
	m.fieldMap["last_gc_at"] = m.LastGCAt
	m.fieldMap["last_commit_at"] = m.LastCommitAt
	m.fieldMap["last_commit_by"] = m.LastCommitBy
	m.fieldMap["last_commit_message"] = m.LastCommitMessage
	m.fieldMap["updated_at"] = m.UpdatedAt
}

func (m mirroredGitRepo) clone(db *gorm.DB) mirroredGitRepo {
	m.mirroredGitRepoDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m mirroredGitRepo) replaceDB(db *gorm.DB) mirroredGitRepo {
	m.mirroredGitRepoDo.ReplaceDB(db)
	return m
}

type mirroredGitRepoDo struct{ gen.DO }

func (m mirroredGitRepoDo) Debug() *mirroredGitRepoDo {
	return m.withDO(m.DO.Debug())
}

func (m mirroredGitRepoDo) WithContext(ctx context.Context) *mirroredGitRepoDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m mirroredGitRepoDo) ReadDB() *mirroredGitRepoDo {
	return m.Clauses(dbresolver.Read)
}

func (m mirroredGitRepoDo) WriteDB() *mirroredGitRepoDo {
	return m.Clauses(dbresolver.Write)
}

func (m mirroredGitRepoDo) Session(config *gorm.Session) *mirroredGitRepoDo {
	return m.withDO(m.DO.Session(config))
}

func (m mirroredGitRepoDo) Clauses(conds ...clause.Expression) *mirroredGitRepoDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m mirroredGitRepoDo) Returning(value interface{}, columns ...string) *mirroredGitRepoDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m mirroredGitRepoDo) Not(conds ...gen.Condition) *mirroredGitRepoDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m mirroredGitRepoDo) Or(conds ...gen.Condition) *mirroredGitRepoDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m mirroredGitRepoDo) Select(conds ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m mirroredGitRepoDo) Where(conds ...gen.Condition) *mirroredGitRepoDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m mirroredGitRepoDo) Order(conds ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m mirroredGitRepoDo) Distinct(cols ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m mirroredGitRepoDo) Omit(cols ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m mirroredGitRepoDo) Join(table schema.Tabler, on ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m mirroredGitRepoDo) LeftJoin(table schema.Tabler, on ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m mirroredGitRepoDo) RightJoin(table schema.Tabler, on ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m mirroredGitRepoDo) Group(cols ...field.Expr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m mirroredGitRepoDo) Having(conds ...gen.Condition) *mirroredGitRepoDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m mirroredGitRepoDo) Limit(limit int) *mirroredGitRepoDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m mirroredGitRepoDo) Offset(offset int) *mirroredGitRepoDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m mirroredGitRepoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *mirroredGitRepoDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m mirroredGitRepoDo) Unscoped() *mirroredGitRepoDo {
	return m.withDO(m.DO.Unscoped())
}

func (m mirroredGitRepoDo) Create(values ...*model.MirroredGitRepo) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m mirroredGitRepoDo) CreateInBatches(values []*model.MirroredGitRepo, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m mirroredGitRepoDo) Save(values ...*model.MirroredGitRepo) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m mirroredGitRepoDo) First() (*model.MirroredGitRepo, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MirroredGitRepo), nil
	}
}

func (m mirroredGitRepoDo) Take() (*model.MirroredGitRepo, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MirroredGitRepo), nil
	}
}

func (m mirroredGitRepoDo) Last() (*model.MirroredGitRepo, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MirroredGitRepo), nil
	}
}

func (m mirroredGitRepoDo) Find() ([]*model.MirroredGitRepo, error) {
	result, err := m.DO.Find()
	return result.([]*model.MirroredGitRepo), err
}

func (m mirroredGitRepoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MirroredGitRepo, err error) {
	buf := make([]*model.MirroredGitRepo, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m mirroredGitRepoDo) FindInBatches(result *[]*model.MirroredGitRepo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m mirroredGitRepoDo) Attrs(attrs ...field.AssignExpr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m mirroredGitRepoDo) Assign(attrs ...field.AssignExpr) *mirroredGitRepoDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m mirroredGitRepoDo) Joins(fields ...field.RelationField) *mirroredGitRepoDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m mirroredGitRepoDo) Preload(fields ...field.RelationField) *mirroredGitRepoDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m mirroredGitRepoDo) FirstOrInit() (*model.MirroredGitRepo, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MirroredGitRepo), nil
	}
}

func (m mirroredGitRepoDo) FirstOrCreate() (*model.MirroredGitRepo, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MirroredGitRepo), nil
	}
}

func (m mirroredGitRepoDo) FindByPage(offset int, limit int) (result []*model.MirroredGitRepo, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m mirroredGitRepoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m mirroredGitRepoDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m mirroredGitRepoDo) Delete(models ...*model.MirroredGitRepo) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *mirroredGitRepoDo) withDO(do gen.Dao) *mirroredGitRepoDo {
	m.DO = *do.(*gen.DO)
	return m
}
