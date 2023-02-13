package cassandra

import (
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"

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

func (c *Cassandra) CreateDriver(driver auth.PostDriverSingUpBody) error {
	var name string
	err := c.session.Query("SELECT name FROM innotaxi.drivers WHERE (phone_number = ? OR email = ?) AND status = ?", driver.PhoneNumber, driver.Email, service.StatusCreated).Scan(&name)
	if err == nil {
		return fmt.Errorf("user: %v: %w", driver.Name, service.ErrUserAlreadyExists)

	}

	err = c.session.Query("INSERT INTO innotaxi.drivers (name, phone_number, email, password, rating, status) VALUES(?, ?, ?, ?, 0.0, ?)", driver.Name, driver.PhoneNumber, driver.Email, []byte(driver.Password), service.StatusCreated).Exec()
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}
