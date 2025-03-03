package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

type StandardResponse struct {
	Error   *string `json:"error"`
	Status  bool    `json:"status"`
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message"`
}

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func WriteJSONResponse(w http.ResponseWriter, message *string, data any) {
	SetJSONHeader(w)

	response := StandardResponse{
		Error:   nil,
		Status:  true,
		Data:    data,
		Message: message,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func ErrorMessage(w http.ResponseWriter, err error, message *string) {
	SetJSONHeader(w)

	var str = "1"
	var hash *string = &str
	if err != nil {
		hash = hashError(err)
	}

	if message == nil {
		msg := "Something Went Wrong"
		message = &msg
	}
	response := StandardResponse{
		Error:   hash,
		Status:  false,
		Data:    nil,
		Message: message,
	}
	log.Println(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func hashError(err error) *string {
	hasher := sha256.New()
	hasher.Write([]byte(err.Error()))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return &hash
}
