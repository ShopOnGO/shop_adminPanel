package dlq

var autoRetryCategories = map[string]bool{
	"MESSAGE": true,
	"EMAIL":   true,
	"ORDER":   true,
}

func ShouldRetry(n Notification) bool {
	return autoRetryCategories[n.Category]
}
