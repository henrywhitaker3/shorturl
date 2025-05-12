package app

import (
	"database/sql"

	"github.com/henrywhitaker3/boiler"
	gocache "github.com/henrywhitaker3/go-cache"
	"github.com/henrywhitaker3/go-template/database/queries"
	"github.com/henrywhitaker3/go-template/internal/config"
	"github.com/henrywhitaker3/go-template/internal/crypto"
	ohttp "github.com/henrywhitaker3/go-template/internal/http"
	"github.com/henrywhitaker3/go-template/internal/metrics"
	"github.com/henrywhitaker3/go-template/internal/postgres"
	"github.com/henrywhitaker3/go-template/internal/probes"
	"github.com/henrywhitaker3/go-template/internal/queue"
	"github.com/henrywhitaker3/go-template/internal/redis"
	"github.com/henrywhitaker3/go-template/internal/storage"
	"github.com/henrywhitaker3/go-template/internal/urls"
	"github.com/henrywhitaker3/go-template/internal/workers"
	"github.com/redis/rueidis"
	"github.com/thanos-io/objstore"
)

func RegisterServe(b *boiler.Boiler) {
	RegisterBase(b)
	boiler.MustRegister(b, RegisterRunner)
	boiler.MustRegister(b, RegisterHTTP)
}

func RegisterBase(b *boiler.Boiler) {
	conf := boiler.MustResolve[*config.Config](b)

	boiler.MustRegister(b, RegisterProbes)
	if *conf.Telemetry.Metrics.Enabled {
		boiler.MustRegister(b, RegisterMetrics)
	}
	if *conf.Database.Enabled {
		boiler.MustRegister(b, RegisterDB)
		boiler.MustRegister(b, RegisterQueries)
		boiler.MustRegister(b, RegisterMigrator)
	}
	if *conf.Redis.Enabled {
		boiler.MustRegister(b, RegisterRedis)
		boiler.MustRegister(b, RegisterCache)
	}
	if *conf.Encryption.Enabled {
		boiler.MustRegister(b, RegisterEncryption)
	}
	if *conf.Storage.Enabled {
		boiler.MustRegister(b, RegisterStorage)
	}
	boiler.MustRegisterDeferred(b, RegisterUrls)
	if *conf.Queue.Enabled {
		boiler.MustRegister(b, RegisterQueue)
	}
}

func RegisterConsumers(b *boiler.Boiler) {
	RegisterBase(b)
	boiler.MustRegisterNamedDefered(b, DefaultQueue, RegisterDefaultQueueWorker)
	boiler.MustRegisterNamedDefered(b, CreateQueue, RegisterCreateQueueWorker)
}

func RegisterDB(b *boiler.Boiler) (*sql.DB, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	db, err := postgres.Open(b.Context(), conf.Database.Uri(), conf.Telemetry.Tracing)
	if err != nil {
		return nil, err
	}
	b.RegisterShutdown(func(b *boiler.Boiler) error {
		return db.Close()
	})
	return db, nil
}

func RegisterMigrator(b *boiler.Boiler) (*postgres.Migrator, error) {
	db, err := boiler.Resolve[*sql.DB](b)
	if err != nil {
		return nil, err
	}

	return postgres.NewMigrator(db)
}

func RegisterRedis(b *boiler.Boiler) (rueidis.Client, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	redis, err := redis.New(b.Context(), conf)
	if err != nil {
		return nil, err
	}
	b.RegisterShutdown(func(b *boiler.Boiler) error {
		redis.Close()
		return nil
	})
	return redis, nil
}

func RegisterCache(b *boiler.Boiler) (*gocache.Cache, error) {
	redis, err := boiler.Resolve[rueidis.Client](b)
	if err != nil {
		return nil, err
	}
	return gocache.NewCache(
		gocache.NewRueidisStore(redis),
	), nil
}

func RegisterQueries(b *boiler.Boiler) (*queries.Queries, error) {
	db, err := boiler.Resolve[*sql.DB](b)
	if err != nil {
		return nil, err
	}
	return queries.New(db), nil
}

func RegisterUrls(b *boiler.Boiler) (urls.Urls, error) {
	q, err := boiler.Resolve[*queries.Queries](b)
	if err != nil {
		return nil, err
	}
	return urls.New(urls.ServiceOpts{
		DB: q,
		Generator: urls.NewAliasGenerator(urls.AliasGeneratorOpts{
			DB: q,
		}),
	}), nil
}

func RegisterHTTP(b *boiler.Boiler) (*ohttp.Http, error) {
	return ohttp.New(b), nil
}

func RegisterEncryption(b *boiler.Boiler) (*crypto.Encrptor, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	return crypto.NewEncryptor(conf.Encryption.Secret)
}

func RegisterProbes(b *boiler.Boiler) (*probes.Probes, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	return probes.New(conf.Probes.Port), nil
}

func RegisterMetrics(b *boiler.Boiler) (*metrics.Metrics, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	return metrics.New(conf.Telemetry.Metrics.Port), nil
}

func RegisterQueue(b *boiler.Boiler) (*queue.Publisher, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	return queue.NewPublisher(queue.PublisherOpts{
		Redis: queue.RedisOpts{
			Addr:        conf.Redis.Addr,
			Password:    conf.Redis.Password,
			DB:          conf.Queue.DB,
			OtelEnabled: *conf.Telemetry.Tracing.Enabled,
		},
	})
}

func RegisterRunner(b *boiler.Boiler) (*workers.Runner, error) {
	redis, err := boiler.Resolve[rueidis.Client](b)
	if err != nil {
		return nil, err
	}
	return workers.NewRunner(b.Context(), redis)
}

func RegisterStorage(b *boiler.Boiler) (objstore.Bucket, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	return storage.New(conf.Storage)
}

const (
	DefaultQueue = "queue:default"
	CreateQueue  = "queue:create"
)

func RegisterDefaultQueueWorker(
	b *boiler.Boiler,
) (*queue.Worker, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	conc := 0
	if conf.Queue.Concurrency != nil {
		conc = *conf.Queue.Concurrency
	}
	return queue.NewWorker(b.Context(), queue.ServerOpts{
		Redis: queue.RedisOpts{
			Addr:        conf.Redis.Addr,
			Password:    conf.Redis.Password,
			DB:          conf.Queue.DB,
			OtelEnabled: *conf.Telemetry.Tracing.Enabled,
		},
		Queues:      []queue.Queue{queue.DefaultQueue},
		Concurrency: conc,
	})
}

func RegisterCreateQueueWorker(
	b *boiler.Boiler,
) (*queue.Worker, error) {
	conf, err := boiler.Resolve[*config.Config](b)
	if err != nil {
		return nil, err
	}
	conc := 0
	if conf.Queue.Concurrency != nil {
		conc = *conf.Queue.Concurrency
	}

	svc, err := boiler.Resolve[urls.Urls](b)
	if err != nil {
		return nil, err
	}
	handler := urls.NewCreateJobHandler(svc)

	worker, err := queue.NewWorker(b.Context(), queue.ServerOpts{
		Redis: queue.RedisOpts{
			Addr:        conf.Redis.Addr,
			Password:    conf.Redis.Password,
			DB:          conf.Queue.DB,
			OtelEnabled: *conf.Telemetry.Tracing.Enabled,
		},
		Queues:      []queue.Queue{queue.Create},
		Concurrency: conc,
	})
	if err != nil {
		return nil, err
	}
	worker.RegisterHandler(queue.CreateTask, handler)
	return worker, nil
}
