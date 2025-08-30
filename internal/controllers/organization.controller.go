package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/models"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/response"
)

func CreateOrganization(c *gin.Context) {
	var organization models.Organization

	if err := c.BindJSON(&organization); err != nil {
		c.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
		return
	}

	existingOrg, _ := repository.FindOrganizationByDomanin(organization.Domain)
	if (existingOrg != models.Organization{}) {
		c.IndentedJSON(http.StatusBadRequest, response.Error(errors.New("domain already exists"), http.StatusBadRequest, "Domain is taken"))
		return
	}

	rowAffected, err := repository.CreateOrganization(organization)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Database error while creating organization"))
		return
	}

	if rowAffected == 0 {
		c.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Creation failed"))
		return
	}

	c.IndentedJSON(http.StatusOK, response.Success(rowAffected, http.StatusOK, "Organization created"))
}