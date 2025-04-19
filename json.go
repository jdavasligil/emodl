package emotedownloader

// Used to unmarshal errors from an API response
//
//easyjson:json
type jsonError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
