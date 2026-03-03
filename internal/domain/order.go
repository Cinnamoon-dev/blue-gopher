package domain

import "errors"

type Order struct {
	ID            int64
	UserID        int64
	TotalSpent    int64
	OrderQuantity int64
}

func NewOrder(id int64, userId int64, totalSpent int64, orderQuantity int64) (*Order, error) {
	order := &Order{
		ID:            id,
		UserID:        userId,
		TotalSpent:    totalSpent,
		OrderQuantity: orderQuantity,
	}

	if err := order.ValidateOrderQuantity(); err != nil {
		return nil, err
	}

	if err := order.ValidateTotalSpent(); err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) ValidateTotalSpent() error {
	if o.TotalSpent < 0 {
		return errors.New("Validate Total Spent: total value spent should be positive")
	}

	return nil
}

func (o *Order) ValidateOrderQuantity() error {
	if o.OrderQuantity < 0 {
		return errors.New("Validate Order Quantity: order quantity should be positive")
	}

	return nil
}
