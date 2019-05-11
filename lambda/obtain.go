package lambda

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore/acmstore"
	"github.com/begmaroman/acme-dns-route53/handler"
	"github.com/begmaroman/acme-dns-route53/notifier/awsns"
)

const (
	// ConfigDir is the default configuration directory
	ConfigDir = "/tmp"
)

var (
	// ErrEmailMissing is the error when email is not provided
	ErrEmailMissing = errors.New("email must be filled")

	// ErrDomainsMissing is the error when the domains list is empty
	ErrDomainsMissing = errors.New("domains list must not be filled")
)

// Payload contains payload data
type Payload struct {
	Domains []string `json:"domains"`
	Email   string   `json:"email"`
	Staging string   `json:"staging"`
}

func HandleLambdaEvent(payload Payload) error {
	conf := InitConfig(payload)

	// Domains list must not be empty
	if len(conf.Domains) == 0 {
		return ErrDomainsMissing
	}

	// Email must be filled
	if len(conf.Email) == 0 {
		return ErrEmailMissing
	}

	// Create options
	certificateHandlerOpts := &handler.CertificateHandlerOptions{
		ConfigDir: ConfigDir,
		Staging:   conf.Staging,
		Log:       logrus.New(),
		SNS:       awsns.New(sns.New(AWSSession)), // Initialize SNS API client
		R53:       route53.New(AWSSession),        // Initialize Route53 API client
		Store:     acmstore.New(AWSSession),       // Initialize ACM client
	}

	// Create a new handler
	certificateHandler := handler.NewCertificateHandler(certificateHandlerOpts)

	for _, domain := range conf.Domains {
		if err := certificateHandler.Obtain([]string{domain}, conf.Email); err != nil {
			logrus.Errorf("[%s] unable to obtain certificate: %s", domain, err)
			continue
		}

		logrus.Infof("[%s] certificate successfully obtained and stored", domain)
	}

	return nil
}
