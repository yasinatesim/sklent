package product

import (
	"net/url"
	"strconv"

	productmodels "github.com/yasinatesim/vela-commerce/api/internal/product/models"
)

func positiveInt(values url.Values, key string, fallback int) int {
	parsed, err := strconv.Atoi(values.Get(key))
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func nonNegativeCents(values url.Values, key string) int64 {
	parsed, err := strconv.ParseInt(values.Get(key), 10, 64)
	if err != nil || parsed < 0 {
		return 0
	}
	return parsed
}

func clampRating(raw string) int {
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	for _, threshold := range productmodels.RatingFacetThresholds {
		if parsed == threshold {
			return parsed
		}
	}
	return 0
}

// ParseListQuery turns raw query params into a bounded, validated query. Every hostile value is
// dropped here so the repository never has to defend itself.
func ParseListQuery(values url.Values) productmodels.ListQuery {
	q := productmodels.ListQuery{
		CategorySlug:  values.Get("category"),
		Q:             values.Get("q"),
		Page:          positiveInt(values, "page", 1),
		PageSize:      positiveInt(values, "pageSize", productmodels.DEFAULT_PAGE_SIZE),
		MinPriceCents: nonNegativeCents(values, "minPrice"),
		MaxPriceCents: nonNegativeCents(values, "maxPrice"),
		InStockOnly:   values.Get("inStock") == "true",
		MinRating:     clampRating(values.Get("minRating")),
	}
	if q.PageSize > productmodels.MAX_PAGE_SIZE {
		q.PageSize = productmodels.MAX_PAGE_SIZE
	}
	if sort := values.Get("sort"); IsValidSort(sort) {
		q.Sort = sort
	}
	// A reversed range is a slip, not a request for an empty page.
	if q.MaxPriceCents > 0 && q.MinPriceCents > q.MaxPriceCents {
		q.MinPriceCents, q.MaxPriceCents = q.MaxPriceCents, q.MinPriceCents
	}
	return q
}
