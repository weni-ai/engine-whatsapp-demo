package metric

// ChannelCreation represents a channel creation metric
type ChannelCreation struct {
	Channel string
}

// ChannelCreation returns new metric struct value representation.
func NewChannelCreation(channel string) *ChannelCreation {
	return &ChannelCreation{Channel: channel}
}

// ContactMessage represents a contact message metric.
type ContactMessage struct {
	Channel string
}

// ContactMessage returns new metric struct value representation.
func NewContactMessage(channel string) *ContactMessage {
	return &ContactMessage{Channel: channel}
}

// ContactActivation represents a contact activation ContactActivated metric.
type ContactActivation struct {
	Channel string
}

// ContactActivation returns new metric struct value representation.
func NewContactActivation(channel string) *ContactActivation {
	return &ContactActivation{Channel: channel}
}

// ContactActivated represents a current contact activated gauge metric.
type ContactActivated struct {
	Channel string
}

// ContactActivated returns new metric struct value representation.
func NewContactActivated(channel string) *ContactActivated {
	return &ContactActivated{Channel: channel}
}

// Metric encapsulates interface metric definitions
type Metric interface {
	SaveChannelCreation(m *ChannelCreation)
	SaveContactMessage(m *ContactMessage)
	SaveContactActivation(m *ContactActivation)
	IncContactActivated(m *ContactActivated)
	DecContactActivated(m *ContactActivated)
}
