package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"vshop/lib/aescbc"
)

var ERR_NO_RESPONSE = errors.New("encrypted response not available")

func (s *Service) Payment(w http.ResponseWriter, r *http.Request) error {

	encResp := r.FormValue("encResp")
	if encResp == "" {
		return ERR_NO_RESPONSE
	}

	payload := strings.TrimSpace(encResp)
	buf, err := hex.DecodeString(payload)
	if err != nil {
		return fmt.Errorf("error decoding hex string %s: %s", payload, err)
	}

	resp, err := aescbc.NewCrypter().Decrypt(buf)
	if err != nil {
		return fmt.Errorf("error decrypting payload: %s", err)
	}
	jsonBytes, err := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

	return nil
}
