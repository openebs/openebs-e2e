package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sync"

	"strings"

	nats "github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
)

var g_conn *nats.Conn
var g_stream nats.JetStreamContext

var g_eventSubscriptions map[*nats.Subscription]EventSubscription
var g_eventServer string
var g_mutex sync.Mutex

const STREAM_NAME = "events-stream"
const SUBJECT_PATTERN = "events.>"

type EventMessage struct {
	subject string
	message string
}

type EventSubscription struct {
	subject_pattern string
	messages        []EventMessage
}

type EventSubscriptionRequest struct {
	EventServerAddr string `json:"eventServerAddr"`
	Subject         string `json:"subject"`
}

type EventPublishRequest struct {
	EventServerAddr string `json:"eventServerAddr"`
	Subject         string `json:"subject"`
	Data            string `json:"data"`
}

type EventRequest struct {
	Subject string `json:"subject"`
}

func init() {
	g_eventSubscriptions = make(map[*nats.Subscription]EventSubscription)
}

func messageHandler(msg *nats.Msg) {
	g_mutex.Lock()
	defer g_mutex.Unlock()

	klog.Info("messageHandler, message subject: ", msg.Subject,
		" data: ", string(msg.Data),
	)

	eventsubscription, ok := g_eventSubscriptions[msg.Sub]
	if !ok {
		klog.Error("messageHandler: internal error: failed to find subscription")
		return
	}

	message := EventMessage{
		subject: msg.Subject,
		message: string(msg.Data),
	}

	eventsubscription.messages = append(eventsubscription.messages, message)
	klog.Info("messageHandler: subscription subject ", eventsubscription.subject_pattern)

	g_eventSubscriptions[msg.Sub] = eventsubscription
}

func EventList(w http.ResponseWriter, r *http.Request) {
	var event_req EventRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&event_req); err != nil {
		msg := fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	klog.Info("EventList: subject: ", event_req.Subject)

	g_mutex.Lock()
	defer g_mutex.Unlock()

	var returnedmessages []string
	// return the list of event messages as a json array of strings
	for _, subscrip := range g_eventSubscriptions {
		klog.Info("EventList, found a subscription to: ", subscrip.subject_pattern)
		for _, msg := range subscrip.messages {
			klog.Info("EventList, found a subject: ", msg)
			klog.Info("EventList, checking for prefix: ", event_req.Subject)
			if strings.HasPrefix(msg.subject, event_req.Subject) {
				klog.Info("EventList, appending messages")
				returnedmessages = append(returnedmessages, msg.message)
				klog.Info("EventList, matched subject: ", event_req.Subject)
			}
		}
	}
	jsn, err := json.Marshal(returnedmessages)
	if err != nil {
		msg := fmt.Sprintf("EventList: failed to marshal messages, Error: %s", err.Error())
		WrapResult(msg, ErrJsonEncode, w)
		return
	}
	WrapResult(string(jsn), ErrNone, w)
}

func ensureConnection(eventServerAddr string) error {
	var err error
	if g_conn == nil {
		g_conn, err = nats.Connect(eventServerAddr)
		if err != nil {
			klog.Error("ensureConnection: failed to connect, error: ", err)
			return err
		} else if g_conn == nil {
			klog.Error("ensureConnection: failed to connect")
			return err
		} else {
			klog.Info("ensureConnection: connected")
		}
		g_eventServer = eventServerAddr
	} else {
		klog.Info("ensureConnection: already connected")
		if eventServerAddr != g_eventServer {
			err = fmt.Errorf("ensureConnection: already subscribed to a different server %s", g_eventServer)
		}
	}
	return err
}

func ensureStream() error {
	var err error
	if g_stream == nil {
		g_stream, err = g_conn.JetStream()
		if err != nil {
			return err
		}
		_, err = g_stream.AddStream(&nats.StreamConfig{
			Name:     STREAM_NAME,
			Subjects: []string{SUBJECT_PATTERN},
		})
		if err != nil {
			// tolerate this but log, the stream probably exists
			klog.Error("ensureStream: failed to add stream, error: ", err.Error())
		} else {
			klog.Info("ensureStream: added stream")
		}
	}
	return nil
}

