package mackerel

import "context"

// Invitation information
type Invitation struct {
	Email     string `json:"email,omitempty"`
	Authority string `json:"authority,omitempty"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

// FindInvitations finds invitations.
func (c *Client) FindInvitations() ([]*Invitation, error) {
	return c.FindInvitationsContext(context.Background())
}

// FindInvitationsContext finds invitations.
func (c *Client) FindInvitationsContext(ctx context.Context) ([]*Invitation, error) {
	data, err := requestGetContext[struct {
		Invitations []*Invitation `json:"invitations"`
	}](ctx, c, "/api/v0/invitations")
	if err != nil {
		return nil, err
	}
	return data.Invitations, nil
}

// CreateInvitation creates a invitation.
func (c *Client) CreateInvitation(param *Invitation) (*Invitation, error) {
	return c.CreateInvitationContext(context.Background(), param)
}

// CreateInvitationContext creates a invitation.
func (c *Client) CreateInvitationContext(ctx context.Context, param *Invitation) (*Invitation, error) {
	return requestPostContext[Invitation](ctx, c, "/api/v0/invitations", param)
}
