package product_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/product"
	productmodels "github.com/yasinatesim/vela-commerce/api/internal/product/models"
)

func TestIsValidSort(t *testing.T) {
	assert.True(t, product.IsValidSort(""), "empty means the default order")
	assert.True(t, product.IsValidSort(product.PRODUCT_SORT_NEWEST))
	assert.True(t, product.IsValidSort(product.PRODUCT_SORT_PRICE_ASC))
	assert.True(t, product.IsValidSort(product.PRODUCT_SORT_PRICE_DESC))
	assert.False(t, product.IsValidSort("price_asc; DROP TABLE products"))
	assert.False(t, product.IsValidSort("random"))
}

func TestParseListQueryDefaults(t *testing.T) {
	q := product.ParseListQuery(url.Values{})

	assert.Equal(t, 1, q.Page)
	assert.Equal(t, productmodels.DEFAULT_PAGE_SIZE, q.PageSize)
	assert.Empty(t, q.Sort)
	assert.False(t, q.InStockOnly)
}

func TestParseListQueryReadsEveryFilter(t *testing.T) {
	q := product.ParseListQuery(url.Values{
		"category":  {"rings"},
		"q":         {"gümüş"},
		"page":      {"3"},
		"pageSize":  {"24"},
		"minPrice":  {"1000"},
		"maxPrice":  {"5000"},
		"sort":      {product.PRODUCT_SORT_PRICE_ASC},
		"inStock":   {"true"},
		"minRating": {"4"},
	})

	assert.Equal(t, "rings", q.CategorySlug)
	assert.Equal(t, "gümüş", q.Q)
	assert.Equal(t, 3, q.Page)
	assert.Equal(t, 24, q.PageSize)
	assert.Equal(t, int64(1000), q.MinPriceCents)
	assert.Equal(t, int64(5000), q.MaxPriceCents)
	assert.Equal(t, product.PRODUCT_SORT_PRICE_ASC, q.Sort)
	assert.True(t, q.InStockOnly)
	assert.Equal(t, 4, q.MinRating)
}

func TestParseListQueryRejectsHostileInput(t *testing.T) {
	q := product.ParseListQuery(url.Values{
		"page":     {"-5"},
		"pageSize": {"100000"},
		"minPrice": {"-1"},
		"sort":     {"price_asc; DROP TABLE products"},
	})

	assert.Equal(t, 1, q.Page, "a negative page must fall back to the first page")
	assert.LessOrEqual(t, q.PageSize, productmodels.MAX_PAGE_SIZE,
		"an unbounded page size is a denial-of-service lever")
	assert.Equal(t, int64(0), q.MinPriceCents)
	assert.Empty(t, q.Sort, "an unknown sort must be dropped, never passed to the query builder")
}

func TestParseListQuerySwapsInvertedPriceBounds(t *testing.T) {
	q := product.ParseListQuery(url.Values{"minPrice": {"9000"}, "maxPrice": {"1000"}})

	assert.Equal(t, int64(1000), q.MinPriceCents)
	assert.Equal(t, int64(9000), q.MaxPriceCents)
}

func TestParseListQueryClampsRatingToTheFacetBuckets(t *testing.T) {
	assert.Equal(t, 0, product.ParseListQuery(url.Values{"minRating": {"9"}}).MinRating)
	assert.Equal(t, 0, product.ParseListQuery(url.Values{"minRating": {"1"}}).MinRating)
	assert.Equal(t, 4, product.ParseListQuery(url.Values{"minRating": {"4"}}).MinRating)
}
