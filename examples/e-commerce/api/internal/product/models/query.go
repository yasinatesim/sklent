package productmodels

const (
	DEFAULT_PAGE_SIZE = 12
	MAX_PAGE_SIZE     = 60
)

type ListQuery struct {
	CategorySlug  string
	Q             string
	Page          int
	PageSize      int
	MinPriceCents int64
	MaxPriceCents int64
	Sort          string
	InStockOnly   bool
	MinRating     int
}

// RatingFacetThresholds: "average approved rating >= N" buckets, highest first.
var RatingFacetThresholds = []int{4, 3, 2}

type RatingFacet struct {
	Threshold int   `json:"threshold"`
	Count     int64 `json:"count"`
}

type CategoryFacet struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type Facets struct {
	Categories    []CategoryFacet `json:"categories"`
	Ratings       []RatingFacet   `json:"ratings"`
	MinPriceCents int64           `json:"minPriceCents"`
	MaxPriceCents int64           `json:"maxPriceCents"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}
