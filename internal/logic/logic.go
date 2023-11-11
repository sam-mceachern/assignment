package logic

import "example.com/internal/storage"

type Client struct {
	storage storage.Storage
}

func NewClient(storage storage.Storage) *Client {
	return &Client{
		storage: storage,
	}
}
