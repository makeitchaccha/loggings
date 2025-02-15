package logging

type Event string

const (
	EventMemberJoin Event = "member_join"
)

var Events = []Event{
	EventMemberJoin,
}

func (e Event) String() string {
	return string(e)
}

func ParseEvent(s string) (Event, bool) {
	for _, e := range Events {
		if e.String() == s {
			return e, true
		}
	}
	return "", false
}

func MustParseEvent(s string) Event {
	e, ok := ParseEvent(s)
	if !ok {
		panic("invalid event")
	}
	return e
}
