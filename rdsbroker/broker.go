package rdsbroker

import (
	"encoding/json"
	"fmt"

	"github.com/frodenas/brokerapi"
	"github.com/mitchellh/mapstructure"
	"code.cloudfoundry.org/lager"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	trusts "github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/rds/v1/instance"

	//"github.com/chenyingkof/rds-broker/awsrds"
	//"github.com/chenyingkof/rds-broker/sqlengine"
	//"github.com/chenyingkof/rds-broker/utils"
)

const defaultUsernameLength = 16
const defaultPasswordLength = 32

const instanceIDLogKey = "instance-id"
const bindingIDLogKey = "binding-id"
const detailsLogKey = "details"
const acceptsIncompleteLogKey = "acceptsIncomplete"

type RDSBroker struct {
	dbPrefix                     string
	allowUserProvisionParameters bool
	allowUserUpdateParameters    bool
	allowUserBindParameters      bool
	catalog                      Catalog
	logger                       lager.Logger
}

func New(
	config Config,
	//dbInstance awsrds.DBInstance,
	//dbCluster awsrds.DBCluster,
	//sqlProvider sqlengine.Provider,
	logger lager.Logger,
) *RDSBroker {
	return &RDSBroker{
		dbPrefix:                     config.DBPrefix,
		allowUserProvisionParameters: config.AllowUserProvisionParameters,
		allowUserUpdateParameters:    config.AllowUserUpdateParameters,
		allowUserBindParameters:      config.AllowUserBindParameters,
		catalog:                      config.Catalog,
		//dbInstance:                   dbInstance,
		//dbCluster:                    dbCluster,
		//sqlProvider:                  sqlProvider,
		logger:                       logger.Session("broker"),
	}
}

func (b *RDSBroker) Services() brokerapi.CatalogResponse {
	fmt.Println("########Call Services#########")
	catalogResponse := brokerapi.CatalogResponse{}

	brokerCatalog, err := json.Marshal(b.catalog)
	if err != nil {
		b.logger.Error("marshal-error", err)
		return catalogResponse
	}

	apiCatalog := brokerapi.Catalog{}
	if err = json.Unmarshal(brokerCatalog, &apiCatalog); err != nil {
		b.logger.Error("unmarshal-error", err)
		return catalogResponse
	}

	catalogResponse.Services = apiCatalog.Services

	return catalogResponse
}

