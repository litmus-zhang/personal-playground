package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

type createTodoReq struct {
	Title string `json:"title"`
}

func (s *Server) createTodo(c *gin.Context) {
	var param createTodoReq
	if err := c.ShouldBindJSON(&param); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// t, err := s.store.CreateTodo(
	// 	c,
	// 	param.Title,
	// )
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, t)

}

func (s *Server) deleteTodo(c *gin.Context) {
	var req updateTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// err := s.store.DeleteTodoByID(c, req.ID)
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "todo deleted",
	// })

}

type updateTodoReq struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) updateTodo(c *gin.Context) {

	var req updateTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var param createTodoReq
	if err := c.ShouldBindJSON(&param); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// args := db.UpdateTodoParams{
	// 	ID:    req.ID,
	// 	Title: param.Title,
	// }
	// t, err := s.store.UpdateTodo(c, args)
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, t)

}

func (s *Server) completeTodo(c *gin.Context) {
	var req updateTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var param struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// args := db.CompleteTodoParams{
	// 	ID:        req.ID,
	// 	Completed: param.Completed,
	// }
	// t, err := s.store.CompleteTodo(c, args)
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, t)
}

type listAccountRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5,max=100"`
}

func (s *Server) getAllTodos(c *gin.Context) {
	var req listAccountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// args := db.GetAllTodoParams{
	// 	Limit:  req.Size,
	// 	Offset: (req.Page - 1) * req.Size,
	// }
	// todos, err := s.store.GetAllTodo(c, args)
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// response := util.Response{
	// 	Data:    todos,
	// 	Message: "All todos fetched successfully",
	// }
	// response.SetMeta(req.Page, req.Size, int32(len(todos)), c)
	// c.JSON(http.StatusOK, response)

}
func (s *Server) getTodo(c *gin.Context) {
	var req updateTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// t, err := s.store.GetTodoByID(c, req.ID)
	// if err != nil {
	// 	errResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, t)

}
