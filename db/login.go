package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func SaveRefresh(ctx context.Context, jti, username string, issuedAt, exp time.Time) error {
	db, _ := GetDB()

	_, err := db.ExecContext(ctx,
		`INSERT INTO refresh_tokens(jti, username, expires_at, revoked, issued_at)
         VALUES(?, ?, ?, 0, ?)`,
		jti, username, exp.Unix(), issuedAt.Unix())
	log.Printf("Saved refresh token: %s, %s, %d, %d", jti, username, exp.Unix(), issuedAt.Unix())
	if err != nil {
		log.Printf("ERROR:Database: Error in saving refresh token. %s", err)
	}
	return err
}

func RevokeRefresh(ctx context.Context, jti string) error {
	db, _ := GetDB()

	_, err := db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked=1 WHERE jti=?`, jti)
	if err != nil {
		log.Printf("ERROR:Database: Error in revoking refresh token. %s", err)
	}
	log.Println("Revoked refresh token:", jti)
	return err
}

func SetValidAfterNow(ctx context.Context, username string) error {
	db, _ := GetDB()

	now := time.Now().UTC().Unix()
	_, err := db.ExecContext(ctx, `
		INSERT INTO user_session (username, valid_after)
		VALUES(?, ?)
		ON CONFLICT(username) DO UPDATE SET valid_after=excluded.valid_after
		`, username, now)
	if err != nil {
		log.Printf("ERROR:Database: Error in setting valid_after for user %s. %s", username, err)
	}
	log.Printf("Set valid_after for user %s to %d", username, now)
	return err
}

func IsRefreshValid(ctx context.Context, jti, username string, now time.Time) (bool, error) {
	db, _ := GetDB()

	var expires int64
	var revoked int
	var user string
	err := db.QueryRowContext(ctx,
		`SELECT username, expires_at, revoked
           FROM refresh_tokens
          WHERE jti=?`, jti).Scan(&user, &expires, &revoked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if revoked != 0 || user != username || now.Unix() >= expires {
		return false, nil
	}
	return true, nil
}

func RevokeAllRefreshForUser(ctx context.Context, username string) error {
	db, _ := GetDB()

	_, err := db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked=1 WHERE username=? AND revoked=0`, username)
	if err != nil {
		log.Printf("ERROR:Database: Error in revoking all refresh tokens for user %s. %s", username, err)
	}
	return err
}

func GetValidAfter(ctx context.Context, username string) (time.Time, error) {
	var ts int64
	db, _ := GetDB()
	err := db.QueryRowContext(ctx, `SELECT valid_after FROM user_session WHERE username=?`, username).Scan(&ts)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	if ts <= 0 {
		return time.Time{}, nil
	}
	return time.Unix(ts, 0).UTC(), nil
}
