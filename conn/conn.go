package conn

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tianhongw/tinyid/internal/util"
)

type DB struct {
	dbs []*sqlx.DB
}

func (db *DB) GetConn() *sqlx.DB {
	if len(db.dbs) == 1 {
		return db.dbs[0]
	}

	return db.dbs[util.RandomInt(len(db.dbs))]
}

type options struct {
	Configs      []config
	MaxIdleConns int
	MaxOpenConns int
}

type config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Option func(*options)

func WithAddressAndAuth(hosts []string, ports []int,
	users []string, passwords []string, dbNames []string) Option {
	return func(o *options) {
		configs := make([]config, 0, len(hosts))
		for i := 0; i < len(hosts); i++ {
			configs = append(configs, config{
				Host:     hosts[i],
				Port:     ports[i],
				User:     users[i],
				Password: passwords[i],
				DBName:   dbNames[i],
			})
		}
		o.Configs = configs
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(o *options) {
		o.MaxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(o *options) {
		o.MaxOpenConns = maxOpenConns
	}
}

func NewConn(opts ...Option) (*DB, error) {
	opt := &options{
		Configs: []config{
			{
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "",
				Password: "",
				DBName:   "",
			},
		},
		MaxIdleConns: 10,
		MaxOpenConns: 10,
	}

	for _, o := range opts {
		o(opt)
	}

	ret := &DB{}

	for _, conf := range opt.Configs {
		db, err := sqlx.Connect("mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC&charset=utf8mb4&collation=utf8mb4_general_ci",
			conf.User, conf.Password, conf.Host, conf.Port, conf.DBName,
		))

		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(opt.MaxIdleConns)
		db.SetMaxOpenConns(opt.MaxOpenConns)

		ret.dbs = append(ret.dbs, db)
	}

	return ret, nil
}
