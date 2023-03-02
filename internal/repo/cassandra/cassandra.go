package cassandra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/google/uuid"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Cassandra struct {
	session *gocql.Session
	m       *migrate.Migrate
	cfg     *config.Config
}

func New(cfg *config.Config) (*Cassandra, error) {
	cluster := gocql.NewCluster(cfg.CASSANDRA_CLUSTER)
	cluster.Keyspace = cfg.CASSANDRA_KEYSPACE
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cfg.CASSANDRA_DB_USERNAME, Password: cfg.CASSANDRA_DB_PASSWORD}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}

	driver, err := cassandra.WithInstance(session, &cassandra.Config{KeyspaceName: cfg.CASSANDRA_KEYSPACE})
	if err != nil {
		return nil, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MIGRATE_PATH, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("new with database instance failed: %w", err)
	}

	return &Cassandra{session, m, cfg}, nil
}

func (c *Cassandra) Close() {
	c.session.Close()
}

func (c *Cassandra) CreateDriver(driver model.Driver) error {
	var name string
	err := c.session.Query("SELECT name FROM innotaxi.drivers WHERE (phone_number = ? OR email = ?) AND status = ?", driver.PhoneNumber, driver.Email, model.StatusCreated).Scan(&name)
	if err == nil {
		return fmt.Errorf("user: %v: %w", driver.Name, service.ErrDriverAlreadyExists)

	}

	err = c.session.Query("INSERT INTO innotaxi.drivers (id, name, phone_number, email, password, raiting, status) VALUES(?, ?, ?, ?, ?, 0.0, ?)", gocql.UUIDFromTime(time.Now()), driver.Name, driver.PhoneNumber, driver.Email, []byte(driver.Password), model.StatusCreated).Exec()
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (c *Cassandra) CheckUserByPhoneNumber(phone_number string) (*model.Driver, error) {
	var driver model.Driver
	var id gocql.UUID
	err := c.session.Query("SELECT id, phone_number, password FROM innotaxi.drivers WHERE phone_number = ? AND status = ? ALLOW FILTERING", phone_number, model.StatusCreated).Scan(&id, &driver.PhoneNumber, &driver.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrDriverDoesNotExists
		}

		return nil, fmt.Errorf("scan failed: %w", err)
	}
	driver.ID = uuid.UUID(id)
	return &driver, nil
}

func (c *Cassandra) DeleteDriverById(ctx context.Context, id string) error {
	err := c.session.Query("UPDATE innotaxi.drivers SET status = ? WHERE id = ?", model.StatusDeleted, id).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}
	return nil
}
