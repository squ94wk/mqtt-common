package topic

import "strings"

//Topic defines an mqtt topic as a structured type.
type Topic struct {
	Levels []string
}

//Filter analogously defines an mqtt topic filter.
type Filter struct {
	Levels []string
}

//String prints the topic as string by inserting back the '/' separator.
func (t Topic) String() string {
	return strings.Join(t.Levels, "/")
}

//String prints the topic filter as string by inserting back the '/' separator.
func (f Filter) String() string {
	return strings.Join(f.Levels, "/")
}

//ParseTopic parses a topic from an input string.
func ParseTopic(input string) (Topic, error) {
	levels := strings.Split(input, "/")
	return Topic{
		Levels: levels,
	}, nil
}

//ParseFilter parses a topic filter from an input string.
func ParseFilter(input string) (Filter, error) {
	levels := strings.Split(input, "/")
	return Filter{
		Levels: levels,
	}, nil
}
