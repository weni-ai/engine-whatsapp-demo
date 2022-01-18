package metric

import "github.com/prometheus/client_golang/prometheus"

// Service implements metric interface
type Service struct {
	channelsCreations   *prometheus.CounterVec
	contactsMessages    *prometheus.CounterVec
	contactsActivations *prometheus.CounterVec
	contactsActivated   *prometheus.GaugeVec
}

// NewPrometheusService returns a new metric service
func NewPrometheusService() (*Service, error) {
	channelsCreations := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "channels_creations",
		Help: "Channel creations counter labeled by channel",
	}, []string{"channel"})

	contactsMessages := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "contacts_messages",
		Help: "Contact messages counter labeled by channel",
	}, []string{"channel"})

	contactsActivations := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "contacts_activations",
		Help: "Contact activation counter labeled by channel",
	}, []string{"channel"})

	contactsActivated := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "contacts_activated",
		Help: "Contact activated gauge labeled by channel",
	}, []string{"channel"})

	s := &Service{
		channelsCreations:   channelsCreations,
		contactsMessages:    contactsMessages,
		contactsActivations: contactsActivations,
		contactsActivated:   contactsActivated,
	}

	err := prometheus.Register(s.channelsCreations)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(s.contactsMessages)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(s.contactsActivations)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(s.contactsActivated)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	return s, nil
}

// receive a *metric.ChannelCreation metric and save to a Counter metric type.
func (s *Service) SaveChannelCreation(cc *ChannelCreation) {
	s.channelsCreations.WithLabelValues(cc.Channel).Inc()
}

// receive a *metric.ContactMessage metric and save to a Histogram metric type.
func (s *Service) SaveContactMessage(cm *ContactMessage) {
	s.contactsMessages.WithLabelValues(cm.Channel).Inc()
}

// receive a *metric.ContactActivation metric and save to a Histogram metric type.
func (s *Service) SaveContactActivation(ca *ContactActivation) {
	s.contactsActivations.WithLabelValues(ca.Channel).Inc()
}

// receive a *metric.ContactActivated metric and increment to a Gauge metric type.
func (s *Service) IncContactActivated(ca *ContactActivated) {
	s.contactsActivated.WithLabelValues(ca.Channel).Inc()
}

// receive a *metric.ContactActivated metric and decrement to a Gauge metric type.
func (s *Service) DecContactActivated(ca *ContactActivated) {
	s.contactsActivated.WithLabelValues(ca.Channel).Dec()
}