func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, acceptsIncomplete bool) (brokerapi.ProvisioningResponse, bool, error) {
	b.logger.Debug("provision", lager.Data{
		instanceIDLogKey:        instanceID,
		detailsLogKey:           details,
		acceptsIncompleteLogKey: acceptsIncomplete,
	})

	fmt.Println("########Call Provision#########")
	provisioningResponse := brokerapi.ProvisioningResponse{}


	//keystoneEndpoint := "https://iam.eu-de.otc.t-systems.com/v3" //TODO set the endpoint of keystone.
	//
	//pc, err := openstack.NewClient(keystoneEndpoint)
	//fmt.Println("pc '%s' not found", pc)
	//if err != nil {
	//	fmt.Println("error 01")
	//	return provisioningResponse, true, nil
	//}
	//eo := gophercloud.EndpointOpts{Region: "eu-de", Availability: "public"} //TODO change the two parameters
	//fmt.Println("eo '%s' not found", eo)
	//
	//err1 := openstack.AuthenticateV3(pc, tokens3.AuthOptions{Username: "", Password: "", DomainID: ""}, eo)
	//fmt.Println("err1 '%s' not found", err1)
	//
	//if err1 != nil {
	//	fmt.Println("error 02")
    	//	return provisioningResponse, true, nil
	//}
	//
	//sc, err2 := openstack.NewRdsServiceV1(pc, eo)
	//fmt.Println("sc '%s' not found", sc)
	//fmt.Println("err2 '%s' not found", err2)

	keystoneEndpoint := "https://iam.eu-de.otc.t-systems.com/v3" //TODO set the endpoint of keystone.
	pc, err := openstack.NewClient(keystoneEndpoint)
	if err != nil {
		fmt.Println("error 01")
		return provisioningResponse, true, nil
	}
	eo := gophercloud.EndpointOpts{Region: "eu-de", Availability: gophercloud.AvailabilityPublic} //TODO change the two parameters
	opts := tokens3.AuthOptions{
		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v3",
		Username:         "grab_tit",
		UserID:           "cb3cfa12219b47f5809e864b3d511ff5",
		Password:         "cnp200@HW",
		DomainID:         "",
		DomainName:       "grab_tit",
		AllowReauth:      true,
	}
	authOptsExt := trusts.AuthOptsExt{
		TrustID:            "", //TODO config the trust id
		AuthOptionsBuilder: &opts,
	}
	//authenticate
	err = openstack.AuthenticateV3(pc, authOptsExt, eo)

	if err != nil {
		fmt.Println("error 02")
		return provisioningResponse, true, nil
	}

	sc, _ := openstack.NewRdsServiceV1(pc, eo)
	fmt.Println("sc '%s' not found", sc)

	r := rds.Create(sc, rds.CreateOps{Name: "", FlavorRef: "", ReplicaOf: ""})
	if r.Err != nil {
		//TODO log the error
		fmt.Println("error 03")
		return provisioningResponse, true, nil
	}

	var instance rds.Instance
	err3 := r.ExtractInto(&instance)
	if err3 != nil {
		//TODO deal with the error
		fmt.Println("error 04")
		return provisioningResponse, true, nil
	}



	//r := instance.Create(sc, instance.CreateOps{Name: "", FlavorRef: "", ReplicaOf: ""})
	//if r.Err != nil {
	//	//TODO log the error
	//	fmt.Println("error 03")
	//	return provisioningResponse, true, nil
	//}

	//instance1 = instance.Instance{}
	//err3 := r.commonResult.ExtractInto(&instance1)
	//if err3 != nil {
	//	//TODO deal with the error
	//	fmt.Println("error 02")
	//	return provisioningResponse, true, nil
	//}


	return provisioningResponse, true, nil
}

func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
	b.logger.Debug("update", lager.Data{
		instanceIDLogKey:        instanceID,
		detailsLogKey:           details,
		acceptsIncompleteLogKey: acceptsIncomplete,
	})
	fmt.Println("########Call update#########")

	return true, nil
}

func (b *RDSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
	b.logger.Debug("deprovision", lager.Data{
		instanceIDLogKey:        instanceID,
		detailsLogKey:           details,
		acceptsIncompleteLogKey: acceptsIncomplete,
	})

	fmt.Println("########Call Deprovision#########")

	return true, nil
}

func (b *RDSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
	b.logger.Debug("bind", lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})
	fmt.Println("########Call Bind#########")

	bindingResponse := brokerapi.BindingResponse{}

	bindParameters := BindParameters{}
	if b.allowUserBindParameters {
		if err := mapstructure.Decode(details.Parameters, &bindParameters); err != nil {
			return bindingResponse, err
		}
	}

	service, ok := b.catalog.FindService(details.ServiceID)
	if !ok {
		return bindingResponse, fmt.Errorf("Service '%s' not found", details.ServiceID)
	}

	if !service.Bindable {
		return bindingResponse, brokerapi.ErrInstanceNotBindable
	}

	servicePlan, ok := b.catalog.FindServicePlan(details.PlanID)
	if !ok {
		return bindingResponse, fmt.Errorf("Service Plan '%s' not found", details.PlanID)
	}

	fmt.Println("Service '%s' not found", servicePlan)


	return bindingResponse, nil
}

func (b *RDSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	b.logger.Debug("unbind", lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})

	fmt.Println("########Call Unbind#########")

	servicePlan, ok := b.catalog.FindServicePlan(details.PlanID)
	if !ok {
		return fmt.Errorf("Service Plan '%s' not found", details.PlanID)
	}

	fmt.Println("Service '%s' not found", servicePlan)

	return nil
}

func (b *RDSBroker) LastOperation(instanceID string) (brokerapi.LastOperationResponse, error) {
	b.logger.Debug("last-operation", lager.Data{
		instanceIDLogKey: instanceID,
	})

	fmt.Println("########Call LastOperation#########")

	lastOperationResponse := brokerapi.LastOperationResponse{State: brokerapi.LastOperationFailed}


	return lastOperationResponse, nil
}
