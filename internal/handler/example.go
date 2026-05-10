package handler

import (
	"github.com/gin-gonic/gin"

	_ "go-skeleton/internal/model"
	"go-skeleton/internal/service"
	"go-skeleton/pkg/response"
)

// ExampleHandler handles HTTP requests for examples.
type ExampleHandler struct {
	svc *service.ExampleService
}

// NewExampleHandler creates an ExampleHandler.
func NewExampleHandler(svc *service.ExampleService) *ExampleHandler {
	return &ExampleHandler{svc: svc}
}

// Create handles POST /examples.
// @Summary      Create an example
// @Description  Create a new example
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        body  body      service.CreateExampleReq  true  "Example to create"
// @Success      200   {object}  response.Response{data=model.Example}
// @Router       /examples [post]
func (h *ExampleHandler) Create(c *gin.Context) {
	var req service.CreateExampleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.BuildValidationErrorResponse(c, err))
		return
	}

	example, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.WriteSuccess(c, example)
}

// List handles GET /examples.
// @Summary      List examples
// @Description  Get a paginated list of examples
// @Tags         examples
// @Produce      json
// @Param        limit   query     int  false  "Limit"   minimum(1)  maximum(100)
// @Param        offset  query     int  false  "Offset"  minimum(0)
// @Success      200     {object}  response.Response{data=service.ListExamplesRes}
// @Router       /examples [get]
func (h *ExampleHandler) List(c *gin.Context) {
	var req service.ListExamplesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.BuildValidationErrorResponse(c, err))
		return
	}

	res, err := h.svc.List(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.WriteSuccess(c, res)
}
