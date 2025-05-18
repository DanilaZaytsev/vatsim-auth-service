package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"vatsim-auth-service/pkg/logger"
)

type UserRepository struct {
	db table.Client
}

func NewUserRepository(db table.Client) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(ctx context.Context, cid uint64, email, fullName, refreshToken, roles string) error {
	query := `
DECLARE $cid AS Uint64;
DECLARE $email AS Utf8;
DECLARE $full_name AS Utf8;
DECLARE $refresh_token AS Utf8;
DECLARE $roles AS Utf8;

UPSERT INTO users (cid, email, full_name, refresh_token, roles)
VALUES ($cid, $email, $full_name, $refresh_token, $roles);
`

	params := table.NewQueryParameters(
		table.ValueParam("$cid", types.Uint64Value(cid)),
		table.ValueParam("$email", types.UTF8Value(email)),
		table.ValueParam("$full_name", types.UTF8Value(fullName)),
		table.ValueParam("$refresh_token", types.UTF8Value(refreshToken)),
		table.ValueParam("$roles", types.UTF8Value(roles)),
	)

	err := r.db.Do(ctx, func(ctx context.Context, sess table.Session) error {
		_, res, err := sess.Execute(ctx, table.DefaultTxControl(), query, params)
		if err != nil {
			logger.Error(err, "Failed to execute UPSERT query")
			return fmt.Errorf("failed to execute query: %w", err)
		}
		defer res.Close()
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateUserRole(ctx context.Context, cid uint64, newRole string) error {
	query := `
DECLARE $cid AS Uint64;
DECLARE $roles AS Utf8;
		UPDATE users
		SET roles = $roles
		WHERE cid = $cid;
	`

	params := table.NewQueryParameters(
		table.ValueParam("$cid", types.Uint64Value(cid)),
		table.ValueParam("$roles", types.UTF8Value(newRole)),
	)

	return r.db.Do(ctx, func(ctx context.Context, sess table.Session) error {
		_, res, err := sess.Execute(ctx, table.DefaultTxControl(), query, params)
		if err != nil {
			log.Printf("Failed to update role: %v", err)
			return err
		}
		defer res.Close()
		return nil
	})
}

func (r *UserRepository) GetUserRole(ctx context.Context, cid uint64) (string, error) {
	query := `
DECLARE $cid AS Uint64;
	SELECT roles
	FROM users
	WHERE cid = $cid;
	`

	params := table.NewQueryParameters(
		table.ValueParam("$cid", types.Uint64Value(cid)),
	)

	var role string

	err := r.db.Do(ctx, func(ctx context.Context, sess table.Session) error {
		_, res, err := sess.Execute(ctx, table.DefaultTxControl(), query, params)
		if err != nil {
			return fmt.Errorf("query execution error: %w", err)
		}
		defer res.Close()

		// Пролистываем сет результатов
		if !res.NextResultSet(ctx) {
			return fmt.Errorf("no result set")
		}

		if !res.NextRow() {
			return fmt.Errorf("user not found")
		}

		if err := res.ScanWithDefaults(&role); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}
	return role, nil
}
