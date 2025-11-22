package server

import (
	"os"

	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pro_interfaces"
)

func NewSubscriptionService(userRepo db.UserManager, optionsRepo db.OptionsManager) pro_interfaces.SubscriptionService {
	return &SubscriptionServiceImpl{}
}

type SubscriptionServiceImpl struct {
}

func (s *SubscriptionServiceImpl) GetToken() (res pro_interfaces.SubscriptionToken, err error) {
	err = db.ErrNotFound
	return
}

func (s *SubscriptionServiceImpl) HasActiveSubscription() bool {
	// Enable PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	return devMode == "true" || devMode == "1"
}

func (s *SubscriptionServiceImpl) CanAddProUser() (ok bool, err error) {
	// Enable PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	return devMode == "true" || devMode == "1", nil
}

func (s *SubscriptionServiceImpl) StartValidationCron() {

}

func (s *SubscriptionServiceImpl) CanAddRole() (ok bool, err error) {
	// Enable PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	return devMode == "true" || devMode == "1", nil
}

func (s *SubscriptionServiceImpl) CanAddRunner() (ok bool, err error) {
	// Enable PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	return devMode == "true" || devMode == "1", nil
}

func (s *SubscriptionServiceImpl) CanAddTerraformHTTPBackend() (ok bool, err error) {
	// Enable PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	return devMode == "true" || devMode == "1", nil
}

func (s *SubscriptionServiceImpl) GetPlan() (plan string, err error) {
	// Return dev plan in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	if devMode == "true" || devMode == "1" {
		return "dev", nil
	}
	return "", db.ErrNotFound
}
