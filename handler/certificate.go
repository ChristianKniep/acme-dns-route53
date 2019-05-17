package handler

import (
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/registration"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore"
	"github.com/begmaroman/acme-dns-route53/notifier"
)

// CertificateHandlerOptions is the options of certificate handler
type CertificateHandlerOptions struct {
	Staging           bool
	ConfigDir         string
	NotificationTopic string
	RenewBefore       int

	Store    certstore.CertStore
	Notifier notifier.Notifier
	DNS01    challenge.Provider

	Log *logrus.Logger
}

// CertificateHandler is the certificates handler
type CertificateHandler struct {
	isStaging         bool
	configDir         string
	notificationTopic string
	renewBefore       int

	store    certstore.CertStore
	notifier notifier.Notifier
	dns01    challenge.Provider
	log      *logrus.Logger
}

// NewCertificateHandler is the constructor of CertificateHandler
func NewCertificateHandler(opts *CertificateHandlerOptions) *CertificateHandler {
	return &CertificateHandler{
		isStaging:         opts.Staging,
		store:             opts.Store,
		notificationTopic: opts.NotificationTopic,
		renewBefore:       opts.RenewBefore,
		notifier:          opts.Notifier,
		dns01:             opts.DNS01,
		configDir:         opts.ConfigDir,
		log:               opts.Log,
	}
}

// toConfigParams creates a new configParams model
func (h *CertificateHandler) toConfigParams(user registration.User) *configParams {
	return &configParams{
		user:      user,
		isStaging: h.isStaging,
		keyType:   certcrypto.RSA2048, // TODO: Create a flag to define key type
	}
}

// toUserParams creates a new userParams model
func (h *CertificateHandler) toUserParams(email string) *userParams {
	return &userParams{
		email:   email,
		keyType: certcrypto.RSA2048, // TODO: Create a flag to define key type
	}
}
