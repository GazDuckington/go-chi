package utils

import (
	"log"
)

func StandardResponse(status bool, data map[string]interface{}) map[string]interface{} {
	var stat string
	if status {
		stat = "true"
	} else {
		stat = "false"
	}

	return map[string]interface{}{
		"status": stat,
		"data":   data,
	}
}

func PaginatedResponse(status bool, data []map[string]interface{}, page, itemsPerPage int) map[string]interface{} {
	totalItems := len(data)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage

	if end > totalItems {
		end = totalItems
	}

	paginatedData := data[start:end]

	var stat string
	if status {
		stat = "true"
	} else {
		stat = "false"
	}

	return map[string]interface{}{
		"status": stat,
		"data":   paginatedData,
		"pagination": map[string]interface{}{
			"current_page":   page,
			"total_pages":    totalPages,
			"items_per_page": itemsPerPage,
			"total_items":    totalItems,
		},
	}
}

func CheckDataType(status bool, data interface{}) interface{} {
	// Check if data is a slice of maps
	if ret, ok := data.([]map[string]interface{}); ok {
		return PaginatedResponse(status, ret, 1, 10)
	}

	// Check if data is a single map
	if ret, ok := data.(map[string]interface{}); ok {
		log.Print("here ya dum dum")
		return StandardResponse(status, ret)
	}

	return data
}
