package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message  string                 `json:"message"`
	Data     interface{}            `json:"data,omitempty"`
	MetaData map[string]interface{} `json:"meta,omitempty"`
}

/*
	SetMeta sets the meta data of the response

It takes in the page number, page size and total count of the data being returned
It sets the meta data of the response
It returns nothing
*/
func (r *Response) SetMeta(page, size, total int32, c *gin.Context) {
	r.MetaData = map[string]interface{}{
		"page":     page,
		"size":     size,
		"total":    total,
		"nextLink": c.Request.URL.String() + "&page=" + strconv.Itoa(int(page+1)),
		"prevLink": c.Request.URL.String() + "&page=" + strconv.Itoa(int(page-1)),
	}
}
