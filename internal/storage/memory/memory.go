package inmemory

type Client struct {
	storage map[string]string
}

func NewClient() *Client {
	return &Client{
		storage: make(map[string]string),
	}
}
