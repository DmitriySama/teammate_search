package teammatesearchapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/DmitriySama/teammate_search/api/swagger"
	tsService "github.com/DmitriySama/teammate_search/internal/services/teammateSearchService"
)

type API struct {
	service     *tsService.Service
	serviceName string
	metrics     *metrics.Collector
	once        sync.Once
	swaggerSpec []byte
}

func New(service *tsService.Service, serviceName string, collector *metrics.Collector) *API {
	return &API{service: service, serviceName: serviceName, metrics: collector}
}

func (a *API) Router() http.Handler {
	router := chi.NewRouter()
	router.Get("/health", a.health)
	router.Get("/metrics", a.metrics.Handler())
	router.Get("/orders/{orderId}/groups", a.getGroups)
	router.Get("/orders/{orderId}/groups/{groupId}/tasks", a.getGroupTasks)
	router.MethodFunc(http.MethodPut, "/orders/{orderId}/tasks/{taskId}/status", a.updateTaskStatus)
	router.MethodFunc(http.MethodPatch, "/orders/{orderId}/tasks/{taskId}/status", a.updateTaskStatus)
	router.Get("/swagger", a.swaggerUI)
	router.Get("/swagger/registry.swagger.json", a.swaggerSpecHandler)
	return router
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	body := map[string]string{
		"service": a.serviceName,
		"status":  "ok",
	}
	writeJSON(w, http.StatusOK, body)
}

func (a *API) getGroups(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	groups, err := a.service.ListGroups(r.Context(), orderID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, groups)
}

func (a *API) getGroupTasks(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	groupID := chi.URLParam(r, "groupId")

	groupName, tasks, err := a.service.ListGroupTasks(r.Context(), orderID, groupID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"order_id":   orderID,
		"group_id":   groupID,
		"group_name": groupName,
		"tasks":      tasks,
	}
	writeJSON(w, http.StatusOK, response)
}

type updateTaskStatusRequest struct {
	Status string `json:"status"`
}

func (a *API) updateTaskStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	taskID := chi.URLParam(r, "taskId")

	var req updateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Status) == "" {
		writeJSONError(w, http.StatusBadRequest, "status is required")
		return
	}

	err := a.service.UpdateTaskStatus(r.Context(), orderID, taskID, req.Status)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeJSONError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidStatus) {
			writeJSONError(w, http.StatusBadRequest, "invalid status")
			return
		}
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload := map[string]interface{}{
		"task_id":  taskID,
		"order_id": orderID,
		"status":   req.Status,
	}
	writeJSON(w, http.StatusOK, payload)
}

func (a *API) swaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
  <title>Task Registry API</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({ url: '/swagger/registry.swagger.json', dom_id: '#swagger-ui' });
    };
  </script>
</body>
</html>`)
}

func (a *API) swaggerSpecHandler(w http.ResponseWriter, _ *http.Request) {
	a.once.Do(func() {
		a.swaggerSpec = swagger.Registry()
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(a.swaggerSpec)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"detail": message})
}
