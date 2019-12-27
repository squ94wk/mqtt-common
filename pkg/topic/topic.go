package topic

import "strings"

//Topic defines an mqtt topic as a structured type.
type Topic struct {
	levels []string
}

//Filter analogously defines an mqtt topic filter.
type Filter struct {
	levels []string
}

//Levels returns the topic levels of a topic.
func (t Topic) Levels() []string {
	return t.levels
}

//Levels returns the topic levels of a topic filter.
func (f Filter) Levels() []string {
	return f.levels
}

//String prints the topic as string by inserting back the '/' separator.
func (t Topic) String() string {
	return strings.Join(t.levels, "/")
}

//String prints the topic filter as string by inserting back the '/' separator.
func (f Filter) String() string {
	return strings.Join(f.levels, "/")
}

//ParseTopic parses a topic from an input string.
func ParseTopic(input string) (Topic, error) {
	levels := strings.Split(input, "/")
	return Topic{
		levels: levels,
	}, nil
}

//ParseFilter parses a topic filter from an input string.
func ParseFilter(input string) (Filter, error) {
	levels := strings.Split(input, "/")
	return Filter{
		levels: levels,
	}, nil
}
