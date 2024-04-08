package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	repo "apna-restaurant-2.0/db/sqlc"
	"apna-restaurant-2.0/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type OrderController struct {
	db *repo.Queries
}

func NewOrderController(db *repo.Queries) *OrderController {
	return &OrderController{db: db}
}

func (oc *OrderController) AddTable(c *gin.Context) {
	var tableReqBody *repo.Table
	if err := c.ShouldBindJSON(&tableReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddTableRequest(tableReqBody, oc.db, ""); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	order := &repo.CreateTableParams{
		TableNumber: tableReqBody.TableNumber,
		OrderIds:    tableReqBody.OrderIds,
	}
	insertedTable, err := oc.db.CreateTable(context.Background(), *order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting table in db"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "table created", "data": insertedTable})
}

func (oc *OrderController) UpdateTable(c *gin.Context) {
	var tableReqBody *repo.Table
	if err := c.ShouldBindJSON(&tableReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	_, err := oc.db.GetTableByID(context.Background(), tableReqBody.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid/missing tableId"})
		return
	}
	if resp, ok := utils.ValidateAddTableRequest(tableReqBody, oc.db, "update"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	existingTable, err := oc.db.GetTableByID(context.Background(), tableReqBody.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if tableReqBody.TableNumber == 0 {
		tableReqBody.TableNumber = existingTable.TableNumber
	}
	if len(tableReqBody.OrderIds) == 0 {
		tableReqBody.OrderIds = existingTable.OrderIds
	}
	table := &repo.UpdateTableParams{
		TableNumber: tableReqBody.TableNumber,
		OrderIds:    tableReqBody.OrderIds,
		ID:          tableReqBody.ID,
	}
	updatedTable, err := oc.db.UpdateTable(context.Background(), *table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table updated", "updated_table": updatedTable})
}

func (oc *OrderController) GetAllTables(c *gin.Context) {
	allTables, err := oc.db.GetAllTables(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tables": allTables})
}

func (oc *OrderController) AddOrder(c *gin.Context) {
	var orderReqBody *utils.Order
	if err := c.ShouldBindJSON(&orderReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddOrderRequest(orderReqBody, oc.db, ""); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	amount, msg := utils.CalculateOrderAmount(orderReqBody, oc.db)
	if msg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	orderItemsJSON, err := json.Marshal(orderReqBody.OrderItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in marshalling orderitems"})
		return
	}
	order := &repo.CreateOrderParams{
		TableID: orderReqBody.TableID,
		Amount:  amount,
		OrderItems: pqtype.NullRawMessage{
			RawMessage: orderItemsJSON,
			Valid:      true,
		},
	}
	insertedOrder, err := oc.db.CreateOrder(context.Background(), *order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	insertedOrderItems, err := utils.MarshalToOrder(insertedOrder.OrderItems.RawMessage)

	resultantOrder := utils.Order{
		ID:          insertedOrder.ID,
		TableID:     insertedOrder.TableID,
		Amount:      int(insertedOrder.Amount.Int32),
		OrderItems:  insertedOrderItems,
		CreatedAt:   insertedOrder.CreatedAt,
		DeliveredAt: insertedOrder.DeliveredAt.Time,
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order created", "data": resultantOrder})
}

func (oc *OrderController) UpdateOrder(c *gin.Context) {
	var orderReqBody *utils.Order
	if err := c.ShouldBindJSON(&orderReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddOrderRequest(orderReqBody, oc.db, "update"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}

	existingOrder, err := oc.db.GetOrderByID(context.Background(), orderReqBody.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if orderReqBody.TableID != (uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000}) && existingOrder.TableID != orderReqBody.TableID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table id cannot be changed"})
		return
	}
	amount, msg := utils.CalculateOrderAmount(orderReqBody, oc.db)
	if msg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if orderReqBody.CreatedAt.IsZero() {
		orderReqBody.CreatedAt = existingOrder.CreatedAt
	}
	orderItemsJSON, err := json.Marshal(orderReqBody.OrderItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in marshalling orderitems"})
		return
	}
	table := &repo.UpdateOrderParams{
		TableID: orderReqBody.TableID,
		Amount:  amount,
		OrderItems: pqtype.NullRawMessage{
			RawMessage: orderItemsJSON,
			Valid:      true,
		},
		CreatedAt: orderReqBody.CreatedAt,
		DeliveredAt: sql.NullTime{
			Time:  orderReqBody.DeliveredAt,
			Valid: true,
		},
		ID: orderReqBody.ID,
	}
	updatedOrder, err := oc.db.UpdateOrder(context.Background(), *table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	resultantOrder := utils.Order{
		ID:          updatedOrder.ID,
		TableID:     updatedOrder.TableID,
		Amount:      int(updatedOrder.Amount.Int32),
		OrderItems:  orderReqBody.OrderItems,
		CreatedAt:   updatedOrder.CreatedAt,
		DeliveredAt: updatedOrder.DeliveredAt.Time,
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated", "updated_order": resultantOrder})
}

func (oc *OrderController) GetOrderDetails(c *gin.Context) {
	orderId := c.Param("id")
	parsedId, err := uuid.Parse(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderId"})
		return
	}
	if resp, ok := utils.ValidateOrderID(parsedId, oc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	requiredOrder, err := oc.db.GetOrderByID(context.Background(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	orderItems, err := utils.MarshalToOrder(requiredOrder.OrderItems.RawMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	resultant := utils.Order{
		ID:          requiredOrder.ID,
		TableID:     requiredOrder.TableID,
		Amount:      int(requiredOrder.Amount.Int32),
		OrderItems:  orderItems,
		CreatedAt:   requiredOrder.CreatedAt,
		DeliveredAt: requiredOrder.DeliveredAt.Time,
	}
	c.JSON(http.StatusOK, gin.H{"order": resultant})
}

func (oc *OrderController) CancelOrder(c *gin.Context) {
	orderId := c.Param("id")
	parsedOrderId, err := uuid.Parse(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tableId"})
		return
	}
	err = oc.db.DeleteOrderByID(context.Background(), parsedOrderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in deleting order from db"})
		return
	}
	err = oc.db.RemoveOrderIDFromTables(context.Background(), parsedOrderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in deleting order from db"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

func (oc *OrderController) GetOrderDetailsForTable(c *gin.Context) {
	tableId := c.Param("table_id")
	parsedTableId, err := uuid.Parse(tableId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tableId"})
		return
	}
	if resp, ok := utils.ValidateTableId(parsedTableId, oc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	requiredOrders := make([]utils.Order, 0)

	existingTable, err := oc.db.GetTableByID(context.Background(), parsedTableId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	orders, err := oc.db.GetAllOrders(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	for _, orderId := range existingTable.OrderIds {

		for _, order := range orders {
			if order.ID == orderId {
				orderItems, err := utils.MarshalToOrder(order.OrderItems.RawMessage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
					return
				}
				requiredOrders = append(requiredOrders, utils.Order{
					ID:          orderId,
					TableID:     existingTable.ID,
					Amount:      int(order.Amount.Int32),
					OrderItems:  orderItems,
					CreatedAt:   order.CreatedAt,
					DeliveredAt: order.DeliveredAt.Time,
				})
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Orders fetched", "data": requiredOrders})
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	allOrders, err := oc.db.GetAllOrders(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	requiredOrders := make([]utils.Order, 0)
	for _, order := range allOrders {
		orderItems, err := utils.MarshalToOrder(order.OrderItems.RawMessage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
			return
		}
		requiredOrders = append(requiredOrders, utils.Order{
			ID:          order.ID,
			TableID:     order.TableID,
			Amount:      int(order.Amount.Int32),
			OrderItems:  orderItems,
			CreatedAt:   order.CreatedAt,
			DeliveredAt: order.DeliveredAt.Time,
		})
	}
	c.JSON(http.StatusOK, gin.H{"orders": requiredOrders})
}
