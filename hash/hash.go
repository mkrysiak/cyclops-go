package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"hash"
	"io"
	"sort"

	log "github.com/sirupsen/logrus"
)

type SentryData struct {
	Exception Exception `json:"exception"`
	Culprit   string    `json:"culprit"`
	Message   string    `json:"message"`

	// These fields are in the Attributes specification, but unused by us
	// for calculating the hash to determine if a given exception matches
	// https://docs.sentry.io/clientdev/attributes/

	// Extra       map[string]interface{} `json:"extra"`
	// EventId     string                 `json:"event_id"`
	// Project     string                 `json:"project"`
	// Logger      string                 `json:"platform"`
	// Request     Request                `json:"request"`
	// Platform    string                 `json:"platform"`
	// Stacktrace  []Stacktrace           `json:"stacktrace"`
	// Timestamp   time.Time              `json:"timestamp"`
	// Sdk         map[string]interface{} `json:"sdk"`
	// Level       string                 `json:"level"`
	// ServerName  string                 `json:"server_name"`
	// Release     string                 `json:"release"`
	// Tags        map[string]interface{} `json:"tags"`
	// Environment string                 `json:"environment"`
	// Modules     map[string]interface{} `json:"modules"`
	// Fingerprint []string               `json:"fingerprint"`
}

type Request struct {
	Headers map[string]string `json:"headers"`
	Url     string            `json:"url"`
}

type Exception struct {
	Values []Value `json:"values"`
}

type Value struct {
	Type       string     `json:"type"`
	Value      string     `json:"value"`
	Module     string     `json:"module"`
	ThreadId   string     `json:"thread_id"`
	Stacktrace Stacktrace `json:"stacktrace`
}

type Stacktrace struct {
	Frames            []Frame `json:"frames"`
	Package           string  `json:"package"`
	Platform          string  `json:"platform"`
	ImageAddr         string  `json:"image_addr"`
	SymbolAddr        string  `json:"symbol_addr"`
	InstructionOffset int     `json:"instruction_offset"`
}

type Frame struct {
	Filename    string            `json:"filename"`
	Lineno      int               `json:"lineno"`
	Colno       int               `json:"colno"`
	Function    string            `json:"function"`
	InApp       bool              `json:"in_app"`
	Module      string            `json:"module"`
	AbsPath     string            `json:"abs_path"`
	ContextLine string            `json:"context_line"`
	PreContext  []string          `json:"pre_context"`
	PostContext []string          `json:"post_context"`
	Vars        map[string]string `json:"vars"`
}

type Template struct {
	AbsPath     string   `json:"abs_path"`
	ContextLine string   `json:"context_line"`
	Filename    string   `json:"filename"`
	Lineno      int      `json:"lineno"`
	PreContext  []string `json:"pre_context"`
	PostContext []string `json:"post_context"`
}

var hasher hash.Hash

// Try to calculate the hash using the most specific data first (stacktraces),
// but if they don't exist fallback to just the message and culprit
func HashForGrouping(requestBody []byte) (string, error) {
	hasher = md5.New()
	var sentryData SentryData
	err := json.Unmarshal(requestBody, &sentryData)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if len(sentryData.Exception.Values) > 0 {
		return sentryData.hashForFramesInException(), nil
	} else if sentryData.Message != "" || sentryData.Culprit != "" {
		return sentryData.hashStrings([]string{sentryData.Message, sentryData.Culprit}), nil
	}
	return "", err
}

func (s *SentryData) hashStrings(str []string) string {
	if len(str) < 1 {
		return ""
	}
	sort.Strings(str)
	for _, v := range str {
		io.WriteString(hasher, v)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *SentryData) hashForFramesInException() string {
	var hashItems []string
	for _, v := range s.Exception.Values {
		hashItems = append(hashItems, v.Type)
		for _, f := range v.Stacktrace.Frames {
			out, err := json.Marshal(f)
			if err == nil {
				//TODO: Waste of cpu cycles to convert "out" from []byte to string back to a []byte
				// in hashStrings()
				hashItems = append(hashItems, string(out))
			}
		}
	}

	h := s.hashStrings(hashItems)
	return h
}

func (s *SentryData) hashForFramesInStackTrace(stacktrace Stacktrace) {
}
