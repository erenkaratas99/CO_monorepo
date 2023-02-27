package genericHandler

import (
	sharedentities "CustomerOrderMonoRepo/shared/entities"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericRepository"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericService"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	repo    *genericRepository.Repository
	service *genericService.Service
}

func NewGenericHandler(repo *genericRepository.Repository, service *genericService.Service) *Handler {
	return &Handler{repo: repo, service: service}
}

// GenericGet godoc
// @Summary      Generic endpoint
// @Description  Gets a generic content by usage
// @Tags         shared
// @Accept       json
// @Produce      json
// @Success      200  {object}  sharedentities.ResponseModel
// @Failure      500      {error}  error         "binding error"
// @Failure      400      {error}  error         "bad request"
// @Router       /.../generic [get]
func (h *Handler) GenericGet(c echo.Context) error {
	jsonMap := sharedentities.JsonMap{}
	err := c.Bind(&jsonMap)
	if err != nil {
		return customErrors.BindErr
	}
	resp, err := h.service.GenericGetService(jsonMap)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
