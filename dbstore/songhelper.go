package dbstore

import "fmt"

type SongFormatter interface {
	Format() string
}

func (s GetRandomSongRow) Format() string {
	return fmt.Sprintf("%s - %s", s.Song, s.Album)
}
