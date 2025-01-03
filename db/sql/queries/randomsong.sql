-- name: GetRandomSong :one
SELECT s.title as song, a.title as album FROM song s
         INNER JOIN main.album a on a.id = s.album_id
    ORDER BY RANDOM() LIMIT 1;