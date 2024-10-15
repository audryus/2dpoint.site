package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/pkg/database/cockroach"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cockroachdb"
	"github.com/testcontainers/testcontainers-go/modules/etcd"
	"github.com/z0ne-dev/mgx/v2"
)

var cockroachDBContainer *cockroachdb.CockroachDBContainer
var etcdContainer testcontainers.Container

// TestMain controls main for the tests and allows for setup and shutdown of tests
func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	// COCKROACH
	envs := testcontainers.WithEnv(map[string]string{"COCKROACH_DATABASE": "cockroach"})

	cockroachDBContainer, err = cockroachdb.Run(ctx, "cockroachdb/cockroach:latest-v23.1", envs, cockroachdb.WithDatabase("cockroach"))
	if err != nil {
		fmt.Println("Could not start cockroach Container")
		os.Exit(1)
	}

	port, err := cockroachDBContainer.MappedPort(ctx, "26257")
	if err != nil {
		fmt.Println("Could not obtain cockroach port")
		os.Exit(1)
	}
	os.Setenv("DPOINT_COCKROACH_PORT", port.Port())

	cockroachDBContainer.Exec(ctx, []string{"sh", "-c", "./cockroach sql --insecure"})

	// ETCD
	envsEtcd := testcontainers.WithEnv(map[string]string{
		"ALLOW_NONE_AUTHENTICATION":  "yes",
		"ETCD_ADVERTISE_CLIENT_URLS": "http://0.0.0.0:2379",
		"ETCD_LISTEN_CLIENT_URLS":    "http://0.0.0.0:2379"})
	etcdContainer, err := etcd.Run(ctx, "quay.io/coreos/etcd:v3.5.16", envsEtcd)

	if err != nil {
		fmt.Println("Could not start etcd Container")
		os.Exit(1)
	}
	time.Sleep(1 * time.Second)
	port, err = etcdContainer.MappedPort(ctx, "2379")
	if err != nil {
		log.Fatal("Could not obtain etcd port ", err)
		os.Exit(1)
	}
	os.Setenv("DPOINT_ETCD_PORT", port.Port())

	l := logger.New()

	conf, err := config.New(l)
	if err != nil {
		log.Fatal("Could not create config ", err)
		os.Exit(1)
	}

	db, err := cockroach.New(conf, l)
	if err != nil {
		fmt.Println("Could not obtain cockroach client")
		os.Exit(1)
	}

	migrator, err := mgx.New(mgx.Migrations(
		mgx.NewRawMigration(
			"raw migration: create database",
			"create database if not exists cockroach"),
		mgx.NewRawMigration(
			"raw migration: use database",
			"use cockroach"),
		mgx.NewRawMigration(
			"raw migration: create table text",
			"CREATE TABLE if not exists cockroach.text (id varchar(16) PRIMARY KEY, hash_code varchar(40) NOT NULL, status varchar(40) not null, kind varchar(32) NOT NULL, content text NOT NULL)"),
	))

	err = db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		return migrator.Migrate(ctx, tx)
	})
	if err != nil {
		fmt.Println("migration error")
		log.Fatal(err)
		os.Exit(1)
	}

	//Catching all panics to once again make sure that shutDown is successfully run
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic")
		}
	}()

	code := m.Run()
	os.Exit(code)
}
