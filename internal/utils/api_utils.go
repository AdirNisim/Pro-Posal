package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshalRequest(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("failed decoding request JSON: %w", err)
	}

	return nil
}
