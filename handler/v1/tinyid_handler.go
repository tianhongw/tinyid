package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tianhongw/tinyid/pkg/errdef"
)

type TinyIdHandler struct {
	*HandlerV1
}

func NewTinyIdHandler(handler *HandlerV1) *TinyIdHandler {
	return &TinyIdHandler{
		HandlerV1: handler,
	}
}

type NextIdResponse struct {
	IDList []int64 `json:"id_list"`
}

const (
	defaultMaxSize = 1000
)

func (h *TinyIdHandler) NextId(c *gin.Context) {
	bizType := c.Query("type")
	if bizType == "" {
		h.ResponseWithError(c, errdef.NewHttpError("empty type", errdef.HttpErrorCodeBadRequest))
		return
	}

	sizeStr := c.Query("size")
	if sizeStr == "" {
		h.ResponseWithError(c, errdef.NewHttpError("empty size", errdef.HttpErrorCodeBadRequest))
		return
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		h.ResponseWithError(c, errdef.HttpErrorBadRequest)
		return
	}

	if size <= 0 || size > defaultMaxSize {
		h.ResponseWithError(c, errdef.NewHttpError(fmt.Sprintf("size should be in [1, %d]", defaultMaxSize), errdef.HttpErrorCodeBadRequest))
		return
	}

	g, err := h.Services.TinyId.GetGenerator(bizType)
	if err != nil {
		h.ResponseWithError(c, err)
		return
	}

	var idList []int64
	if size == 1 {
		id, err := g.NextID()
		if err != nil {
			h.ResponseWithError(c, err)
			return
		}
		idList = append(idList, id)
	} else {
		ids, err := g.NextBatchIDs(int(size))
		if err != nil {
			h.ResponseWithError(c, err)
			return
		}
		idList = append(idList, ids...)
	}

	c.JSON(http.StatusOK, NextIdResponse{IDList: idList})
}
