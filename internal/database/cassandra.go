package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/codevault-llc/fingerprint/config"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type Database struct {
	Db *gocql.Session
}

func NewDatabase() (*Database, error) {
	cluster := gocql.NewCluster(config.Config.DatabaseHost)
	cluster.Consistency = gocql.Quorum

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Config.DatabaseUser,
		Password: config.Config.DatabasePass,
	}
	cluster.Port = 9042

	tempSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer tempSession.Close()

	err = createSchema(tempSession)
	if err != nil {
		return nil, err
	}

	cluster.Keyspace = config.Config.DatabaseKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Log.Error("Failed to create session", zap.Error(err))
		return nil, err
	}

	return &Database{Db: session}, nil
}

func createSchema(session *gocql.Session) error {
	err := createKeyspace(session)
	if err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}

	err = session.Query(entities.CreateFingerprintSchema()).Exec()
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	return nil
}

func createKeyspace(session *gocql.Session) error {
	query := `
	CREATE KEYSPACE IF NOT EXISTS fingerprint
	WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`

	return session.Query(query).Exec()
}

func (d *Database) BulkInsert(fingerprints []entities.Fingerprint) error {
	batch := d.Db.NewBatch(gocql.LoggedBatch)

	for _, fingerprint := range fingerprints {
		batch.Query(entities.InsertFingerprintQuery(), fingerprint.Id, fingerprint.Name, fingerprint.Description, fingerprint.Pattern, fingerprint.Type, fingerprint.Keywords, fingerprint.CreatedAt, fingerprint.UpdatedAt)
	}

	if err := d.Db.ExecuteBatch(batch); err != nil {
		logger.Log.Error("Failed to execute batch insert", zap.Error(err))
		return err
	}

	return nil
}

func (d *Database) TableExists(table string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM system_schema.tables WHERE keyspace_name = '%s' AND table_name = '%s'", config.Config.DatabaseKeyspace, table)
	iter := d.Db.Query(query).Iter()
	defer iter.Close()

	var count int
	iter.Scan(&count)

	if count > 0 {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", config.Config.DatabaseKeyspace, table)
		iter = d.Db.Query(query).Iter()
		defer iter.Close()

		iter.Scan(&count)

		return count > 0
	}

	return false
}

func (d *Database) Close() {
	d.Db.Close()
}

func (d *Database) GetDatabase() *gocql.Session {
	return d.Db
}

func (db *Database) Select(ctx context.Context, query string, result interface{}, args ...interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a pointer to a slice")
	}

	iter := db.Db.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	sliceElemType := rv.Elem().Type().Elem()

	if sliceElemType.Kind() == reflect.Struct {
		fieldMap := make(map[string]int)
		for i := 0; i < sliceElemType.NumField(); i++ {
			field := sliceElemType.Field(i)
			dbTag := field.Tag.Get("ctx")
			if dbTag != "" {
				fieldMap[dbTag] = i
			}
		}

		for {
			columns := iter.Columns()
			row := reflect.New(sliceElemType).Elem()
			fieldValues := make([]interface{}, len(columns))

			for i, column := range columns {
				if fieldIndex, ok := fieldMap[column.Name]; ok {
					field := row.Field(fieldIndex)

					switch field.Type() {
					case reflect.TypeOf(time.Time{}):
						var temp int64
						fieldValues[i] = &temp
					default:
						fieldValues[i] = field.Addr().Interface()
					}
				} else {
					var temp interface{}
					fieldValues[i] = &temp
				}
			}

			if !iter.Scan(fieldValues...) {
				break
			}

			for i, column := range columns {
				if fieldIndex, ok := fieldMap[column.Name]; ok {
					field := row.Field(fieldIndex)
					if field.Type() == reflect.TypeOf(time.Time{}) {
						intValue := *fieldValues[i].(*int64)
						field.Set(reflect.ValueOf(time.Unix(intValue, 0)))
					}
				}
			}

			rv.Elem().Set(reflect.Append(rv.Elem(), row))
		}
	} else {
		for {
			row := reflect.New(sliceElemType).Elem().Addr().Interface()
			if !iter.Scan(row) {
				break
			}
			rv.Elem().Set(reflect.Append(rv.Elem(), reflect.ValueOf(row).Elem()))
		}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}
