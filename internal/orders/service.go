package orders

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

type CreateOrderRequest struct {
	TaxID         *string             `json:"tax_id"`
	PaymentMethod string              `json:"payment_method"`
	Items         []CreateOrderItemReq `json:"items"`
}

type CreateOrderItemReq struct {
	ProductID string   `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Modifiers []string `json:"modifiers"` // list of modifier IDs
}

type Service interface {
	CreateOrder(ctx context.Context, cashierID string, req CreateOrderRequest) (*Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, cashierID string, req CreateOrderRequest) (*Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	if req.PaymentMethod == "" {
		return nil, errors.New("payment method is required")
	}

	orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixMilli())

	order := &Order{
		OrderNumber:   orderNumber,
		CashierID:     cashierID,
		TaxID:         req.TaxID,
		PaymentMethod: req.PaymentMethod,
		Status:        "PAID", // Defaulting to PAID as per requirements, could be dynamic
	}

	var subtotal int64
	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}

		productPrice, err := s.repo.GetProductPrice(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product price for %s: %w", itemReq.ProductID, err)
		}

		itemSubtotal := productPrice
		var modifiers []OrderItemModifier
		for _, modID := range itemReq.Modifiers {
			modPrice, err := s.repo.GetModifierPrice(ctx, modID)
			if err != nil {
				return nil, fmt.Errorf("failed to get modifier price for %s: %w", modID, err)
			}
			modifiers = append(modifiers, OrderItemModifier{
				ModifierID: modID,
				Price:      modPrice,
			})
			itemSubtotal += modPrice
		}

		itemSubtotal *= int64(itemReq.Quantity)

		order.Items = append(order.Items, OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: productPrice,
			Subtotal:  itemSubtotal,
			Modifiers: modifiers,
		})

		subtotal += itemSubtotal
	}

	order.Subtotal = subtotal

	var taxAmount int64
	if req.TaxID != nil && *req.TaxID != "" {
		taxRate, err := s.repo.GetTaxRate(ctx, *req.TaxID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tax rate for %s: %w", *req.TaxID, err)
		}
		// calculate tax: subtotal * (rate / 100)
		taxVal := float64(subtotal) * (taxRate / 100.0)
		taxAmount = int64(math.Round(taxVal))
	}
	
	order.TaxAmount = taxAmount
	order.TotalAmount = subtotal + taxAmount

	err := s.repo.CreateOrderTx(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
