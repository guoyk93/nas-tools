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

func newPhotoFile(db *gorm.DB, opts ...gen.DOOption) photoFile {
	_photoFile := photoFile{}

	_photoFile.photoFileDo.UseDB(db, opts...)
	_photoFile.photoFileDo.UseModel(&model.PhotoFile{})

	tableName := _photoFile.photoFileDo.TableName()
	_photoFile.ALL = field.NewAsterisk(tableName)
	_photoFile.ID = field.NewString(tableName, "id")
	_photoFile.Group_ = field.NewString(tableName, "group")
	_photoFile.Path = field.NewString(tableName, "path")
	_photoFile.Md5 = field.NewString(tableName, "md5")

	_photoFile.fillFieldMap()

	return _photoFile
}

type photoFile struct {
	photoFileDo

	ALL    field.Asterisk
	ID     field.String
	Group_ field.String
	Path   field.String
	Md5    field.String

	fieldMap map[string]field.Expr
}

func (p photoFile) Table(newTableName string) *photoFile {
	p.photoFileDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p photoFile) As(alias string) *photoFile {
	p.photoFileDo.DO = *(p.photoFileDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *photoFile) updateTableName(table string) *photoFile {
	p.ALL = field.NewAsterisk(table)
	p.ID = field.NewString(table, "id")
	p.Group_ = field.NewString(table, "group")
	p.Path = field.NewString(table, "path")
	p.Md5 = field.NewString(table, "md5")

	p.fillFieldMap()

	return p
}

func (p *photoFile) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *photoFile) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 4)
	p.fieldMap["id"] = p.ID
	p.fieldMap["group"] = p.Group_
	p.fieldMap["path"] = p.Path
	p.fieldMap["md5"] = p.Md5
}

func (p photoFile) clone(db *gorm.DB) photoFile {
	p.photoFileDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p photoFile) replaceDB(db *gorm.DB) photoFile {
	p.photoFileDo.ReplaceDB(db)
	return p
}

type photoFileDo struct{ gen.DO }

func (p photoFileDo) Debug() *photoFileDo {
	return p.withDO(p.DO.Debug())
}

func (p photoFileDo) WithContext(ctx context.Context) *photoFileDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p photoFileDo) ReadDB() *photoFileDo {
	return p.Clauses(dbresolver.Read)
}

func (p photoFileDo) WriteDB() *photoFileDo {
	return p.Clauses(dbresolver.Write)
}

func (p photoFileDo) Session(config *gorm.Session) *photoFileDo {
	return p.withDO(p.DO.Session(config))
}

func (p photoFileDo) Clauses(conds ...clause.Expression) *photoFileDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p photoFileDo) Returning(value interface{}, columns ...string) *photoFileDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p photoFileDo) Not(conds ...gen.Condition) *photoFileDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p photoFileDo) Or(conds ...gen.Condition) *photoFileDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p photoFileDo) Select(conds ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p photoFileDo) Where(conds ...gen.Condition) *photoFileDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p photoFileDo) Order(conds ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p photoFileDo) Distinct(cols ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p photoFileDo) Omit(cols ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p photoFileDo) Join(table schema.Tabler, on ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p photoFileDo) LeftJoin(table schema.Tabler, on ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p photoFileDo) RightJoin(table schema.Tabler, on ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p photoFileDo) Group(cols ...field.Expr) *photoFileDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p photoFileDo) Having(conds ...gen.Condition) *photoFileDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p photoFileDo) Limit(limit int) *photoFileDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p photoFileDo) Offset(offset int) *photoFileDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p photoFileDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *photoFileDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p photoFileDo) Unscoped() *photoFileDo {
	return p.withDO(p.DO.Unscoped())
}

func (p photoFileDo) Create(values ...*model.PhotoFile) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p photoFileDo) CreateInBatches(values []*model.PhotoFile, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p photoFileDo) Save(values ...*model.PhotoFile) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p photoFileDo) First() (*model.PhotoFile, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.PhotoFile), nil
	}
}

func (p photoFileDo) Take() (*model.PhotoFile, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.PhotoFile), nil
	}
}

func (p photoFileDo) Last() (*model.PhotoFile, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.PhotoFile), nil
	}
}

func (p photoFileDo) Find() ([]*model.PhotoFile, error) {
	result, err := p.DO.Find()
	return result.([]*model.PhotoFile), err
}

func (p photoFileDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PhotoFile, err error) {
	buf := make([]*model.PhotoFile, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p photoFileDo) FindInBatches(result *[]*model.PhotoFile, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p photoFileDo) Attrs(attrs ...field.AssignExpr) *photoFileDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p photoFileDo) Assign(attrs ...field.AssignExpr) *photoFileDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p photoFileDo) Joins(fields ...field.RelationField) *photoFileDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p photoFileDo) Preload(fields ...field.RelationField) *photoFileDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p photoFileDo) FirstOrInit() (*model.PhotoFile, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.PhotoFile), nil
	}
}

func (p photoFileDo) FirstOrCreate() (*model.PhotoFile, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.PhotoFile), nil
	}
}

func (p photoFileDo) FindByPage(offset int, limit int) (result []*model.PhotoFile, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p photoFileDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p photoFileDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p photoFileDo) Delete(models ...*model.PhotoFile) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *photoFileDo) withDO(do gen.Dao) *photoFileDo {
	p.DO = *do.(*gen.DO)
	return p
}