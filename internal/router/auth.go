package router

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"spire/lobby/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	AccountId uint64 `json:"account_id" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func HandleAuthBot(c *gin.Context, ctx *core.Context) {
	var r AuthRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	row := ctx.Db.QueryRow("SELECT 1 FROM accounts WHERE account_id = ?", r.AccountId)
	if err := row.Scan(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": r.AccountId,
	})
	s, err := t.SignedString(ctx.Settings.AuthKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token signing failed"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: s,
	})
}

func HandleAccountBot(c *gin.Context, ctx *core.Context) {
	type Request struct {
		BotId uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		AccountId uint64 `json:"account_id"`
	}

	var r Request
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	row := ctx.Db.QueryRow("SELECT account_id FROM bots WHERE bot_id = ?", r.BotId)
	var accountId uint64
	if err := row.Scan(&accountId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, Response{AccountId: 0})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, Response{AccountId: accountId})
}

func HandleRegisterBot(c *gin.Context, ctx *core.Context) {
	type Request struct {
		BotId uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		AccountId uint64 `json:"account_id"`
	}

	var r Request
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// Register bot
	tx, err := ctx.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		log.Fatalf("Error beginning transaction: %v", err)
		return
	}

	cleanup := func() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		err := tx.Rollback()
		if err != nil {
			log.Fatalf("Error rolling back transaction: %v", err)
		}
	}

	row := tx.QueryRow("SELECT 1 FROM bots WHERE bot_id = ?", r.BotId)
	if err := row.Scan(); !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bot already registered"})
		return
	}

	res, err := tx.Exec("INSERT INTO accounts (account_id) VALUES (?)", r.BotId)
	if err != nil {
		cleanup()
		return
	}
	accountId, err := res.LastInsertId()
	if err != nil {
		cleanup()
		return
	}

	_, err = tx.Exec("INSERT INTO bots (bot_id, account_id) VALUES (?, ?)", r.BotId, accountId)
	if err != nil {
		cleanup()
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}

	c.JSON(http.StatusOK, Response{
		AccountId: uint64(accountId),
	})
}
