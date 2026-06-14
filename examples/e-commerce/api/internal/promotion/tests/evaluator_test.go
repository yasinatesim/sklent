package promotion_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/promotion"
)

func lines() []promotion.CartLine {
	cat := "cat-1"
	return []promotion.CartLine{
		{ItemID: "i1", ProductID: "p1", CategoryID: &cat, UnitPrice: 10000, Quantity: 2},
		{ItemID: "i2", ProductID: "p2", UnitPrice: 5000, Quantity: 1},
	}
}

func TestEvaluate_NoDiscounts(t *testing.T) {
	r := promotion.Evaluate(lines(), nil, nil, 2000)
	assert.Equal(t, int64(25000), r.Subtotal)
	assert.Equal(t, int64(27000), r.Total)
	assert.Empty(t, r.Promotions)
	assert.Nil(t, r.Coupon)
}

func TestEvaluate_PercentAllScope(t *testing.T) {
	promo := promotion.PromotionInput{
		ID: "promo-10", DiscountType: constants.DISCOUNT_TYPE_PERCENT, DiscountValue: 10,
		ScopeType: constants.SCOPE_TYPE_ALL,
	}
	r := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, nil, 0)
	assert.Equal(t, int64(2500), r.Promotions[0].AmountCents)
	assert.Equal(t, int64(22500), r.Total)
}

func TestEvaluate_FixedCappedAtLine(t *testing.T) {
	promo := promotion.PromotionInput{
		ID: "promo-big", DiscountType: constants.DISCOUNT_TYPE_FIXED_TRY, DiscountValue: 9999999,
		ScopeType: constants.SCOPE_TYPE_PRODUCTS, ProductIDs: []string{"p2"},
	}
	r := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, nil, 0)
	assert.Equal(t, int64(5000), r.Promotions[0].AmountCents)
}

func TestEvaluate_CategoryScopeAndCoupon(t *testing.T) {
	promo := promotion.PromotionInput{
		ID: "promo-cat", DiscountType: constants.DISCOUNT_TYPE_PERCENT, DiscountValue: 50,
		ScopeType: constants.SCOPE_TYPE_CATEGORIES, CategoryIDs: []string{"cat-1"},
	}
	coupon := &promotion.CouponInput{
		Code: "SAVE20", DiscountType: constants.DISCOUNT_TYPE_PERCENT, DiscountValue: 20,
		ScopeType: constants.SCOPE_TYPE_ALL,
	}
	r := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, coupon, 0)
	assert.Equal(t, int64(10000), r.Promotions[0].AmountCents)
	assert.NotNil(t, r.Coupon)
	assert.Equal(t, "SAVE20", r.Coupon.Source)
}

func TestEvaluate_MinCartNotMet(t *testing.T) {
	promo := promotion.PromotionInput{
		ID: "promo-min", DiscountType: constants.DISCOUNT_TYPE_PERCENT, DiscountValue: 10,
		ScopeType: constants.SCOPE_TYPE_ALL, MinCartCents: 999999,
	}
	r := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, nil, 0)
	assert.Empty(t, r.Promotions)
}

func TestEvaluate_Deterministic(t *testing.T) {
	promo := promotion.PromotionInput{ID: "p", DiscountType: constants.DISCOUNT_TYPE_PERCENT, DiscountValue: 15, ScopeType: constants.SCOPE_TYPE_ALL}
	a := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, nil, 1500)
	b := promotion.Evaluate(lines(), []promotion.PromotionInput{promo}, nil, 1500)
	assert.Equal(t, a, b)
}
