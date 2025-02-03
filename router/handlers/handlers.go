package handlers

import (
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type handler struct {
	path     string
	function gin.HandlerFunc
}

type Handlers struct {
	mutex    sync.RWMutex
	handlers []handler
}

func (h *Handlers) Add(path string, handlers ...gin.HandlerFunc) {
	if len(handlers) == 0 {
		return
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	fn := func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, path) {
			for _, handler := range handlers {
				handler(c)

				if c.IsAborted() {
					break
				}
			}

			c.Abort()
			return
		}
	}

	h.handlers = append(h.handlers, handler{
		path:     path,
		function: fn,
	})
	h.sortLocked()
}

func (h *Handlers) sortLocked() {
	sort.Slice(h.handlers, func(i, j int) bool {
		return len(h.handlers[i].path) > len(h.handlers[j].path)
	})
}

func (h *Handlers) Handlers() gin.HandlersChain {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	handlers := make(gin.HandlersChain, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.function)
	}

	return handlers
}
