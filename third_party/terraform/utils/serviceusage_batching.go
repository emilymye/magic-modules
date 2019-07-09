package google

import (
	"fmt"
	"time"
)

const (
	batchTypeServiceUsageEnableServices = "project/services:batchEnable"
)

// globalBatchEnableServices can be used to batch requests to enable services
// across resource nodes, i.e. to batch creation of several
// google_project_service(s) resources.
func globalBatchEnableServices(services []string, project string, config *Config) error {
	req := &BatchRequest{
		DebugId:      fmt.Sprintf("Enable Project Services %s: %+v", project, services),
		ResourceName: project,
		Body:         services,
		CombineF:     combineServiceUsageServicesBatches,
		SendF:        sendBatchFuncEnableServices(config),
	}

	_, err := config.requestBatcherServiceUsage.SendRequestWithTimeout(
		batchTypeServiceUsageEnableServices,
		req,
		time.Minute*10)
	return err
}

func combineServiceUsageServicesBatches(srvsRaw interface{}, toAddRaw interface{}) (interface{}, error) {
	srvs, ok := srvsRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", srvsRaw)
	}
	toAdd, ok := toAddRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected new request body type to be []string, got %v. This is a provider error.", toAdd)
	}

	return append(srvs, toAdd...), nil
}

func sendBatchFuncEnableServices(config *Config) batcherSendFunc {
	return func(project string, toEnableRaw interface{}) (interface{}, error) {
		toEnable, ok := toEnableRaw.([]string)
		if !ok {
			return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", toEnableRaw)
		}
		return nil, enableServiceUsageProjectServices(toEnable, project, config)
	}
}