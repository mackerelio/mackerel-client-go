package mackerel

// TODO

// Channel represents a Mackerel notification channel.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	// Exists when the type is "email"
	Emails  []string `json:"emails,omitempty"`
	UserIDs []string `json:"userIds,omitempty"`

	// Exists when the type is "slack"
	Mentions struct {
		OK       string `json:"ok,omitempty"`
		Warning  string `json:"warning,omitempty"`
		Critical string `json:"critical,omitempty"`
	} `json:"mentions,omitempty"`
	EnabledGraphImage bool `json:"enabledGraphImage,omitempty"`

	// Exists when the type is "slack" or "webhook"
	URL string `json:"url,omitempty"`

	// Exists when the type is "email", "slack", or "webhook"
	Events []string `json:"events,omitempty"`
}