func EventSubscribe(w http.ResponseWriter, r *http.Request) {
	var subscription_req EventSubscriptionRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&subscription_req); err != nil {
		msg := fmt.Sprintf("EventSubscribe: failed to read JSON encoded data, Error: %s", err.Error())
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	g_mutex.Lock()
	defer g_mutex.Unlock()

	klog.Info("EventSubscribe: attempting to subscribe to ", subscription_req.Subject)
	var err error
	if err = ensureConnection(subscription_req.EventServerAddr); err != nil {
		msg := fmt.Sprintf("EventSubscribe: failed to connect to %s", subscription_req.EventServerAddr)
		WrapResult(msg, ErrConnectFail, w)
		return
	}

	if err = ensureStream(); err != nil {
		msg := fmt.Sprintf("EventSubscribe: failed to create stream, error: %s", err.Error())
		WrapResult(msg, ErrStreamCreateFail, w)
		return
	}

	// if already subscribed to this subject, error
	for _, v := range g_eventSubscriptions {
		if v.subject_pattern == subscription_req.Subject {
			msg := fmt.Sprintf("EventSubscribe: already subscribed to this subject %s", v.subject_pattern)
			WrapResult(msg, ErrSubscribedAlready, w)
			return
		}
	}

	// add subscription to map
	subscription, err := g_stream.Subscribe(subscription_req.Subject, messageHandler)
	if err != nil {
		msg := fmt.Sprintf("EventSubscribe: failed to subscribe, error: %s", err)
		WrapResult(msg, ErrSubscribeFail, w)
	} else {
		klog.Info("EventSubscribe: new subscription to ", subscription_req.Subject)
		g_eventSubscriptions[subscription] = EventSubscription{
			subject_pattern: subscription_req.Subject,
			messages:        []EventMessage{},
		}
		WrapResult("", ErrNone, w)
	}
}

func EventPublish(w http.ResponseWriter, r *http.Request) {
	var publish_req EventPublishRequest
	d := json.NewDecoder(r.Body)
	var err error
	if err = d.Decode(&publish_req); err != nil {
		msg := fmt.Sprintf("EventPublish: failed to read JSON encoded data, Error: %s", err.Error())
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	klog.Info("EventPublish: attempting to publish to ", publish_req.Subject)
	g_mutex.Lock()
	defer g_mutex.Unlock()

	if err = ensureConnection(publish_req.EventServerAddr); err != nil {
		msg := fmt.Sprintf("EventPublish: failed to connect to %s", publish_req.EventServerAddr)
		WrapResult(msg, ErrConnectFail, w)
		return
	}

	if err = ensureStream(); err != nil {
		msg := fmt.Sprintf("EventPublish: failed to create stream, error: %s", err.Error())
		WrapResult(msg, ErrStreamCreateFail, w)
		return
	}

	_, err = g_stream.Publish(publish_req.Subject, []byte(publish_req.Data))
	if err != nil {
		msg := fmt.Sprintf("EventPublish: failed to publish to nats service, error: %s", err.Error())
		WrapResult(msg, ErrPublishFail, w)
	} else {
		WrapResult("", ErrNone, w)
	}
}

func EventUnsubscribe(w http.ResponseWriter, r *http.Request) {
	var event_req EventRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&event_req); err != nil {
		msg := fmt.Sprintf("EventSubscribe: failed to read JSON encoded data, Error: %s", err.Error())
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	g_mutex.Lock()
	defer g_mutex.Unlock()

	// get subject
	// find subscription in map
	// if not there, error
	for k, v := range g_eventSubscriptions {
		if v.subject_pattern == event_req.Subject {
			err := k.Unsubscribe()
			delete(g_eventSubscriptions, k)
			if err != nil {
				WrapResult(fmt.Sprintf("EventUnsubscribe: unsubscribe failed %v", err), ErrUnSubscribeFail, w)
			} else {
				WrapResult("", ErrNone, w)
			}
			return
		}
	}

	msg := fmt.Sprintf("EventUnsubscribe: subscription %s not found", event_req.Subject)
	WrapResult(msg, ErrSubscriptionNotFound, w)
}

func EventUnsubscribeAll(w http.ResponseWriter, r *http.Request) {
	// iterate through subscriptions in map and unsubscribe all
	// empty map
	klog.Info("EventUnsubscribeAll")

	g_mutex.Lock()
	defer g_mutex.Unlock()

	if g_conn != nil {
		var err error
		klog.Info("EventUnsubscribeAll: was connected")
		for k := range g_eventSubscriptions {
			err = k.Unsubscribe()
			if err != nil {
				break
			}
			delete(g_eventSubscriptions, k)
		}
		g_conn.Close()
		g_conn = nil
		g_eventServer = ""
		g_stream = nil
		if err != nil {
			WrapResult(fmt.Sprintf("EventUnsubscribeAll: unsubscribe failed %v", err), ErrUnSubscribeFail, w)
		} else {
			WrapResult("", ErrNone, w)
		}
	} else {
		WrapResult("EventUnsubscribeAll: was not connected", ErrNotConnected, w)
	}
}
