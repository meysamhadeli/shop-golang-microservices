package gorm_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/uptrace/bun/driver/pgdriver"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Gorm struct {
	DB     *gorm.DB
	config *config.Config
}

func NewGorm(config *config.Config) (*gorm.DB, error) {

	var dataSourceName string

	if config.GormPostgres.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(config)

	if err != nil {
		return nil, err
	}

	dataSourceName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.GormPostgres.Host,
		config.GormPostgres.Port,
		config.GormPostgres.User,
		config.GormPostgres.DBName,
		config.GormPostgres.Password,
	)

	gormDb, err := gorm.Open(gorm_postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return gormDb, nil
}

func (db *Gorm) Close() {
	d, _ := db.DB.DB()
	_ = d.Close()
}

func createDB(cfg *config.Config) error {

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.GormPostgres.User,
		cfg.GormPostgres.Password,
		cfg.GormPostgres.Host,
		cfg.GormPostgres.Port,
		"postgres",
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	var exists int
	rows, err := sqldb.Query(fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.GormPostgres.DBName))
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = sqldb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.GormPostgres.DBName))
	if err != nil {
		return err
	}

	defer sqldb.Close()

	return nil
}

func Migrate(gorm *gorm.DB, types ...interface{}) error {

	for _, t := range types {
		err := gorm.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}

//Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {

	var items []T
	var totalRows int64
	db.Model(items).Count(&totalRows)

	// generate where query
	query := db.Offset(listQuery.GetOffset()).Limit(listQuery.GetLimit()).Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			column := filter.Field
			action := filter.Comparison
			value := filter.Value

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				query = query.Where(whereQuery, value)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.Where(whereQuery, "%"+value+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.Where(whereQuery, queryArray)
				break

			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, errors.Wrap(err, "error in finding products.")
	}

	return utils.NewListResult[T](items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}
