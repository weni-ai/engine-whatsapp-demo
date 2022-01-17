package metric

import "github.com/prometheus/client_golang/prometheus"

// Service implements metric interface
type Service struct {
	channelCreations   prometheus.Counter
	clientMessages     *prometheus.HistogramVec
	contactActivations *prometheus.HistogramVec
	contactActivated   *prometheus.GaugeVec
}

// NewPrometheusService returns a new metric service
func NewPrometheusService() (*Service, error) {
	channelCreations := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "channel_creations",
		Help: "Channel creations counter",
	})

	clientMessages := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "client_messages",
		Help: "Client messages labeled by channel",
	}, []string{"channel"})

	contactActivations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "contact_activations",
		Help: "Contact activation labeled by channel",
	}, []string{"channel"})

	contactActivated := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "contact_activated",
		Help: "Contact activated gauge labeled by channel",
	}, []string{"channel"})

	s := &Service{
		channelCreations:   channelCreations,
		clientMessages:     clientMessages,
		contactActivations: contactActivations,
		contactActivated:   contactActivated,
	}

	err := prometheus.Register(s.channelCreations)
	if err != nil {
		return nil, err
	}

	err = prometheus.Register(s.clientMessages)
	if err != nil {
		return nil, err
	}

	err = prometheus.Register(s.contactActivations)
	if err != nil {
		return nil, err
	}

	err = prometheus.Register(s.contactActivated)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// receive a *metric.ChannelCreation metric and save to a Counter metric type.
func (s *Service) SaveChannelCreation(cc *ChannelCreation) {
	s.channelCreations.Inc()
}

// receive a *metric.ClientMessage metric and save to a Histogram metric type.
func (s *Service) SaveClientMessage(cm *ClientMessage) {
	s.clientMessages.WithLabelValues(cm.Channel).Observe(cm.Duration)
}

// receive a *metric.ContactActivation metric and save to a Histogram metric type.
func (s *Service) SaveContactActivation(ca *ContactActivation) {
	s.contactActivations.WithLabelValues(ca.Channel).Observe(ca.Duration)
}

// receive a *metric.ContactActivated metric and increment to a Gauge metric type.
func (s *Service) IncContactActivated(ca *ContactActivated) {
	s.contactActivated.WithLabelValues(ca.Channel).Inc()
}

// receive a *metric.ContactActivated metric and decrement to a Gauge metric type.
func (s *Service) DecContactActivated(ca *ContactActivated) {
	s.contactActivated.WithLabelValues(ca.Channel).Dec()
}
