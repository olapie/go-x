package xmobile

import (
	"database/sql"
	"errors"
	"log/slog"

	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xsqlite"
)

type SecretManager interface {
	Get(key string) string
	Set(key, data string) bool
	Del(key string) bool
}

type secretManager struct {
	kv *xsqlite.KVTable
}

func NewSecretManager(db *sql.DB) SecretManager {
	return &secretManager{
		kv: xsqlite.NewKVTable(db, "secret_manager"),
	}
}

func (s *secretManager) Get(key string) string {
	v, err := s.kv.String(key)
	if err != nil {
		if errors.Is(err, xerror.DBNoRecords) {
			return ""
		}
		logger.Error("cannot read from db", slog.String("key", key))
	}
	return v
}

func (s *secretManager) Set(key, data string) bool {
	err := s.kv.SaveString(key, data)
	if err != nil {
		logger.Error("cannot save into db", slog.String("key", key))
		return false
	}
	return true
}

func (s *secretManager) Del(key string) bool {
	err := s.kv.Delete(key)
	if err != nil {
		logger.Error("cannot delete from db", slog.String("key", key))
		return false
	}
	return true
}
