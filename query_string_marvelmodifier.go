package querystring

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/martian"
	"github.com/google/martian/parse"
)

func init() {
	parse.Register("querystring.MarvelModifier", marvelModifierFromJSON)
}

// MarvelModifier contains the private and public Marvel API key
type MarvelModifier struct {
	public, private string
}

// MarvelModifierJSON to Unmarshal the JSON configuration
type MarvelModifierJSON struct {
	Public  string               `json:"public"`
	Private string               `json:"private"`
	Scope   []parse.ModifierType `json:"scope"`
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *MarvelModifier) ModifyRequest(req *http.Request) error {
	query := req.URL.Query()
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	hash := GetMD5Hash(ts + m.private + m.public)
	query.Set("apikey", m.public)
	query.Set("ts", ts)
	query.Set("hash", hash)
	req.URL.RawQuery = query.Encode()

	return nil
}

// GetMD5Hash returns the md5 hash from a string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// MarvelNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func MarvelNewModifier(public, private string) martian.RequestModifier {
	return &MarvelModifier{
		public:  public,
		private: private,
	}
}

// marvelModifierFromJSON takes a JSON message as a byte slice and returns
// a querystring.modifier and an error.
//
// Example JSON:
// {
//  "public": "apikey",
//  "private": "apikey",
//  "scope": ["request", "response"]
// }
func marvelModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &MarvelModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(MarvelNewModifier(msg.Public, msg.Private), msg.Scope)
}
