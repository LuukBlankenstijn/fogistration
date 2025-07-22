package syncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client/wrapper"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DomJudgeSyncer struct {
	ctx    context.Context
	client *wrapper.Client
	db     *pgxpool.Pool
}

func NewSyncer(ctx context.Context, client *wrapper.Client, db *pgxpool.Pool) *DomJudgeSyncer {
	return &DomJudgeSyncer{
		ctx:    ctx,
		client: client,
		db:     db,
	}
}

func (s *DomJudgeSyncer) Sync() error {
	tx, err := s.db.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(s.ctx)
		if err != nil && err != pgx.ErrTxClosed {
			logging.Error("failed to rollback transaction: %w", err)
		}
	}()
	queries := database.New(tx)

	err = s.syncContests(queries)
	if err != nil {
		return err
	}

	err = s.syncTeams(queries)
	if err != nil {
		return err
	}

	if err := tx.Commit(s.ctx); err != nil {
		logging.Error("failed to commit transaction: %w", err)
	}

	return nil
}

func computeHash(data any) string {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// Handle pointer types
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	hashData := make(map[string]any)

	for i := range val.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Check if field has hash:"exclude" tag
		hashTag := field.Tag.Get("hash")
		if hashTag == "exclude" {
			continue // Skip this field
		}

		// Get field name for JSON (use json tag if available, otherwise field name)
		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		// Remove ",omitempty" etc from json tag
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}

		hashData[fieldName] = fieldValue.Interface()
	}

	jsonData, _ := json.Marshal(hashData)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}
