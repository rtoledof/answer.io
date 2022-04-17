package handler

import (
	"answer.io/pkg/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type response struct {
	Key   model.Key   `json:"key"`
	Value model.Value `json:"value"`
}

func (r *response) Marshal(q model.Question) {
	r.Key = q.Key
	r.Value = q.Value
}

type handler struct {
	manager QuestionManager
}

func NewQuestionHandler(e *echo.Echo, manager QuestionManager) {
	h := &handler{manager: manager}
	g := e.Group("questions")
	g.POST("/", h.post)
	g.POST("", h.post)
	g.GET("/", h.list)
	g.GET("", h.list)
	g.GET("/:key/history", h.history)
	g.GET("/:key", h.get)
	g.PUT("/:key", h.put)
	g.DELETE("/:key", h.delete)
}

func (h *handler) post(c echo.Context) error {
	key := c.FormValue("key")
	value := c.FormValue("value")
	q, err := h.manager.New(model.Key(key), model.Value(value))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusNoContent, response{
		Key:   q.Key,
		Value: q.Value,
	})
}

func (h *handler) put(c echo.Context) error {
	key := c.Param("key")
	newValue := c.FormValue("value")
	if err := h.manager.Update(model.Key(key), model.Value(newValue)); err != nil {
		return echo.ErrBadRequest
	}
	return c.String(http.StatusNoContent, "")
}

func (h *handler) delete(c echo.Context) error {
	key := c.Param("key")
	if err := h.manager.Delete(model.Key(key)); err != nil {
		return echo.ErrNotFound
	}
	return c.String(http.StatusNoContent, "")
}

func (h *handler) get(c echo.Context) error {
	key := c.Param("key")
	q, err := h.manager.Get(model.Key(key))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	rsp := response{
		Key:   q.Key,
		Value: q.Value,
	}
	return c.JSON(http.StatusOK, rsp)
}

func (h *handler) history(c echo.Context) error {
	key := c.Param("key")
	q, err := h.manager.Get(model.Key(key))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var events []struct {
		Event string `json:"event"`
		Data  model.Data
	}
	for i := len(q.History) - 1; i >= 0; i-- {
		var ev = struct {
			Event string `json:"event"`
			Data  model.Data
		}{
			Event: q.History[i].String(),
			Data:  q.History[i].Data(),
		}

		events = append(events, ev)
	}
	return c.JSON(http.StatusOK, events)
}

func (h *handler) list(c echo.Context) error {
	list, err := h.manager.List()
	if err != nil {
		return echo.ErrBadRequest
	}
	var l = make([]response, len(list))
	for i, v := range list {
		l[i].Marshal(v)
	}
	return c.JSON(http.StatusOK, l)
}
