package server

import (
	"encoding/json"

	"github.com/dapr/go-sdk/service/common"
)

func JSON(in any) (*common.Content, error) {
	data, err := json.Marshal(&in)
	if err != nil {
		return nil, err
	}
	return &common.Content{
		ContentType: "application/json",
		Data:        data,
	}, nil
}
