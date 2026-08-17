package product

import (
	"context"

	"gorm.io/gorm"

	productmodels "github.com/yasinatesim/vela-commerce/api/internal/product/models"
)

// filtered builds the shared WHERE for listing and faceting, so a facet count can never disagree
// with the list it describes.
func (r *Repo) filtered(ctx context.Context, q productmodels.ListQuery) *gorm.DB {
	tx := r.db.WithContext(ctx).Model(&productmodels.Product{}).Where("published = ?", true)

	if q.CategorySlug != "" && q.CategorySlug != "all" {
		tx = tx.Where("category_slug = ?", q.CategorySlug)
	}
	if q.Q != "" {
		// unaccent() makes the match Turkish diacritic-insensitive, so "gumus" finds "gümüş";
		// both locale titles and the description are searched, not just the leading title.
		needle := "%" + q.Q + "%"
		tx = tx.Where(
			`(unaccent(title_tr) ILIKE unaccent(?)
				OR unaccent(title_en) ILIKE unaccent(?)
				OR unaccent(description_tr) ILIKE unaccent(?)
				OR unaccent(slug) ILIKE unaccent(?))`,
			needle, needle, needle, needle,
		)
	}
	if q.MinPriceCents > 0 {
		tx = tx.Where("price_cents >= ?", q.MinPriceCents)
	}
	if q.MaxPriceCents > 0 {
		tx = tx.Where("price_cents <= ?", q.MaxPriceCents)
	}
	if q.InStockOnly {
		tx = tx.Where("stock > 0")
	}
	return tx
}

func (r *Repo) List(ctx context.Context, q productmodels.ListQuery) ([]productmodels.Product, productmodels.Pagination, error) {
	var total int64
	if err := r.filtered(ctx, q).Count(&total).Error; err != nil {
		return nil, productmodels.Pagination{}, err
	}

	var out []productmodels.Product
	err := r.filtered(ctx, q).
		Order(OrderClause(q.Sort)).
		Limit(q.PageSize).
		Offset((q.Page - 1) * q.PageSize).
		Find(&out).Error
	if err != nil {
		return nil, productmodels.Pagination{}, err
	}

	return out, productmodels.Pagination{Page: q.Page, PageSize: q.PageSize, Total: total}, nil
}

// Facets counts the same filtered set the list uses, except that the category facet ignores the
// selected category — otherwise picking one category would hide every other option.
func (r *Repo) Facets(ctx context.Context, q productmodels.ListQuery) (productmodels.Facets, error) {
	out := productmodels.Facets{Categories: []productmodels.CategoryFacet{}, Ratings: []productmodels.RatingFacet{}}

	categoryQuery := q
	categoryQuery.CategorySlug = ""
	rows := []struct {
		CategorySlug string
		Count        int64
	}{}
	err := r.filtered(ctx, categoryQuery).
		Select("category_slug, COUNT(*) AS count").
		Group("category_slug").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return productmodels.Facets{}, err
	}
	for _, row := range rows {
		out.Categories = append(out.Categories, productmodels.CategoryFacet{
			Slug: row.CategorySlug, Name: row.CategorySlug, Count: row.Count,
		})
	}

	bounds := struct {
		Min int64
		Max int64
	}{}
	if err := r.filtered(ctx, q).Select("COALESCE(MIN(price_cents),0) AS min, COALESCE(MAX(price_cents),0) AS max").Scan(&bounds).Error; err != nil {
		return productmodels.Facets{}, err
	}
	out.MinPriceCents, out.MaxPriceCents = bounds.Min, bounds.Max

	return out, nil
}
