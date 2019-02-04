package config

// Config holds all environment vars needed to run this service
type Config struct {
	Port                      int    `envconfig:"PORT" required:"true"`
	Environment               string `envconfig:"ENVIRONMENT" default:"development"`
	BrokerDSN                 string `envconfig:"BROKER_DSN" required:"true"`
	VerificationWorkers       int    `envconfig:"VERIFICATION_WORKERS" default:"5"`
	InvitationWorkers         int    `envconfig:"INVITATION_WORKERS" default:"5"`
	AlertWorkers              int    `envconfig:"ALERT_WORKERS" default:"10"`
	City                      string `envconfig:"CITY" required:"true"`
	Locale                    string `envconfig:"LOCALE" default:"en"`
	TwilioVerificationAPIHost string `envconfig:"TWILIO_VERIFICATION_API_HOST" required:"true"`
	TwilioVerificationAPIKey  string `envconfig:"TWILIO_VERIFICATION_API_KEY" required:"true"`
}
