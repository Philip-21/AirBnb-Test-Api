package handlers

import (
	"airbnb/middleware"
	"airbnb/models"
	"airbnb/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	DbRepo *repository.UserRepo
}

func NewUserHandlers(repo *repository.UserRepo) *UserHandlers {
	return &UserHandlers{
		DbRepo: repo,
	}
}

// @Tags		   User
// @Summary		   SignUp user
// @Description    A User signups
// @Success        200 "user created successfully"
// @Param          CreateUser body models.CreateUser true "Create User Request"
// @Router         /user/signup [post]
func (h *UserHandlers) CreateUser(ctx *gin.Context) {
	var req models.CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.UserRole,
	}

	if err := h.DbRepo.CreateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userToken, err := middleware.GenerateUserToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"user_id": user.ID,
		"token":   userToken,
	})
}

// @Tags		   User
// @Summary		   Signin user
// @Description    A User log's in
// @Success        200 "login successful"
// @Param          CreateUser body models.LoginUser true "Create User Request"
// @Router         /user/login [post]
func (h *UserHandlers) LoginUser(ctx *gin.Context) {
	var req models.LoginUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.DbRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	userToken, err := middleware.GenerateUserToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   userToken,
		"user_id": user.ID,
		"role":    user.Role,
	})
}
