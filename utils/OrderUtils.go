package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	repo "apna-restaurant-2.0/db/sqlc"
	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID         `json:"id"`
	TableID     uuid.UUID         `json:"table_id"`
	Amount      int               `json:"amount"`
	OrderItems  map[uuid.UUID]int `json:"order_items"`
	CreatedAt   time.Time         `json:"created_at"`
	DeliveredAt time.Time         `json:"delivered_at"`
}

func ValidateOrderID(id uuid.UUID, db *repo.Queries) (string, bool) {
	orderItemCount, err := db.CheckExistingOrder(context.Background(), id)
	if err != nil {
		return "Internal server error", false
	}
	if orderItemCount != 1 {
		return "One of the Orderitems does not exist", false
	}
	return "", true
}

func ValidateAddTableRequest(table *repo.Table, db *repo.Queries, flag string) (string, bool) {
	if table.TableNumber <= 0 && flag != "update" {
		return "table num required", false
	} else if !IsValidUUID(table.ID) {
		return "Missing/Invalid id", false
	}
	if (len(table.OrderIds)) > 0 {
		for _, orderId := range table.OrderIds {
			if resp, ok := ValidateOrderID(orderId, db); !ok {
				return resp, false
			}
		}
	}
	return "requirement passed", true
}

func ValidateTableId(tableId uuid.UUID, db *repo.Queries) (string, bool) {
	tableCount, err := db.CheckExistingTable(context.Background(), tableId)
	if err != nil {
		return "Internal server error", false
	}
	if tableCount != 1 {
		return "Table does not exist", false
	}
	return "", true
}

func ValidateAddOrderRequest(order *Order, db *repo.Queries, flag string) (string, bool) {
	applyFlag := flag != ""
	if !IsValidUUID(order.ID) && applyFlag {
		return "Missing/Invalid id", false
	} else if !IsValidUUID(order.TableID) && !applyFlag {
		return "Missing/Invalid tableId", false
	} else if resp, ok := ValidateTableId(order.TableID, db); !ok && !applyFlag {
		return resp, false
	} else if order.Amount < 0 {
		return "Amount should be positive", false
	} else if order.CreatedAt.IsZero() && !applyFlag {
		return "Created at must be a valid time", false
	} else if order.CreatedAt.After(order.DeliveredAt) && applyFlag {
		return "Created time must be earlier than delivered time", false
	}

	for orderItem, quantity := range order.OrderItems {
		if resp, ok := ValidateMenuItem(orderItem, db); !ok {
			return resp, false
		}
		if quantity < 0 {
			return "One of orderitem's quantity is negative", false
		}
	}
	return "requirement passed", true
}

func CalculateOrderAmount(order *Order, db *repo.Queries) (sql.NullInt32, string) {
	amount := 0
	for orderItem, quantity := range order.OrderItems {
		existingOrderItem, err := db.GetMenuitemsById(context.Background(), orderItem)
		if err != nil {
			return sql.NullInt32{
				Int32: int32(0),
				Valid: true,
			}, "Internal server error"
		}
		amount += int(existingOrderItem.Price) * int(quantity)
	}
	return sql.NullInt32{
		Int32: int32(amount),
		Valid: true,
	}, ""
}
func MarshalToOrder(data json.RawMessage) (map[uuid.UUID]int, error) {
	orderItems := make(map[uuid.UUID]int)
	if err := json.Unmarshal(data, &orderItems); err != nil {
		return nil, err
	}
	return orderItems, nil
}
