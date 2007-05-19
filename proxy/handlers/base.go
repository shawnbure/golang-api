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
	endpointHandlersMap map[string][]EndpointGroupHandler
}

func NewGroupHandler() *groupHandler {
	return &groupHandler{
		endpointHandlersMap: make(map[string][]EndpointGroupHandler),
	}
}

func (g *groupHandler) RegisterEndpoints(r *gin.Engine) {
	for groupRoot, handlersGroups := range g.endpointHandlersMap {
		for _, handlersGroup := range handlersGroups {
			routerGroup := r.Group(groupRoot).Use(handlersGroup.Middlewares...)
			{
				for _, h := range handlersGroup.EndpointHandlers {
					routerGroup.Handle(h.Method, h.Path, h.HandlerFunc)
				}
			}
		}
	}
}

func (g *groupHandler) AddEndpointGroupHandler(endpointHandler EndpointGroupHandler) {
	g.endpointHandlersMap[endpointHandler.Root] = append(g.endpointHandlersMap[endpointHandler.Root], endpointHandler)
}
