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

const (
	REVIEW_STATUS_PENDING  = "pending"
	REVIEW_STATUS_APPROVED = "approved"
	REVIEW_STATUS_REJECTED = "rejected"
)

const (
	RETURN_STATUS_REQUESTED = "requested"
	RETURN_STATUS_APPROVED  = "approved"
	RETURN_STATUS_REJECTED  = "rejected"
	RETURN_STATUS_REFUNDED  = "refunded"
)

const (
	RETURN_REASON_DAMAGED      = "damaged"
	RETURN_REASON_WRONG_ITEM   = "wrong_item"
	RETURN_REASON_NOT_AS_SHOWN = "not_as_shown"
	RETURN_REASON_CHANGED_MIND = "changed_mind"
	RETURN_REASON_OTHER        = "other"
)

// ORDER_SOURCE_* names the channel an order arrived from; notification rules key off it.
const (
	ORDER_SOURCE_SITE = "site"
	ORDER_SOURCE_HB   = "HB"
	ORDER_SOURCE_TY   = "TY"
)
