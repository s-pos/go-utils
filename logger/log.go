package logger

import (
	"context"
	"fmt"
)

const (
	BODY        = "body"
	QUERYSTRING = "query_string"
	PATH        = "path"
	HEADER      = "header"
)

type msg struct {
	text         string
	field        string
	fieldMessage string
}

// LogMessage create contract for log any message
type LogMessage interface {
	// Send message log To context and will show at the end of request
	To(ctx context.Context)

	// create mandatory fields on response with required request on Body param's
	Body(ctx context.Context)

	// create mandatory fields on response with required request on Path param's
	Path(ctx context.Context)

	// create mandatory fields on repsonse with required request on QueryString param's
	QueryString(ctx context.Context)

	// create mandatory fields on repsonse with required request on Header param's
	Header(ctx context.Context)
}

func (m msg) To(ctx context.Context) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.StoreMessage(m.text)
}

func (m msg) Body(ctx context.Context) {
	var Fields []MandatoryField

	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_MandatoryField)
	if ok {
		Fields = tmp.([]MandatoryField)
	}

	mf := MandatoryField{
		Field:    m.field,
		Location: BODY,
		Message:  m.fieldMessage,
	}

	Fields = append(Fields, mf)

	value.Set(_MandatoryField, Fields)
}

func (m msg) Path(ctx context.Context) {
	var Fields []MandatoryField

	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_MandatoryField)
	if ok {
		Fields = tmp.([]MandatoryField)
	}

	mf := MandatoryField{
		Field:    m.field,
		Location: PATH,
		Message:  m.fieldMessage,
	}

	Fields = append(Fields, mf)

	value.Set(_MandatoryField, Fields)
}

func (m msg) QueryString(ctx context.Context) {
	var Fields []MandatoryField

	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_MandatoryField)
	if ok {
		Fields = tmp.([]MandatoryField)
	}

	mf := MandatoryField{
		Field:    m.field,
		Location: QUERYSTRING,
		Message:  m.fieldMessage,
	}

	Fields = append(Fields, mf)

	value.Set(_MandatoryField, Fields)
}

func (m msg) Header(ctx context.Context) {
	var Fields []MandatoryField

	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_MandatoryField)
	if ok {
		Fields = tmp.([]MandatoryField)
	}

	mf := MandatoryField{
		Field:    m.field,
		Location: HEADER,
		Message:  m.fieldMessage,
	}

	Fields = append(Fields, mf)

	value.Set(_MandatoryField, Fields)
}

// Message default log message
func Message(log ...interface{}) LogMessage {
	return msg{
		text: fmt.Sprint(log...),
	}
}

// Messagef log message with format
func Messagef(format string, log ...interface{}) LogMessage {
	return msg{
		text: fmt.Sprintf(format, log...),
	}
}

// FieldMandatory record any field mandatory (requeired)
func FieldMandatory(field, fieldMessage string) LogMessage {
	return msg{
		field:        field,
		fieldMessage: fieldMessage,
	}
}
