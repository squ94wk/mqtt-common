package topic

type Filter struct {
	topic string
}

func ParseFilter(filter string) Filter {
	return Filter{filter}
}

func (f Filter) String() string {
	return f.topic
}
