package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"projects/emergency-messages/internal/services"
)

type Receiver struct {
	receiverService *services.ReceiverService
	log             *slog.Logger
}

func NewReceiver(receiverService *services.ReceiverService, log *slog.Logger) *Receiver {
	return &Receiver{
		receiverService: receiverService,
		log:             log,
	}
}

func (c *Receiver) GetByCity(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	ctx := context.Background()
	receivers, err := c.receiverService.FindByCity(ctx, city)
	if assertError(err, w) {
		c.log.Error("getting by city", err)
		return
	}

	receiversBytes, err := json.Marshal(receivers)
	if err != nil {
		c.log.Error("cannot marshalling receivers")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(receiversBytes)
	w.WriteHeader(http.StatusOK)
}

func (c *Receiver) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	receiversCreated, err := c.receiverService.Upload(file)
	if assertError(err, w) {
		c.log.Error("uploading receiver", err)
		return
	}

	receiversBytes, err := json.Marshal(receiversCreated)
	if err != nil {
		c.log.Error("cannot marshalling receivers")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(receiversBytes)
	w.WriteHeader(http.StatusCreated)
}
