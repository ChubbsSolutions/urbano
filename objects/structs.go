package objects

//WordDataSlice represents a list of words in JSON.
type WordDataSlice struct {
	List []WordData `json:"list"`
}

//WordData represents the JSON struct sent by Urban Dictionary with the word.
type WordData struct {
	Author      string `json:"author"`
	CurrentVote string `json:"current_vote"`
	Defid       int    `json:"defid"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	Permalink   string `json:"permalink"`
	ThumbsUp    int    `json:"thumbs_up"`
	ThumbsDown  int    `json:"thumbs_down"`
	Word        string `json:"word"`
}
