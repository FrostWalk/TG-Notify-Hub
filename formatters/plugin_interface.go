package formatters

type Formatter interface {
	Format(input string) string
	Slug() string
}
