package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/models"
	"github.com/gin-gonic/gin"
)

// FollowUser godoc
// @Summary Follow a user
// @Description Current user follows another user
// @Tags followers
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param followRequest body models.FollowRequest true "User to follow"
// @Success 201 {object} models.FollowerResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/followers/follow [post]
func FollowUser(c *gin.Context) {
	// Get current user from context (set by auth middleware)
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Can't follow yourself
	if followerID.(uint) == req.FollowedID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// Check if user to follow exists
	var followedUser models.User
	if err := database.DB.First(&followedUser, req.FollowedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}

	// Check if already following
	var existingFollow models.Follower
	err := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, req.FollowedID).First(&existingFollow).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already following this user"})
		return
	}

	// Create follow relationship
	follower := models.Follower{
		FollowerID:    followerID.(uint),
		FollowedID:    req.FollowedID,
		FollowedSince: time.Now(),
	}

	if err := database.DB.Create(&follower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	// Load relations for response
	database.DB.Preload("Follower").Preload("Followed").First(&follower, follower.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Successfully followed user",
		"follower": follower.ToResponse(),
	})
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Current user unfollows another user
// @Tags followers
// @Produce json
// @Security CookieAuth
// @Param id path int true "User ID to unfollow"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/followers/unfollow/{id} [delete]
func UnfollowUser(c *gin.Context) {
	// Get current user from context
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	followedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find follow relationship
	var follower models.Follower
	err = database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follower).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not following this user"})
		return
	}

	// Delete relationship
	if err := database.DB.Delete(&follower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

// GetMyFollowers godoc
// @Summary Get my followers
// @Description Get list of users following the current user
// @Tags followers
// @Produce json
// @Security CookieAuth
// @Success 200 {array} models.FollowerResponse
// @Failure 401 {object} map[string]string
// @Router /api/followers/my-followers [get]
func GetMyFollowers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var followers []models.Follower
	if err := database.DB.
		Preload("Follower").
		Preload("Followed").
		Where("followed_id = ?", userID).
		Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	responses := make([]models.FollowerResponse, len(followers))
	for i, f := range followers {
		responses[i] = f.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}

// GetMyFollowing godoc
// @Summary Get users I follow
// @Description Get list of users the current user follows
// @Tags followers
// @Produce json
// @Security CookieAuth
// @Success 200 {array} models.FollowerResponse
// @Failure 401 {object} map[string]string
// @Router /api/followers/my-following [get]
func GetMyFollowing(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var following []models.Follower
	if err := database.DB.
		Preload("Follower").
		Preload("Followed").
		Where("follower_id = ?", userID).
		Find(&following).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}

	responses := make([]models.FollowerResponse, len(following))
	for i, f := range following {
		responses[i] = f.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}
