package main

import (
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-cf/cf-redis-broker/brokerconfig"
	"net/http"
)

type MyServiceBroker struct{
	InstanceCreators map[string]InstanceCreator
	InstanceBinders  map[string]InstanceBinder
	Config           brokerconfig.Config

}

type InstanceCredentials struct {
	Host     string
	Port     int
	Password string
}


type InstanceCreator interface {
	Create(instanceID string) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

type InstanceBinder interface {
	Bind(instanceID string, bindingID string) (InstanceCredentials, error)
	Unbind(instanceID string, bindingID string) error
	InstanceExists(instanceID string) (bool, error)
}

func (*MyServiceBroker) Services() []brokerapi.Service {
	// Return a []brokerapi.Service here, describing your service(s) and plan(s)
	free := false
	services := []brokerapi.Service{
		brokerapi.Service{
			ID:             "test_id",
			Name:           "test_redis",
			Description:	"test_desc",
			Bindable:        true,
			Tags:            []string{"tag1","tag2",},
			PlanUpdatable:   true,
			Plans:	[]brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          "test_plan_id",
					Name:        "small_plan",
					Description: "some_plan_description",
					Free:        &free,
					Metadata:    nil,
				},
			},
			Requires:        nil,
			Metadata:        nil,
			DashboardClient: nil,
		},
	}
	return services
}

func (*MyServiceBroker) Provision(
instanceID string,
details brokerapi.ProvisionDetails,
asyncAllowed bool,
) (brokerapi.ProvisionedServiceSpec, error) {

	return brokerapi.ProvisionedServiceSpec{},nil
}

func (*MyServiceBroker) LastOperation(instanceID string) (brokerapi.LastOperation, error) {
	// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
	// for the status of the provisioning operation.
	// This also applies to deprovisioning (work in progress).
	return brokerapi.LastOperation{},nil
}

func (*MyServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {
	// Deprovision a new instance here. If async is allowed, the broker can still
	// chose to deprovision the instance synchronously, hence the first return value.
	//for _, instanceCreator := range MyServiceBroker.InstanceCreators {
	//	instanceExists, _ := instanceCreator.InstanceExists(instanceID)
	//	if instanceExists {
	//		return instanceCreator.Destroy(instanceID)
	//	}
	//}
	return false, brokerapi.ErrInstanceDoesNotExist
}

func (*MyServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// Bind to instances here
	// Return a binding which contains a credentials object that can be marshalled to JSON,
	// and (optionally) a syslog drain URL.
	return brokerapi.Binding{},nil
}

func (*MyServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	// Unbind from instances here
	return nil
}

func (*MyServiceBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {
	// Update instance here
	return false,nil
}

func (redisServiceBroker *MyServiceBroker) instanceExists(instanceID string) bool {
	for _, instanceCreator := range redisServiceBroker.InstanceCreators {
		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
		if instanceExists {
			return true
		}
	}
	return false
}

func main() {
	serviceBroker := &MyServiceBroker{}
	logger := lager.NewLogger("my-service-broker")
	credentials := brokerapi.BrokerCredentials{
		Username: "admin",
		Password: "admin",
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":3000", nil)
}