package router

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"spire/lobby/internal/core"
)

func HandleBotAccount(c *gin.Context, ctx *core.Context) {
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
			c.JSON(http.StatusOK, Response{AccountId: 0})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, Response{AccountId: accountId})
}

func HandleBotRegister(c *gin.Context, ctx *core.Context) {
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

func HandleBotCharacterList(c *gin.Context, ctx *core.Context) {
	type Request struct {
		AccountId uint64 `json:"account_id" binding:"required"`
	}

	type Response struct {
		Characters []uint64 `json:"characters"`
	}

	var r Request
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	rows, err := ctx.Db.Query("SELECT character_id FROM characters WHERE account_id = ?", r.AccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var characters []uint64
	for rows.Next() {
		var characterId uint64
		if err := rows.Scan(&characterId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		characters = append(characters, characterId)
	}

	c.JSON(http.StatusOK, Response{Characters: characters})
}

func HandleBotCharacterCreate(c *gin.Context, ctx *core.Context) {
	type Request struct {
		AccountId     uint64 `json:"account_id" binding:"required"`
		CharacterName string `json:"character_name" binding:"required"`
	}

	type Response struct {
		CharacterId uint64 `json:"character_id"`
	}

	var r Request
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	//TODO
}
