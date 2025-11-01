package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DaniyarDaniyar/go-practice5/config"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	City        string `db:"city" json:"city"`
	TotalOrders int    `db:"total_orders" json:"total_orders"`
}

func GetUsers(c *gin.Context) {
	start := time.Now()
	db := config.DB

	city := c.Query("city")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	query := `
		SELECT 
			u.id, 
			u.name, 
			u.city, 
			COALESCE(COUNT(o.id), 0) AS total_orders
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
	`

	conditions := []string{}
	args := []interface{}{}

	if city != "" {
		conditions = append(conditions, "u.city = ?")
		args = append(args, city)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " GROUP BY u.id ORDER BY total_orders DESC, u.id DESC"

	if limitStr != "" {
		query += " LIMIT ?"
		if limit, err := strconv.Atoi(limitStr); err == nil {
			args = append(args, limit)
		}
	}

	if offsetStr != "" {
		query += " OFFSET ?"
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			args = append(args, offset)
		}
	}

	var users []UserResponse
	err := db.Select(&users, db.Rebind(query), args...)
	if err != nil {
		log.Printf("Ошибка запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	elapsed := time.Since(start)
	c.Header("X-Query-Time", elapsed.String())
	c.JSON(http.StatusOK, users)
}
