package handlers

import "github.com/gin-gonic/gin"

type EndpointGroupHandler struct {
	Root             string
	Middlewares      []gin.HandlerFunc
	EndpointHandlers []EndpointHandler
}

type EndpointHandler struct {
	Path        string
	Method      string
	HandlerFunc gin.HandlerFunc
}

type groupHandler struct {
	endpointHandlersMap map[string]EndpointGroupHandler
}

func NewGroupHandler() *groupHandler {
	return &groupHandler{
		endpointHandlersMap: make(map[string]EndpointGroupHandler),
	}
}

func (g *groupHandler) RegisterEndpoints(r *gin.Engine) {
	for groupRoot, handlersGroup := range g.endpointHandlersMap {
		routerGroup := r.Group(groupRoot).Use(handlersGroup.Middlewares...)
		{
			for _, h := range handlersGroup.EndpointHandlers {
				routerGroup.Handle(h.Method, h.Path, h.HandlerFunc)
			}
		}
	}
}

func (g *groupHandler) AddEndpointGroupHandler(endpointHandler EndpointGroupHandler) {
	g.endpointHandlersMap[endpointHandler.Root] = endpointHandler
}

type ApiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func JsonResponse(c *gin.Context, status int, data interface{}, error string) {
	c.JSON(status, ApiResponse{
		Data:  data,
		Error: error,
	})
}
