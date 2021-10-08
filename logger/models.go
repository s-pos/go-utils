package logger

import (
	"sync"
	"time"
)

// Locker is container data
type Locker struct {
	data sync.Map
}

type (
	// Key of context
	Key int

	// Flags is key for store context
	Flags string
)

const (
	logKey                 = Key(27)
	_StatusCode      Flags = "StatusCode"
	_Response        Flags = "Response"
	_Messages        Flags = "Messages"
	_ThirdParties    Flags = "ThirdParties"
	_ErrorMessage    Flags = "ErrorMessage"
	_ErrorLocation   Flags = "ErrorLocation"
	_ResponseMessage Flags = "ResponseMessage"
	_MandatoryField  Flags = "MandatoryFields"
)

// DataLogger is standard output to terminal
type DataLogger struct {
	RequestID       string                 `json:"request_id"`
	Service         string                 `json:"service"`
	TimeStart       time.Time              `json:"time_start"`
	Host            string                 `json:"host"`
	Endpoint        string                 `json:"endpoint"`
	RequestMethod   string                 `json:"request_method"`
	RequestHeader   map[string]interface{} `json:"request_header"`
	RequestBody     map[string]interface{} `json:"request_body"`
	StatusCode      int                    `json:"status_code"`
	Response        interface{}            `json:"response_body"`
	ResponseMessage string                 `json:"response_message"`
	ErrorMessage    string                 `json:"error_message"`
	ErrorLocation   string                 `json:"error_location"`
	ExecTime        float64                `json:"exec_time"`
	Messages        []string               `json:"log_messages"`
	ThirdParties    []ThirdParty           `json:"outgoing_log"`

	// not used for logging. use only for response
	MandatoryFields MandatoryField `json:"-"`
}

// ThirdParty is data logging for any request to third party (outside)
type ThirdParty struct {
	URL           string      `json:"url"`
	RequestHeader interface{} `json:"request_header"`
	RequestBody   interface{} `json:"request_body"`
	Response      interface{} `json:"response"`
	Method        string      `json:"method"`
	StatusCode    int         `json:"status_code"`
	ExecTime      float64     `json:"execution_time"`
}

// MandatoryField is response for any field requeired and need to save into context
type MandatoryField struct {
	Field    string `json:"field"`
	Location string `json:"location"`
	Message  string `json:"message"`
}
