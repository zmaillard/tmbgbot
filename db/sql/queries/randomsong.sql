-- name: GetRandomSong :one
SELECT * FROM song ORDER BY RANDOM() LIMIT 1;