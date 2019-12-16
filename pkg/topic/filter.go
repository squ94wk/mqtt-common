package topic

//Filter defines a topic filter.
type Filter struct {
	topic string
}

//ParseFilter parses a Filter from a string.
func ParseFilter(filter string) Filter {
	return Filter{filter}
}

//String formats Filter as a string.
func (f Filter) String() string {
	return f.topic
}
