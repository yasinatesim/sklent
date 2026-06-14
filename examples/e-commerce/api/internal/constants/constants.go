package constants

const (
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

const (
	SCOPE_TYPE_ALL        = "all"
	SCOPE_TYPE_PRODUCTS   = "products"
	SCOPE_TYPE_CATEGORIES = "categories"
)

const (
	DISCOUNT_TYPE_PERCENT   = "percent"
	DISCOUNT_TYPE_FIXED_TRY = "fixed_try"
)

const (
	ORDER_STATUS_PENDING   = "pending"
	ORDER_STATUS_PAID      = "paid"
	ORDER_STATUS_SHIPPED   = "shipped"
	ORDER_STATUS_CANCELLED = "cancelled"
)

const (
	PAYMENT_METHOD_CARD          = "card"
	PAYMENT_METHOD_BANK_TRANSFER = "bank_transfer"
)

const (
	RESERVATION_TTL_MINUTES = 15
)

const (
	CSRF_HEADER = "X-CSRF-Token"
	CSRF_COOKIE = "csrf_token"
)
