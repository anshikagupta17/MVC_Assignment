package repositories

import (
	"context"
	"fmt"

	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type MockRow struct {
	values []any
	err    error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	for i, d := range dest {
		switch v := d.(type) {
		case *int:
			*v = m.values[i].(int)
		case *int64:
			*v = m.values[i].(int64)
		case *float64:
			*v = m.values[i].(float64)
		case *string:
			*v = m.values[i].(string)
		case *pgtype.Timestamp:
			*v = m.values[i].(pgtype.Timestamp)
		case **time.Time:
			val := m.values[i]
			if val == nil {
				*v = nil
			} else {
				switch t := val.(type) {
				case *time.Time:
					*v = t
				case time.Time:
					*v = &t
				}
			}
		}
	}
	return nil
}

type MockRows struct {
	rows  [][]any
	index int
}

func (m *MockRows) Next() bool {
	return m.index < len(m.rows)
}

func (m *MockRows) Scan(dest ...any) error {
	current := m.rows[m.index]
	m.index++
	for i, d := range dest {
		switch v := d.(type) {
		case *int:
			*v = current[i].(int)
		case *int64:
			*v = current[i].(int64)
		case *float64:
			*v = current[i].(float64)
		case *pgtype.Timestamp:
			*v = current[i].(pgtype.Timestamp)
		case **time.Time:
			val := current[i]
			if val == nil {
				*v = nil
			} else {
				switch t := val.(type) {
				case *time.Time:
					*v = t
				case time.Time:
					*v = &t
				}
			}
		}
	}
	return nil
}

func (m *MockRows) Close() {}

func (m *MockRows) Err() error { return nil }

func (m *MockRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }

func (m *MockRows) Values() ([]any, error) { return nil, nil }

func (m *MockRows) RawValues() [][]byte { return nil }

func (m *MockRows) Conn() *pgx.Conn { return nil }

type MockDBExecutor struct {
	queryRows        []*MockRow
	queryRowIndex    int
	queryResults     []*MockRows
	queryIndex       int
	execErr          error
	queryErr         error
	execRowsAffected int64
}

func (m *MockDBExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	row := m.queryRows[m.queryRowIndex]
	m.queryRowIndex++
	return row
}

func (m *MockDBExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}

	result := m.queryResults[m.queryIndex]
	m.queryIndex++
	return result, nil
}

func (m *MockDBExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if m.execErr != nil {
		return pgconn.CommandTag{}, m.execErr
	}
	return pgconn.NewCommandTag(fmt.Sprintf("UPDATE %d", m.execRowsAffected)), nil

}
