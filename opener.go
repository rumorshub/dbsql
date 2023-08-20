package dbsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/azuread"
	_ "github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/sqliteshim"
)

var (
	_ Opener    = (*DBOpener)(nil)
	_ io.Closer = (*DBOpener)(nil)
	_ io.Closer = (*Channel)(nil)
)

var ErrConfigNotFound = errors.New("sql driver config not found")

type Opener interface {
	OpenDB(name string) (*sql.DB, string, error)
}

type Channel struct {
	io.Closer
	Config Config
	once   sync.Once
	db     *sql.DB
}

type DBOpener struct {
	io.Closer
	Channels map[string]*Channel
	mu       sync.Mutex
}

func NewOpener() *DBOpener {
	return &DBOpener{Channels: map[string]*Channel{}}
}

func (c *Channel) DB() (*sql.DB, string, error) {
	var err error
	c.once.Do(func() {
		c.db, err = sql.Open(c.Config.DriverName, c.Config.DataSourceName)
		if err != nil {
			return
		}

		if c.Config.MaxIdleConns != nil {
			c.db.SetMaxIdleConns(*c.Config.MaxIdleConns)
		}

		if c.Config.MaxOpenConns != nil {
			c.db.SetMaxOpenConns(*c.Config.MaxOpenConns)
		}

		if c.Config.ConnMaxLifetime != nil {
			c.db.SetConnMaxLifetime(*c.Config.ConnMaxLifetime)
		}

		if c.Config.ConnMaxIdleTime != nil {
			c.db.SetConnMaxIdleTime(*c.Config.ConnMaxIdleTime)
		}

		if c.Config.Ping {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()

			err = c.db.PingContext(ctx)
		}
	})
	return c.db, c.Config.DriverName, err
}

func (c *Channel) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}

func (o *DBOpener) AddChannel(name string, config Config) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.Channels[name] = &Channel{Config: config}
}

func (o *DBOpener) OpenDB(name string) (*sql.DB, string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if channel, ok := o.Channels[name]; ok {
		return channel.DB()
	}
	return nil, "", fmt.Errorf("%w: `%s`", ErrConfigNotFound, name)
}

func (o *DBOpener) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	var err error
	for _, channel := range o.Channels {
		if err1 := channel.Close(); err1 != nil {
			if err == nil {
				err = err1
			} else {
				err = fmt.Errorf("%w; %w", err, err1)
			}
		}
	}
	return err
}
