package logger

import "context"

// Values create contract store value into context
type Values interface {
	Set(key Flags, value interface{})
	Load(key Flags) (interface{}, bool)
	LoadAndDelete(key Flags) (interface{}, bool)
	StoreMessage(message string)
}

// Set set value to keys
func (l *Locker) Set(key Flags, value interface{}) {
	l.data.Store(key, value)
}

// Load load value from key
func (l *Locker) Load(key Flags) (interface{}, bool) {
	return l.data.Load(key)
}

// LoadAndDelete load and delete keys
func (l *Locker) LoadAndDelete(key Flags) (interface{}, bool) {
	return l.data.LoadAndDelete(key)
}

// StoreMessage storing logmessage into key
func (l *Locker) StoreMessage(message string) {
	var msg []string

	if len(message) < 1 {
		return
	}

	tmp, ok := l.data.LoadAndDelete(_Messages)
	if ok {
		msg = tmp.([]string)
	}

	msg = append(msg, message)

	l.Set(_Messages, msg)
}

func extract(ctx context.Context) (Values, bool) {
	var (
		lock = new(Locker)
		ok   bool
	)
	if ctx == nil {
		return lock, false
	}

	lock, ok = ctx.Value(logKey).(*Locker)
	return lock, ok
}
