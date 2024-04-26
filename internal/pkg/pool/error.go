package pool

import (
	"sync"

	"github.com/go-sql-driver/mysql"
	storage_go "github.com/supabase-community/storage-go"
)

var (
	MySQLErr = sync.Pool{
		New: func() interface{} { return new(mysql.MySQLError) },
	}
	SupabaseErr = sync.Pool{
		New: func() any { return new(storage_go.StorageError) },
	}
)
