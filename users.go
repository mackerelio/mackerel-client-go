package mackerel

import (
	"context"
	"fmt"
)

// User information
type User struct {
	ID         string `json:"id,omitempty"`
	ScreenName string `json:"screenName,omitempty"`
	Email      string `json:"email,omitempty"`
	Authority  string `json:"authority,omitempty"`

	IsInRegistrationProcess bool     `json:"isInRegistrationProcess,omitempty"`
	IsMFAEnabled            bool     `json:"isMFAEnabled,omitempty"`
	AuthenticationMethods   []string `json:"authenticationMethods,omitempty"`
	JoinedAt                int64    `json:"joinedAt,omitempty"`
}

// FindUsers finds users.
func (c *Client) FindUsers() ([]*User, error) {
	return c.FindUsersContext(context.Background())
}

// FindUsersContext finds users.
func (c *Client) FindUsersContext(ctx context.Context) ([]*User, error) {
	data, err := requestGetContext[struct {
		Users []*User `json:"users"`
	}](ctx, c, "/api/v0/users")
	if err != nil {
		return nil, err
	}
	return data.Users, nil
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(userID string) (*User, error) {
	path := fmt.Sprintf("/api/v0/users/%s", userID)
	return requestDeleteContext[User](context.Background(), c, path)
}

// DeleteUserContext is like [DeleteUser].
func (c *Client) DeleteUserContext(ctx context.Context, userID string) (*User, error) {
	path := fmt.Sprintf("/api/v0/users/%s", userID)
	return requestDeleteContext[User](ctx, c, path)
}
