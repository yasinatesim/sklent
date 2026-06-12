package promotion

import (
	"github.com/yasinatesim/vela-commerce/api/internal/constants"
)

type CartLine struct {
	ItemID     string
	ProductID  string
	CategoryID *string
	UnitPrice  int64
	Quantity   int
}

type AppliedDiscount struct {
	Kind        string   `json:"kind"`
	Source      string   `json:"source"`
	TargetLines []string `json:"targetLines"`
	AmountCents int64    `json:"amountCents"`
}

type PricingResult struct {
	Subtotal      int64             `json:"subtotal"`
	Promotions    []AppliedDiscount `json:"promotions,omitempty"`
	Coupon        *AppliedDiscount  `json:"coupon,omitempty"`
	ShippingCents int64             `json:"shippingCents"`
	Total         int64             `json:"total"`
}

type PromotionInput struct {
	ID            string
	DiscountType  string
	DiscountValue int
	ScopeType     string
	ProductIDs    []string
	CategoryIDs   []string
	MinCartCents  int64
}

type CouponInput struct {
	ID            string
	Code          string
	DiscountType  string
	DiscountValue int
	ScopeType     string
	ProductIDs    []string
	CategoryIDs   []string
	MinCartCents  int64
}

// Evaluate is pure: same inputs always produce the same PricingResult. No DB, no clock.
func Evaluate(lines []CartLine, promos []PromotionInput, coupon *CouponInput, shippingCents int64) PricingResult {
	subtotal := int64(0)
	for _, l := range lines {
		subtotal += lineSubtotal(l)
	}

	result := PricingResult{Subtotal: subtotal, ShippingCents: shippingCents}

	for _, p := range promos {
		if subtotal < p.MinCartCents {
			continue
		}
		amount, targets := applyDiscount(lines, p.ScopeType, p.ProductIDs, p.CategoryIDs, p.DiscountType, p.DiscountValue)
		if amount > 0 {
			result.Promotions = append(result.Promotions, AppliedDiscount{
				Kind: p.DiscountType, Source: p.ID, TargetLines: targets, AmountCents: amount,
			})
		}
	}

	if coupon != nil && subtotal >= coupon.MinCartCents {
		amount, targets := applyDiscount(lines, coupon.ScopeType, coupon.ProductIDs, coupon.CategoryIDs, coupon.DiscountType, coupon.DiscountValue)
		if amount > 0 {
			result.Coupon = &AppliedDiscount{
				Kind: coupon.DiscountType, Source: coupon.Code, TargetLines: targets, AmountCents: amount,
			}
		}
	}

	discounted := subtotal
	for _, d := range result.Promotions {
		discounted -= d.AmountCents
	}
	if result.Coupon != nil {
		discounted -= result.Coupon.AmountCents
	}
	if discounted < 0 {
		discounted = 0
	}
	result.Total = discounted + shippingCents
	return result
}

func applyDiscount(lines []CartLine, scope string, productIDs, categoryIDs []string, discountType string, discountValue int) (int64, []string) {
	pset := toSet(productIDs)
	cset := toSet(categoryIDs)
	var total int64
	var targets []string
	for _, l := range lines {
		if !applies(scope, pset, cset, l) {
			continue
		}
		amount := discountFor(l, discountType, discountValue)
		if amount > 0 {
			total += amount
			targets = append(targets, l.ItemID)
		}
	}
	return total, targets
}

func lineSubtotal(l CartLine) int64 { return l.UnitPrice * int64(l.Quantity) }

func applies(scope string, productIDs, categoryIDs map[string]struct{}, line CartLine) bool {
	switch scope {
	case constants.SCOPE_TYPE_ALL:
		return true
	case constants.SCOPE_TYPE_PRODUCTS:
		_, ok := productIDs[line.ProductID]
		return ok
	case constants.SCOPE_TYPE_CATEGORIES:
		if line.CategoryID == nil {
			return false
		}
		_, ok := categoryIDs[*line.CategoryID]
		return ok
	}
	return false
}

func discountFor(line CartLine, discountType string, discountValue int) int64 {
	base := lineSubtotal(line)
	switch discountType {
	case constants.DISCOUNT_TYPE_PERCENT:
		return base * int64(discountValue) / 100
	case constants.DISCOUNT_TYPE_FIXED_TRY:
		if int64(discountValue) > base {
			return base
		}
		return int64(discountValue)
	}
	return 0
}

func toSet(ids []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		m[id] = struct{}{}
	}
	return m
}
