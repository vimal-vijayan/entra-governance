package controller

import "time"

const (
	defaultRequeueDuration           = 10 * time.Minute
	faildStatusUpdateRequeueDuration = 10 * time.Second

	// Entra group constants
	entraSecurityGroupFinalizer = "finalizer.entraSecurityGroup.iam.entra.governance.com"

	// Entra app registration constants
	entraAppRegistrationFinalizer = "finalizer.entraAppRegistration.iam.entra.governance.com"
)
