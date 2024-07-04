package event

type eventsFilter struct {
	events []EventMessage
	err    error
}

func (context *EventContext) GetFilteredEvents(subject_pattern string) *eventsFilter {
	events, err := context.GetEvents(subject_pattern)
	return &eventsFilter{
		events,
		err,
	}
}

func (eventsfilter *eventsFilter) Build() ([]EventMessage, error) {
	if eventsfilter.err != nil {
		return []EventMessage{}, eventsfilter.err
	}
	return eventsfilter.events, eventsfilter.err
}

func (eventsfilter *eventsFilter) WithAction(action Action) *eventsFilter {
	if eventsfilter.err == nil {
		var newList []EventMessage
		for _, e := range eventsfilter.events {
			if e.Action == action {
				newList = append(newList, e)
			}
		}
		eventsfilter.events = newList
	}
	return eventsfilter
}

func (eventsfilter *eventsFilter) WithCategory(category Category) *eventsFilter {
	if eventsfilter.err == nil {
		var newList []EventMessage
		for _, e := range eventsfilter.events {
			if e.Category == category {
				newList = append(newList, e)
			}
		}
		eventsfilter.events = newList
	}
	return eventsfilter
}

func (eventsfilter *eventsFilter) WithTarget(target string) *eventsFilter {
	if eventsfilter.err == nil {
		var newList []EventMessage
		for _, e := range eventsfilter.events {
			if e.Target == target {
				newList = append(newList, e)
			}
		}
		eventsfilter.events = newList
	}
	return eventsfilter
}
