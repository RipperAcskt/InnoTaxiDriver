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
	M       *migrate.Migrate
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

	m, err := migrate.NewWithDatabaseInstance(cfg.MIGRATE_PATH, "cassandra", driver)
	if err != nil {
		return nil, fmt.Errorf("new with database instance failed: %w", err)
	}

	return &Cassandra{session, m, cfg}, nil
}

func (c *Cassandra) Close() {
	c.session.Close()
}

func (c *Cassandra) CreateDriver(ctx context.Context, driver model.Driver) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var name string
	err := c.session.Query("SELECT name FROM innotaxi.drivers WHERE (phone_number = ? OR email = ?) AND status = ?", driver.PhoneNumber, driver.Email, model.StatusCreated).WithContext(queryCtx).Scan(&name)
	if err == nil {
		return fmt.Errorf("user: %v: %w", driver.Name, service.ErrDriverAlreadyExists)

	}

	err = c.session.Query("INSERT INTO innotaxi.drivers (id, name, phone_number, email, password, raiting, taxi_type, status) VALUES(?, ?, ?, ?, ?, 0.0, ?, ?)", gocql.UUIDFromTime(time.Now()), driver.Name, driver.PhoneNumber, driver.Email, []byte(driver.Password), driver.TaxiType, model.StatusCreated).Exec()
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (c *Cassandra) CheckUserByPhoneNumber(ctx context.Context, phone_number string) (*model.Driver, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var driver model.Driver
	var id gocql.UUID
	err := c.session.Query("SELECT id, phone_number, password FROM innotaxi.drivers WHERE phone_number = ? AND status = ? ALLOW FILTERING", phone_number, model.StatusCreated).WithContext(queryCtx).Scan(&id, &driver.PhoneNumber, &driver.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrDriverDoesNotExists
		}

		return nil, fmt.Errorf("scan failed: %w", err)
	}
	driver.ID = uuid.UUID(id)
	return &driver, nil
}

func (c *Cassandra) GetUserById(ctx context.Context, id string) (*model.Driver, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	driver := &model.Driver{}
	var driverID gocql.UUID
	err := c.session.Query("SELECT id, name, phone_number, email, raiting FROM innotaxi.drivers WHERE id = ? AND status = ? ALLOW FILTERING", id, model.StatusCreated).WithContext(queryCtx).Scan(&driverID, &driver.Name, &driver.PhoneNumber, &driver.Email, &driver.Raiting)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrDriverDoesNotExists
		}
		return nil, fmt.Errorf("query row context failed: %w", err)
	}
	driver.ID = uuid.UUID(driverID)
	return driver, err
}
