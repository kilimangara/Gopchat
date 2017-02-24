package builder

import (
	"gopkg.in/redis.v5"
	"fmt"
	"database/sql"
	_"github.com/lib/pq"
)

const(
	DATABASES_HOST="localhost"

	POSTGRES_DB="nikitazlain"
	POSTGRES_NAME="nikitazlain"
	POSTGRES_PASSWORD=""
	POSTGRES_PORT="5432"

	REDIS_PORT="6379"
	REDIS_PASSWORD=""
	REDIS_DB=0
)


func BuildRedisClient()(*redis.Client, error){
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", DATABASES_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB: REDIS_DB,
	})
}

func BuildPostgresClient()(*sql.DB, error){
	return sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s", POSTGRES_DB, POSTGRES_NAME, POSTGRES_PASSWORD, DATABASES_HOST, POSTGRES_PORT))
}
