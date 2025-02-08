package router

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"spire/lobby/internal/core"
)

func HandleBotAccount(c *gin.Context, x *core.Context) {
	type Request struct {
		BotId uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		AccountId uint64 `json:"account_id"`
	}

	var r Request
	if !check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	row := x.D.QueryRow("SELECT account_id FROM bots WHERE bot_id = ?", r.BotId)
	var accountId uint64
	if err := row.Scan(&accountId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			check(err, c, http.StatusInternalServerError)
			return
		}

		accountId = 0
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
	if !check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	tx, err := ctx.D.Begin()
	if !check(err, c, http.StatusInternalServerError) {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT 1 FROM bots WHERE bot_id = ?", r.BotId)
	if err := row.Scan(); !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bot already registered"})
		return
	}

	res, err := tx.Exec("INSERT INTO accounts (account_id) VALUES (?)", r.BotId)
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	accountId, err := res.LastInsertId()
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	_, err = tx.Exec("INSERT INTO bots (bot_id, account_id) VALUES (?, ?)", r.BotId, accountId)
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	if !check(tx.Commit(), c, http.StatusInternalServerError) {
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
	if !check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	rows, err := ctx.D.Query("SELECT character_id FROM characters WHERE account_id = ?", r.AccountId)
	if !check(err, c, http.StatusInternalServerError) {
		return
	}
	defer rows.Close()

	var characters []uint64
	for rows.Next() {
		var characterId uint64
		if !check(rows.Scan(&characterId), c, http.StatusInternalServerError) {
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
	if !check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	tx, err := ctx.D.Begin()
	if !check(err, c, http.StatusInternalServerError) {
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO characters (account_id, character_name) VALUES (?, ?)", r.AccountId, r.CharacterName)
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	characterId, err := res.LastInsertId()
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	if !check(tx.Commit(), c, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, Response{CharacterId: uint64(characterId)})
}
