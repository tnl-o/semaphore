package features

import (
	"os"

	"github.com/semaphoreui/semaphore/db"
)

func GetFeatures(user *db.User, plan string) map[string]bool {
	// Enable all PRO features in development mode
	devMode := os.Getenv("SEMAPHORE_DEV_MODE")
	enablePro := devMode == "true" || devMode == "1"

	return map[string]bool{
		"project_runners":   enablePro,
		"terraform_backend": enablePro,
		"task_summary":      enablePro,
		"secret_storages":   enablePro,
	}
}
