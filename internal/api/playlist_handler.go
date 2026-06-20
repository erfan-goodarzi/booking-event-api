package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internal/messages"
	"github.com/erfan-goodarzi/booking-event-api/internal/models"
	"github.com/erfan-goodarzi/booking-event-api/internal/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/pkg/validation"
	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	playlistStore store.PlaylistStore
	logger        *log.Logger
	response      *APIResponse
}

func NewPlaylistHandler(PlaylistStore store.PlaylistStore, logger *log.Logger, response *APIResponse) *PlaylistHandler {
	return &PlaylistHandler{
		PlaylistStore,
		logger,
		response,
	}
}

// GetAll godoc
// @Summary Playlists
// @Description Get all Playlists
// @Tags Playlist
// @Produce json
// @Success 200 {object} models.PlaylistListResponse
// @Failure 500 {object} models.ErrorInternalServer
// @Router /playlist [get]
func (h *PlaylistHandler) GetAll(c *gin.Context) {
	playlists, err := h.playlistStore.GetAll()

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondRetrievedSuccess(c, http.StatusOK, playlists)
}

// GetById godoc
// @Summary Get playlist by ID
// @Description Get an playlist by its ID
// @Tags Playlist
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} models.PlaylistResponse
// @Failure 404 {object} models.ErrorNotFound
// @Failure 500 {object} models.ErrorInternalServer
// @Router /playlist/{id} [get]
func (h *PlaylistHandler) GetById(c *gin.Context) {
	id, err := apiUtils.ParseID(c)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	playlist, err := h.playlistStore.GetById(id)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, err.Error())
		return
	}

	h.response.RespondRetrievedSuccess(c, http.StatusOK, playlist)
}

// CreatePlaylist godoc
// @Summary Create a Playlist
// @Description Create a new Playlist (authenticated)
// @Tags Playlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param playlist body models.CreatePlaylistRequest true "Playlist payload"
// @Success 201 {object} models.PlaylistResponse
// @Failure 422 {object} models.ErrorBadRequest
// @Failure 404 {object} models.ErrorNotFound
// @Failure 500 {object} models.ErrorInternalServer
// @Router /playlist [post]
func (h *PlaylistHandler) Create(c *gin.Context) {
	var payload models.CreatePlaylistRequest

	err := c.ShouldBindJSON(&payload)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	playlist := &models.Playlist{
		Name:   payload.Name,
		Color:  payload.Color,
		UserId: c.GetString("userId"),
	}

	err = validation.Validate.Struct(payload)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	createdPlaylist, err := h.playlistStore.Create(playlist)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "FAILED_TO_REGISTER")
		return
	}

	h.response.RespondSuccess(c, http.StatusCreated, messages.CreateTicketSuccess, createdPlaylist)
}

// Update godoc
// @Summary Update an playlist
// @Description Update an existing playlist (authenticated, owner only)
// @Tags Playlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Playlist ID"
// @Param playlist body models.PatchPlaylistRequest true "Patch payload"
// @Success 200 {object} models.Playlist
// @Failure 401 {object} models.ErrorUnauthorized
// @Failure 403 {object} models.ErrorForbidden
// @Failure 404 {object} models.ErrorNotFound
// @Failure 422 {object} models.ErrorBadRequest
// @Router /playlist/{id} [put]
func (h *PlaylistHandler) Update(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	existingPlaylist, err := h.playlistStore.GetById(id)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "PLAYLIST_NOT_FOUND")
		return
	}

	if existingPlaylist == nil {
		h.response.RespondError(c, http.StatusNotFound, "PLAYLIST_NOT_FOUND")
		return
	}

	var partialPlaylist models.PatchPlaylistRequest

	err = c.ShouldBindJSON(&partialPlaylist)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(partialPlaylist)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	owner, err := h.playlistStore.GetOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PLAYLIST_NOT_EXIST")
		return
	}

	if owner != currentUserId {
		h.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	store.ApplyPlaylistPatch(existingPlaylist, partialPlaylist)

	updatedPlaylist, err := h.playlistStore.Update(existingPlaylist)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondSuccess(c, http.StatusOK, messages.UpdatePlaylistSuccess, updatedPlaylist)
}

// Delete godoc
// @Summary Delete an playlist
// @Description Delete an playlist by ID (authenticated, owner only)
// @Tags Playlist
// @Produce json
// @Security BearerAuth
// @Param id path string true "Playlist ID"
// @Success 200 {object} models.PlaylistDeleteSuccess
// @Failure 401 {object} models.ErrorUnauthorized
// @Failure 403 {object} models.ErrorForbidden
// @Failure 422 {object} models.ErrorBadRequest
// @Router /playlist/{id} [delete]
func (h *PlaylistHandler) Delete(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	owner, err := h.playlistStore.GetOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PLAYLIST_NOT_EXIST")
		return
	}

	if owner != currentUserId {
		h.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	err = h.playlistStore.Delete(id)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondSuccess(c, http.StatusOK, messages.DeletesPlaylistSuccess)
}
