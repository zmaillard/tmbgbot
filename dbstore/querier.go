// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbstore

import (
	"context"
)

type Querier interface {
	GetRandomSong(ctx context.Context) (GetRandomSongRow, error)
}

var _ Querier = (*Queries)(nil)
