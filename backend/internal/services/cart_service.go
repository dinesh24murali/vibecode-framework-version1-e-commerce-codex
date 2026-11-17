package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// ErrInvalidQuantity is returned for a quantity below 1.
var ErrInvalidQuantity = errors.New("quantity must be at least 1")

// CartService manages the authenticated customer cart.
type CartService struct {
	repo     repositories.CartRepository
	products *ProductService
}

// NewCartService returns a CartService.
func NewCartService(repo repositories.CartRepository, products *ProductService) *CartService {
	return &CartService{repo: repo, products: products}
}

// Cart returns the user cart. Totals are computed here from current discounted
// prices; lines whose product is no longer active are left out.
func (s *CartService) Cart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	rows, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ProductID)
	}
	sums, err := s.products.SummariesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	cart := &models.Cart{ID: userID, Items: []models.CartItem{}}
	for _, r := range rows {
		sum, ok := sums[r.ProductID]
		if !ok {
			continue
		}
		line := round2(sum.DiscountedPrice * float64(r.Quantity))
		cart.Items = append(cart.Items, models.CartItem{
			ID: r.ID, Product: sum, Quantity: r.Quantity, UnitPrice: sum.DiscountedPrice, LineTotal: line,
		})
		cart.Subtotal += line
		if r.UpdatedAt.After(cart.UpdatedAt) {
			cart.UpdatedAt = r.UpdatedAt
		}
	}
	cart.Subtotal = round2(cart.Subtotal)
	cart.ItemCount = len(cart.Items)
	if cart.UpdatedAt.IsZero() {
		cart.UpdatedAt = time.Now().UTC()
	}
	return cart, nil
}

// Add puts a product in the cart; an existing line has its quantity incremented.
func (s *CartService) Add(ctx context.Context, userID, productID uuid.UUID, quantity int32) (*models.Cart, error) {
	if quantity < 1 {
		return nil, ErrInvalidQuantity
	}
	sums, err := s.products.SummariesByIDs(ctx, []uuid.UUID{productID})
	if err != nil {
		return nil, err
	}
	if _, ok := sums[productID]; !ok {
		return nil, ErrNotFound
	}
	total := quantity
	if existing, err := s.repo.GetByProduct(ctx, userID, productID); err == nil {
		total += existing.Quantity
	} else if !errors.Is(mapNotFound(err), ErrNotFound) {
		return nil, err
	}
	if err := s.repo.Upsert(ctx, userID, productID, total); err != nil {
		return nil, err
	}
	return s.Cart(ctx, userID)
}

// SetQuantity changes the quantity of one of the user cart lines.
func (s *CartService) SetQuantity(ctx context.Context, userID, itemID uuid.UUID, quantity int32) (*models.Cart, error) {
	if quantity < 1 {
		return nil, ErrInvalidQuantity
	}
	if err := s.repo.SetQuantity(ctx, itemID, userID, quantity); err != nil {
		return nil, mapNotFound(err)
	}
	return s.Cart(ctx, userID)
}

// Remove deletes one of the user cart lines.
func (s *CartService) Remove(ctx context.Context, userID, itemID uuid.UUID) (*models.Cart, error) {
	n, err := s.repo.Delete(ctx, itemID, userID)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrNotFound
	}
	return s.Cart(ctx, userID)
}

// Clear empties the user cart.
func (s *CartService) Clear(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Clear(ctx, userID)
}
