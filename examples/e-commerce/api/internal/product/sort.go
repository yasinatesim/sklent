package product

const (
	PRODUCT_SORT_NEWEST     = "newest"
	PRODUCT_SORT_PRICE_ASC  = "price_asc"
	PRODUCT_SORT_PRICE_DESC = "price_desc"
)

// The sort value reaches an ORDER BY clause, so it is matched against a closed set and never
// interpolated from the request.
var validSorts = map[string]bool{
	PRODUCT_SORT_NEWEST:     true,
	PRODUCT_SORT_PRICE_ASC:  true,
	PRODUCT_SORT_PRICE_DESC: true,
}

func IsValidSort(s string) bool {
	return s == "" || validSorts[s]
}

var orderClauseBySort = map[string]string{
	PRODUCT_SORT_NEWEST:     "created_at DESC",
	PRODUCT_SORT_PRICE_ASC:  "price_cents ASC",
	PRODUCT_SORT_PRICE_DESC: "price_cents DESC",
}

func OrderClause(sort string) string {
	if clause, ok := orderClauseBySort[sort]; ok {
		return clause
	}
	return orderClauseBySort[PRODUCT_SORT_NEWEST]
}
