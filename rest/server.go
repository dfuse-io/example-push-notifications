package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	msignotify "github.com/dfuse-io/example-push-notifications"
)

type Server struct {
	addr    string
	storage msignotify.Storage
}

func NewServer(addr string, storage msignotify.Storage) *Server {
	return &Server{
		addr:    addr,
		storage: storage,
	}
}
func (s *Server) Listen() error {
	http.HandleFunc("/optin", s.OptinDevice)
	fmt.Println("Listening at:", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

type OptinRequest struct {
	Token      string `json:"token"`
	Account    string `json:"account"`
	DeviceType string `json:"device_type"`
}

func (s *Server) OptinDevice(w http.ResponseWriter, r *http.Request) {

	var optinRequest *OptinRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&optinRequest)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Failed to decode optin request"))
	}

	fmt.Println("Received an opt in request:", optinRequest)
	s.storage.OptInDeviceToken(optinRequest.Account, optinRequest.Token, optinRequest.DeviceType)
}
