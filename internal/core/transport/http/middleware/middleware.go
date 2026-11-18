package core_http_middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		// первый элемент - первые луковое кольцо, поэтому начинаем с конца
		// в конце ...m будет условная middleware Trace(), которая станет первым кольцом
		h = m[i](h) // создаем лук из уровней middleware,
	}

	return h
}
