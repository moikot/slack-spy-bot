package main

import (
	"reflect"
	"regexp"

	"github.com/nlopes/slack"
)

type Router interface {
	RouteEvent(event interface{})
	RouteMessage(event *slack.MessageEvent)

	AddEventHandler(handler HandlerFunc)
	AddMessageHandler(regexp *regexp.Regexp, handler HandlerFunc)
}

func NewRouter() Router {
	return &router{
		eventHandlers:   map[string]HandlerFunc{},
		messageHandlers: map[*regexp.Regexp]HandlerFunc{},
	}
}

type router struct {
	eventHandlers   map[string]HandlerFunc
	messageHandlers map[*regexp.Regexp]HandlerFunc
}

func (r *router) AddEventHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventType := handlerType.In(0).Elem().Name()
	r.eventHandlers[eventType] = handler
}

func (r *router) AddMessageHandler(regexp *regexp.Regexp, handler HandlerFunc) {
	r.messageHandlers[regexp] = handler
}

func (r *router) RouteMessage(event *slack.MessageEvent) {
	for pattern, handler := range r.messageHandlers {
		match := pattern.FindAllStringSubmatch(event.Text, -1)
		pattern.SubexpNames()
		if match != nil {
			for _, captures := range match {
				var params = make([]reflect.Value, len(captures))
				params[0] = reflect.ValueOf(event)

				for index, capture := range captures {
					if index == 0 {
						continue
					}
					params[index] = reflect.ValueOf(capture)
				}

				reflect.ValueOf(handler).Call(params)
			}
		}
	}
}

func (r *router) RouteEvent(event interface{}) {
	var eventType = reflect.TypeOf(event).Elem().Name()

	var handler = r.eventHandlers[eventType]
	if handler == nil {
		return
	}

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(event)

	reflect.ValueOf(handler).Call(params)
}
