package main

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

// для санитайзинга на сторое фронта используйте https://github.com/cure53/DOMPurify

func main() {
	sanitizer := bluemonday.UGCPolicy()

	comment := `<a onclick="alert(document.сookie)" href="https://www.mail.ru">Mail.ru</a>`
	fmt.Printf("comment before sanitizing: %s\n", comment)

	comment = sanitizer.Sanitize(comment)
	fmt.Printf("comment after sanitizing: %s\n", comment)
}
