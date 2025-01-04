package dbstore

import "fmt"

type SongFormatter interface {
	Format() string
}

func (s GetRandomSongRow) Format() string {
	return fmt.Sprintf("%s - %s (%v)", s.Song, s.Album, s.Year)
}
