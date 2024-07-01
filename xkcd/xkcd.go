package xkcd

const XkcdLink = "https://xkcd.com/"

type Comics struct {
	Num        int
	Month      string
	Year       string
	SafeTitle  string `json:"safe_title"`
	transcript string
}
