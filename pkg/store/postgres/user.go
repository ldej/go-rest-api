package postgres

import (
	"context"
	"database/sql"

	"github.com/ldej/go-rest-example/internal/api"
	"github.com/ldej/go-rest-example/pkg/store"
)

func (c *Client) UserGetByUID(ctx context.Context, uid string) (*api.User, error) {
	u := new(api.User)
	err := c.db.QueryRowContext(ctx, `SELECT uid, name, email_address, password FROM users WHERE uid = $1`, uid).Scan(u.UID, u.Name, u.EmailAddress, u.EncryptedPassword)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) UserCreate(ctx context.Context, user *api.User) (*api.User, error) {
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO users (uid, name, email_address, password)
		VALUES($1, $2, $3, $4)
	`, user.UID, user.Name, user.EmailAddress, user.EncryptedPassword)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (c *Client) UserGetByEmailAddress(ctx context.Context, emailAddress string) (*api.User, error) {
	u := new(api.User)
	err := c.db.QueryRowContext(ctx, `SELECT uid, name, email_address, password FROM users WHERE email_address = $1`, emailAddress).Scan(u.UID, u.Name, u.EmailAddress, u.EncryptedPassword)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, nil
}
