package utils

import (
	"context"
	"strings"

	repo "apna-restaurant-2.0/db/sqlc"
	"github.com/google/uuid"
)

func ValidateMenuItem(item uuid.UUID, db *repo.Queries) (string, bool) {
	menuitemCount, err := db.CheckExistingMenuitem(context.Background(), item)
	if err != nil {
		return "Internal server error", false
	}
	if menuitemCount != 1 {
		return "One of the Menuitems does not exist", false
	}
	return "", true
}

func ValidateAddMenuRequest(menu *repo.Menu, db *repo.Queries, flag string) (string, bool) {
	if len(strings.TrimSpace(menu.Category)) == 0 && flag != "update" {
		return "Category required", false
	} else if len(strings.TrimSpace(menu.Category)) < 5 && flag != "update" {
		return "Name should be at least 5 chars", false
	} else if !IsValidUUID(menu.ID) {
		return "Missing/Invalid id", false
	} else if resp, ok := ValidateMenu(menu.ID, db); !ok {
		return resp, false
	}
	if (len(menu.MenuItemIds)) > 0 {
		for _, menuitem := range menu.MenuItemIds {
			if resp, ok := ValidateMenuItem(menuitem, db); !ok {
				return resp, false
			}
		}
	}
	return "requirement passed", true
}

func ValidateMenu(menuId uuid.UUID, db *repo.Queries) (string, bool) {
	menuCount, err := db.CheckExistingMenu(context.Background(), menuId)
	if err != nil {
		return "Internal server error", false
	}
	if menuCount != 1 {
		return "Menuid does not exist", false
	}
	return "", true
}

func ValidateAddMenuitemRequest(menuitem *repo.Menuitem, db *repo.Queries) (string, bool) {
	if len(strings.TrimSpace(menuitem.Name)) == 0 {
		return "Name required", false
	} else if len(strings.TrimSpace(menuitem.Name)) < 5 {
		return "Name should be at least 5 chars", false
	} else if menuitem.Price < 0 {
		return "Invalid price", false
	} else if len(strings.TrimSpace(menuitem.ImageUrl)) == 0 {
		return "Image url required", false
	} else if !IsUrlValid(menuitem.ImageUrl) {
		return "Invalid image url", false
	} else if !IsValidUUID(menuitem.MenuID) {
		return "Invalid menuid", false
	} else if resp, ok := ValidateMenu(menuitem.MenuID, db); !ok {
		return resp, false
	}
	return "requirement passed", true
}

func ValidateUpdateMenuItemRequest(menuitem *repo.Menuitem, db *repo.Queries) (string, bool) {
	if len(strings.TrimSpace(menuitem.ID.String())) == 0 {
		return "menuitem id missing", false
	} else if !IsValidUUID(menuitem.ID) {
		return "Invalid menuitem id", false
	} else if resp, ok := ValidateMenuItem(menuitem.ID, db); !ok {
		return resp, false
	} else if len(strings.TrimSpace(menuitem.Name)) != 0 && len(strings.TrimSpace(menuitem.Name)) < 5 {
		return "Name should be at least 5 chars", false
	} else if menuitem.Price != 0 && menuitem.Price < 0 {
		return "Invalid price", false
	} else if len(strings.TrimSpace(menuitem.ImageUrl)) != 0 && !IsUrlValid(menuitem.ImageUrl) {
		return "Invalid image url", false
	}
	return "requirement passed", true
}
