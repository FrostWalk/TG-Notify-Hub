package formatters

type Formatter interface {
	// Format takes whatever string is sent to the api and return a valid markdown
	Format(input string) string
	// Slug identify the plugin, must be equal to the topic slug
	Slug() string
}
