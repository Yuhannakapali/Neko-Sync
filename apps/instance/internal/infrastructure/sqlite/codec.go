package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"
)

func toNanos(t time.Time) int64 { return t.UnixNano() }

func fromNanos(n int64) time.Time { return time.Unix(0, n).UTC() }

func nullInt(p *int) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}

func nullFloat(p *float64) sql.NullFloat64 {
	if p == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *p, Valid: true}
}

func intFromNull(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}

func floatFromNull(n sql.NullFloat64) *float64 {
	if !n.Valid {
		return nil
	}
	v := n.Float64
	return &v
}

func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
