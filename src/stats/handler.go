package stats

import (
	"encoding/json"
	"net/http"
)

// ApplyRoutes - Apply the routes
func ApplyRoutes(mux *http.ServeMux, settings *Settings) *http.ServeMux {
	service := NewService(settings)
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusSeeOther)
	})
	mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.json")
	})
	mux.HandleFunc("GET /target/status/{targetName}", TargetStatusHandler(service))
	mux.HandleFunc("GET /instance/status/simple/{instanceName}", InstanceStatusSimpleHandler(service))
	mux.HandleFunc("GET /server/status/{serverName}", ServerStatusHandler(service))
	mux.HandleFunc("GET /server/status/simple/{serverName}", ServerStatusSimpleHandler(service))
	return mux
}

// TargetStatusHandler - Target status handler
func TargetStatusHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := r.PathValue("targetName")
		_, ok := service.Targets[targetName]
		if !ok {
			http.Error(w, "Target not found", http.StatusNotFound)
			return
		}
		status, err := service.GetTargetStatus(targetName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// InstanceStatusSimpleHandler - Instance status simple handler
func InstanceStatusSimpleHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceName := r.PathValue("instanceName")
		_, ok := service.Instances[instanceName]
		if !ok {
			http.Error(w, "Instance not found", http.StatusNotFound)
			return
		}
		status, err := service.InstanceStatusSimple(instanceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status))
	}
}

// ServerStatusHandler - Server status handler
func ServerStatusHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverName := r.PathValue("serverName")
		_, ok := service.Instances[serverName]
		if !ok {
			http.Error(w, "Server not found", http.StatusNotFound)
			return
		}
		status, err := service.GetServerStatus(serverName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// ServerStatusSimpleHandler - Server status simple handler
func ServerStatusSimpleHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverName := r.PathValue("serverName")
		_, ok := service.Instances[serverName]
		if !ok {
			http.Error(w, "Server not found", http.StatusNotFound)
			return
		}
		status, err := service.ServerStatusSimple(serverName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status))
	}
}
