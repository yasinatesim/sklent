package constants

// Frontend paths the API redirects to. Kept locale-neutral and in one place so a
// route rename cannot drift between the two codebases.
const (
	FRONTEND_PATH_CHECKOUT_SUCCESS = "/checkout/success"
	FRONTEND_PATH_CHECKOUT_ERROR   = "/checkout/error"
)
