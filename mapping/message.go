package mapping

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// Message is a structure to be mapped into JSON
type Message struct {
	Headers map[string]interface{} `json:",omitempty"`

	ContentType     string `json:",omitempty"` // MIME content type
	ContentEncoding string `json:",omitempty"` // MIME content encoding
	DeliveryMode    uint8  `json:",omitempty"` // queue implementation use - non-persistent (1) or persistent (2)
	Priority        uint8  `json:",omitempty"` // queue implementation use - 0 to 9
	CorrelationId   string `json:",omitempty"` // application use - correlation identifier
	ReplyTo         string `json:",omitempty"` // application use - address to reply to (ex: RPC)
	Expiration      string `json:",omitempty"` // implementation use - message expiration spec
	MessageId       string `json:",omitempty"` // application use - message identifier
	Timestamp       int64  `json:",omitempty"` // application use - message timestamp
	Type            string `json:",omitempty"` // application use - message type name
	UserId          string `json:",omitempty"` // application use - creating user - should be authenticated user
	AppId           string `json:",omitempty"` // application use - creating application id

	DeliveryTag uint64 `json:",omitempty"`
	Redelivered bool   `json:",omitempty"`
	Exchange    string `json:",omitempty"` // basic.publish exchange
	RoutingKey  string `json:",omitempty"` // basic.publish routing key

	BodyBase64 string `json:",omitempty"` // Body in base64 format

	TakenAt int64 `json:",omitempty"` // Time this message was taken from queue
}

// FromJSON constructs message from JSON string
func FromJSON(s string) (*Message, error) {
	var m Message
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// FromDelivery constructs message from delivery
func FromDelivery(a amqp.Delivery) Message {
	m := Message{
		Headers: a.Headers,

		ContentType:     a.ContentType,
		ContentEncoding: a.ContentEncoding,
		DeliveryMode:    a.DeliveryMode,
		Priority:        a.Priority,
		CorrelationId:   a.CorrelationId,
		ReplyTo:         a.ReplyTo,
		Expiration:      a.Expiration,
		MessageId:       a.MessageId,
		Timestamp:       a.Timestamp.Unix(),
		Type:            a.Type,
		UserId:          a.UserId,
		AppId:           a.AppId,

		DeliveryTag: a.DeliveryTag,
		Redelivered: a.Redelivered,
		Exchange:    a.Exchange,
		RoutingKey:  a.RoutingKey,

		BodyBase64: base64.StdEncoding.EncodeToString(a.Body),

		TakenAt: time.Now().Unix(),
	}

	return m
}

// GetTimestamp returns message timestamp
func (m Message) GetTimestamp() time.Time {
	return time.Unix(m.Timestamp, 0)
}

// GetTakenAt return time this message was taken from queue
func (m Message) GetTakenAt() time.Time {
	return time.Unix(m.TakenAt, 0)
}

// GetBody returns body as byte slice
func (m Message) GetBody() []byte {
	bts, err := base64.StdEncoding.DecodeString(m.BodyBase64)
	if err != nil {
		return []byte{}
	}
	return bts
}

// Fprint outputs mapped contents into given writer
func (m Message) Fprint(w io.Writer, formatted bool) (err error) {
	if w == nil {
		w = os.Stdout
	}

	var bts []byte
	if formatted {
		bts, err = json.MarshalIndent(m, "", "  ")
	} else {
		bts, err = json.Marshal(m)
	}
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, string(bts))
	return
}

// ToPublishing converts message to AMQP publishing
func (m Message) ToPublishing() amqp.Publishing {
	return amqp.Publishing{
		Headers:         stripUnsupportedPublishingHeaders(m.Headers),
		ContentType:     m.ContentType,
		ContentEncoding: m.ContentEncoding,
		DeliveryMode:    m.DeliveryMode,
		Priority:        m.Priority,
		CorrelationId:   m.CorrelationId,
		ReplyTo:         m.ReplyTo,
		Expiration:      m.Expiration,
		MessageId:       m.MessageId,
		Timestamp:       m.GetTimestamp(),
		Type:            m.Type,
		UserId:          m.UserId,
		AppId:           m.AppId,
		Body:            m.GetBody(),
	}
}

// Publish publishes message to given exchange using given channel
func (m Message) Publish(ch *amqp.Channel, exchange, routingKey string) error {
	if ch == nil {
		return errors.New("nil channel")
	}
	if len(exchange) == 0 && len(routingKey) == 0 {
		return errors.New("empty both exchange and routing key")
	}
	return ch.Publish(exchange, m.prepareRoutingKey(routingKey), true, false, m.ToPublishing())
}

func (m Message) prepareRoutingKey(routingKey string) string {
	if len(routingKey) < 2 {
		return routingKey
	}
	if len(routingKey) > 1 && (routingKey[0] == '@' || routingKey[0] == '#' || routingKey[0] == '$') {
		// Reading routing key from header
		header := routingKey[1:]
		if len(m.Headers) > 0 {
			if h, ok := m.Headers[header]; ok {
				// Header is present
				return fmt.Sprint(h)
			}
		}
		return ""
	}
	return routingKey
}

// stripUnsupportedPublishingHeaders removes unsupported headers
func stripUnsupportedPublishingHeaders(m map[string]interface{}) map[string]interface{} {
	if len(m) == 0 {
		// No data
		return m
	}
	if _, ok := m["x-death"]; !ok {
		// No unsupported header
		return m
	}

	out := map[string]interface{}{}
	for k, v := range m {
		if k == "x-death" {
			// Skip
		} else {
			out[k] = v
		}
	}
	return out
}
