package checkgrp

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type Handlers struct {
	Build  string
	Logger *zap.SugaredLogger
}

// Readiness is responsible for checking all the dependencies we expect
// to be up and running during the runtime e.x database.
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	// Set status for success branch.
	statusCode := http.StatusOK
	status := "OK"

	// Check db.
	err := check(ctx, h.Logger, dbCheck)
	if err != nil {
		// Set status for error branch
		statusCode = http.StatusInternalServerError
		status = "db not ready"
	}

	// Prepare response.
	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	// Send response to the client.
	err = response(w, statusCode, data)
	if err != nil {
		h.Logger.Errorw("readiness", "ERROR", err)
		return
	}

	h.Logger.Infow("readiness", "status", status, "code", statusCode)
}

// Liveness is responsible for returning info about the running service.
func (h Handlers) Liveness(w http.ResponseWriter, req *http.Request) {
	// Set current hostname
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	// Set status for success branch.
	statusCode := http.StatusOK
	status := "OK"

	// Prepare response.
	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    status,
		Build:     h.Build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	// Send response to the client.
	err = response(w, http.StatusOK, data)
	if err != nil {
		h.Logger.Errorw("liveness", "ERROR", err)
		return
	}

	h.Logger.Infow("liveness", "status", status, "code", statusCode)
}

// private

type checkFunc func(ctx context.Context, logger *zap.SugaredLogger) error

func check(ctx context.Context, logger *zap.SugaredLogger, fn checkFunc) error {
	return fn(ctx, logger)
}

func dbCheck(ctx context.Context, logger *zap.SugaredLogger) error {
	logger.Infow("readiness", "check", "db")
	<-time.After(100 * time.Millisecond)
	return nil
}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
