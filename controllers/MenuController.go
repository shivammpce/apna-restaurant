package controllers

import (
	"context"
	"net/http"
	"strings"

	repo "apna-restaurant-2.0/db/sqlc"
	"apna-restaurant-2.0/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuController struct {
	db *repo.Queries
}

func NewMenuController(db *repo.Queries) *MenuController {
	return &MenuController{db}
}

func (mc *MenuController) AddMenu(c *gin.Context) {
	var menuReqBody *repo.Menu
	if err := c.ShouldBindJSON(&menuReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddMenuRequest(menuReqBody, mc.db, ""); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	menu := &repo.CreateMenuParams{
		Category:    menuReqBody.Category,
		MenuItemIds: menuReqBody.MenuItemIds,
	}
	insertedMenu, err := mc.db.CreateMenu(context.Background(), *menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting menu in db"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "menu created", "data": insertedMenu})
}

func (mc *MenuController) UpdateMenu(c *gin.Context) {
	var menuReqBody *repo.Menu
	if err := c.ShouldBindJSON(&menuReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddMenuRequest(menuReqBody, mc.db, "update"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	existingMenu, err := mc.db.GetMenuByID(context.Background(), menuReqBody.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if len(strings.TrimSpace(menuReqBody.Category)) == 0 {
		menuReqBody.Category = existingMenu.Category
	}
	if len(menuReqBody.MenuItemIds) == 0 {
		menuReqBody.MenuItemIds = existingMenu.MenuItemIds
	}
	menu := &repo.UpdateMenuParams{
		Category:    menuReqBody.Category,
		MenuItemIds: menuReqBody.MenuItemIds,
		ID:          menuReqBody.ID,
	}
	updatedMenu, err := mc.db.UpdateMenu(context.Background(), *menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menu updated", "updated_menu": updatedMenu})
}

// TODO - Make this public later
func (mc *MenuController) GetAllMenus(c *gin.Context) {
	allMenus, err := mc.db.GetAllMenus(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"menus": allMenus})
}

func (mc *MenuController) GetMenuByID(c *gin.Context) {
	menuId := c.Param("id")
	parsedId, err := uuid.Parse(menuId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menuId"})
		return
	}
	if resp, ok := utils.ValidateMenu(parsedId, mc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	requiredMenu, err := mc.db.GetMenuByID(context.Background(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"menu": requiredMenu})
}

func (mc *MenuController) AddMenuItem(c *gin.Context) {
	var menuitemReqBody *repo.Menuitem
	if err := c.ShouldBindJSON(&menuitemReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateAddMenuitemRequest(menuitemReqBody, mc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	menuitem := &repo.CreateMenuitemParams{
		Name:     menuitemReqBody.Name,
		Price:    menuitemReqBody.Price,
		ImageUrl: menuitemReqBody.ImageUrl,
		MenuID:   menuitemReqBody.MenuID,
	}
	createdMenuitem, err := mc.db.CreateMenuitem(context.Background(), *menuitem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Menuitem added", "data": createdMenuitem})
}

func (mc *MenuController) GetAllMenuItems(c *gin.Context) {
	allMenuitems, err := mc.db.ListMenuitems(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"menusitems": allMenuitems})
}

func (mc *MenuController) GetMenuitemByID(c *gin.Context) {
	menuitemId := c.Param("id")
	parsedId, err := uuid.Parse(menuitemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menuitemId"})
		return
	}
	if resp, ok := utils.ValidateMenuItem(parsedId, mc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
	requiredMenu, err := mc.db.GetMenuitemsById(context.Background(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"menuitem": requiredMenu})
}

func (mc *MenuController) UpdateMenuitem(c *gin.Context) {
	var menuitemReqBody *repo.Menuitem
	if err := c.ShouldBindJSON(&menuitemReqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateUpdateMenuItemRequest(menuitemReqBody, mc.db); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	existingMenuitem, err := mc.db.GetMenuitemsById(context.Background(), menuitemReqBody.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if menuitemReqBody.Name == "" {
		menuitemReqBody.Name = existingMenuitem.Name
	}
	if menuitemReqBody.Price == 0 {
		menuitemReqBody.Price = existingMenuitem.Price
	}
	if menuitemReqBody.ImageUrl == "" {
		menuitemReqBody.ImageUrl = existingMenuitem.ImageUrl
	}
	menuitem := &repo.UpdateMenuitemParams{
		ID:       menuitemReqBody.ID,
		Name:     menuitemReqBody.Name,
		Price:    menuitemReqBody.Price,
		ImageUrl: menuitemReqBody.ImageUrl,
	}
	updatedMenuitem, err := mc.db.UpdateMenuitem(context.Background(), *menuitem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menuitem updated", "data": updatedMenuitem})
}
