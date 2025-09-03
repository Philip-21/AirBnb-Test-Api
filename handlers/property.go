package handlers

import (
	"airbnb/middleware"
	"airbnb/models"
	"airbnb/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PropertyHandlers struct {
	DbRepo *repository.PropertyRepo
}

func NewPropertyHandlers(repo *repository.PropertyRepo) *PropertyHandlers {
	return &PropertyHandlers{
		DbRepo: repo,
	}
}

// @Tags		   Property Owner
// @Summary		   SignUp Property Owner
// @Description    A Property Owner signups
// @Success        200 "property owner successfully"
// @Param          Owner body models.CreatePropertyOwner true "Create Property Owner Request"
// @Router         /property/owner/signup [post]
func (h *PropertyHandlers) CreatePropertyOwner(ctx *gin.Context) {
	var req models.CreatePropertyOwner
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}
	owner := models.PropertyOwner{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.PropertyRole,
	}
	if err := h.DbRepo.CreatePropertyOwner(ctx, &owner); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tokenString, err := middleware.GeneratePropertyOwnerToken(owner.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "property owner created successfully",
		"owner_id": owner.ID,
		"role":     owner.Role,
		"token":    tokenString,
	})
}

// @Tags		   Property Owner
// @Summary		   SignIn Property Owner
// @Description    A Property Owner signs-in
// @Success        200 "property owner successfully"
// @Param          Owner body models.LoginPropertyOwner true "Create Property Owner Request"
// @Router         /property/owner/login [post]
func (h *PropertyHandlers) LoginPropertyOwner(ctx *gin.Context) {
	var req models.LoginPropertyOwner
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	owner, err := h.DbRepo.GetPropertyOwnerByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	tokenString, err := middleware.GeneratePropertyOwnerToken(owner.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "login successful",
		"token":    tokenString,
		"owner_id": owner.ID,
		"role":     owner.Role,
	})
}

// @Tags		   Property Owner
// @Summary		   Create Property
// @Description    A Property Owner creates a property
// @Success        200 "property created  successfully"
// @Param          Owner body models.CreateProperty true "Create Property Request"
// @Router         /property/create [post]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *PropertyHandlers) CreateProperty(ctx *gin.Context) {
	var req models.CreateProperty
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	owner, err := middleware.GetPropertyOwner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	property := models.Property{
		Name:        req.PropertyName,
		Description: req.Description,
		Price:       req.Price,
		OwnerID:     owner.ID,
	}
	if err := h.DbRepo.CreateProperty(ctx, &property); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "property created successfully",
		"property_id": property.ID,
	})
}

// @Tags		   Property Owner
// @Summary		   Get a  Property
// @Description    A Property Owner gets a property details
// @Success        200 {object} models.GetProperty "property created  successfully"
// @Router         /property/{propertyid} [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *PropertyHandlers) GetPropertyByID(ctx *gin.Context) {
	idParam := ctx.Param("propertyid")
	propertyID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
		return
	}
	_, err = middleware.GetPropertyOwner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	property, err := h.DbRepo.GetPropertyByID(ctx, propertyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	response := models.GetProperty{
		PropertyID:   property.ID,
		PropertyName: property.Name,
		Description:  property.Description,
		Price:        property.Price,
		PropertyOwner: models.GetPropertyOwner{
			OwnerID: property.Owner.ID,
			Name:    property.Owner.Name,
			Email:   property.Owner.Email,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// @Tags		   Property Owner
// @Summary		   Get all Property
// @Description    A Property Owner gets all his  properties and its details
// @Success        200 {object} []models.GetProperty
// @Router         /property/owner [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *PropertyHandlers) GetAllProperties(ctx *gin.Context) {
	owner, err := middleware.GetPropertyOwner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	properties, err := h.DbRepo.GetAllProperties(ctx, owner.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response models.GetAllProperties
	for _, prop := range properties {
		response.Properties = append(response.Properties, models.GetProperty{
			PropertyID:   prop.ID,
			PropertyName: prop.Name,
			Description:  prop.Description,
			Price:        prop.Price,
			PropertyOwner: models.GetPropertyOwner{
				OwnerID: prop.Owner.ID,
				Name:    prop.Owner.Name,
				Email:   prop.Owner.Email,
			},
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// @Tags		   Property Owner
// @Summary		   Get all Property
// @Description    A User Owner gets all  properties  available
// @Success        200 {object} []models.GetProperty
// @Router         /property/all [get]
func (h *PropertyHandlers) GetProperties(ctx *gin.Context) {
	properties, err := h.DbRepo.GetProperties(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response models.GetAllProperties
	for _, prop := range properties {
		response.Properties = append(response.Properties, models.GetProperty{
			PropertyID:   prop.ID,
			PropertyName: prop.Name,
			Description:  prop.Description,
			Price:        prop.Price,
			PropertyOwner: models.GetPropertyOwner{
				OwnerID: prop.Owner.ID,
				Name:    prop.Owner.Name,
				Email:   prop.Owner.Email,
			},
		})
	}

	ctx.JSON(http.StatusOK, response)
}
