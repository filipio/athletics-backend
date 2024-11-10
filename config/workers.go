package config

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/filipio/athletics-backend/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"gorm.io/gorm"
)

// wrapper around river.Client so it is easier to insert jobs with gorm transactions
type InsertWorkerClient struct {
	*river.Client[*sql.Tx]
}

func (c *InsertWorkerClient) InsertTx(tx *gorm.DB, args river.JobArgs) (*rivertype.JobInsertResult, error) {
	ctx := context.Background() // for now new context is created for every job, will be changed as needed in future
	sqlTx := tx.Statement.ConnPool.(*sql.Tx)
	return c.Client.InsertTx(ctx, sqlTx, args, nil)
}

type WorkersClients struct {
	InsertClient    *InsertWorkerClient
	executionClient *river.Client[pgx.Tx]
}

func (c *WorkersClients) Shutdown(ctx context.Context) error {
	return c.executionClient.Stop(ctx)
}

var workersClientsInstance *WorkersClients
var workersOnce sync.Once

func SetupWorkersClient(ctx context.Context, db *gorm.DB, appWorkers *river.Workers) *WorkersClients {

	workersOnce.Do(func() {
		dbExecutionPool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
		if err != nil {
			log.Fatal("error creating database pool for river: ", err)
		}
		executionClient, err := river.NewClient(riverpgxv5.New(dbExecutionPool), &river.Config{
			Queues: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: 100},
			},
			Workers: appWorkers,
		})
		if err != nil {
			log.Fatal("error creating river execution client: ", err)
		}

		gormSqlDb, _ := db.DB()
		gormRiverClient, err := river.NewClient(riverdatabasesql.New(gormSqlDb), &river.Config{
			Workers: appWorkers,
		})
		if err != nil {
			log.Fatal("error creating river insert client: ", err)
		}

		// context with gorm db, so it can be used in workers code
		ctxWithGorm := context.WithValue(ctx, utils.DbContextKey, db)
		// starting workers goroutines to listen for jobs
		executionClient.Start(ctxWithGorm)

		insertClient := &InsertWorkerClient{gormRiverClient}

		workersClientsInstance = &WorkersClients{
			InsertClient:    insertClient,
			executionClient: executionClient,
		}
	})

	return workersClientsInstance
}
