package metric

// ChannelCreation represents a channel creation metric
type ChannelCreation struct {
}

// ChannelCreation returns new metric struct value representation.
func NewChannelCreation() *ChannelCreation {
	return &ChannelCreation{}
}

// ClientMessage represents a client message metric.
type ClientMessage struct {
	Channel  string
	Duration float64
}

// ClientMessage returns new metric struct value representation.
func NewClientMessage() *ClientMessage {
	return &ClientMessage{}
}

// ContactActivation represents a contact activation ContactActivated metric.
type ContactActivation struct {
	Channel  string
	Duration float64
}

// ContactActivation returns new metric struct value representation.
func NewContactActivation() *ContactActivation {
	return &ContactActivation{}
}

// ContactActivated represents a current contact activated gauge metric.
type ContactActivated struct {
	Channel  string
	Duration float64
}

// ContactActivated returns new metric struct value representation.
func NewContactActivated() *ContactActivated {
	return &ContactActivated{}
}

// Metric encapsulates interface metric definitions
type Metric interface {
	SaveChannelCreation(m *ChannelCreation)
	SaveClientMessage(m *ClientMessage)
	SaveContactActivation(m *ContactActivation)
	IncContactActivated(m *ContactActivated)
	DecContactActivated(m *ContactActivated)
}
