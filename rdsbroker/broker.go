package rdsbroker

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/frodenas/brokerapi"
	"github.com/mitchellh/mapstructure"

	"crypto/tls"
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	trusts "github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/rds/v1/instance"
	netutil "k8s.io/apimachinery/pkg/util/net"
	certutil "k8s.io/client-go/util/cert"
)

const instanceIDLogKey = "instance-id"
const bindingIDLogKey = "binding-id"
const detailsLogKey = "details"
const acceptsIncompleteLogKey = "acceptsIncomplete"

type RDSInstance struct {
	instanceID      string
	rdsInstanceID   string
	rdsInstanceName string
}

type RDSBroker struct {
	IdentityEndpoint             string
	Ca                           string
	Username                     string
	Password                     string
	DomainName                   string
	ProjectName                  string
	ProjectID                    string
	Region                       string
	dbPrefix                     string
	allowUserProvisionParameters bool
	allowUserUpdateParameters    bool
	allowUserBindParameters      bool
	catalog                      Catalog
	logger                       lager.Logger
	rdsInstances                 []RDSInstance
}

func New(
	config Config,
	logger lager.Logger,
) *RDSBroker {
	return &RDSBroker{
		IdentityEndpoint: config.IdentityEndpoint,
		Ca:               config.Ca,
		Username:         config.Username,
		Password:         config.Password,
		DomainName:       config.DomainName,
		ProjectName:      config.ProjectName,
		ProjectID:        config.ProjectID,
		Region:           config.Region,
		//
		dbPrefix:                     config.DBPrefix,
		allowUserProvisionParameters: config.AllowUserProvisionParameters,
		allowUserUpdateParameters:    config.AllowUserUpdateParameters,
		allowUserBindParameters:      config.AllowUserBindParameters,
		catalog:                      config.Catalog,
		logger:                       logger.Session("broker"),
		rdsInstances:                 []RDSInstance{},
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
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("Provision: Prefixed instanceID:", instanceID)

	servicePlan, ok := b.catalog.FindServicePlan(details.PlanID)
	if !ok {
		return provisioningResponse, false, fmt.Errorf("Service Plan %s not found", details.PlanID)
	}

	rds_prop := servicePlan.RDSProperties
	fmt.Println("Provision:Getting Plan RDSProperties:", rds_prop)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("Creating RDS Client failed. Error:", err)
		return provisioningResponse, false, err
	}

	//dbList := rds.List(rdsClient)
	//fmt.Println("call db_list Provision call list", dbList)

	createResult := rds.Create(rdsClient, rds.CreateOps{
		Name:             instanceID,
		Datastore:        map[string]string{"type": rds_prop.DatastoreType, "version": rds_prop.DatastoreVersion},
		FlavorRef:        rds_prop.FlavorId,
		Volume:           map[string]interface{}{"type": rds_prop.VolumeType, "size": rds_prop.VolumeSize},
		Region:           rds_prop.Region,
		AvailabilityZone: rds_prop.AvailabilityZone,
		Vpc:              rds_prop.VpcId,
		Nics:             map[string]string{"subnetId": rds_prop.SubnetId},
		SecurityGroup:    map[string]string{"id": rds_prop.SecurityGroupId},
		DbPort:           rds_prop.Dbport,
		BackupStrategy:   map[string]interface{}{"startTime": rds_prop.BackupStrategyStarttime, "keepDays": rds_prop.BackupStrategyKeepdays},
		DbRtPd:           rds_prop.Dbpassword})
	if createResult.Err != nil {
		fmt.Println("Creating db intersence failed. Error:", createResult.Err)
		return provisioningResponse, false, createResult.Err
	}

	fmt.Println("call db_create Provision create_result:", createResult)

	dbInstance, err := createResult.Extract()
	if err != nil {
		fmt.Println("Converting db intersence failed. Error:", err)
		return provisioningResponse, false, err
	}

	fmt.Println("instance call create:", dbInstance)

	rdsInstance := RDSInstance{
		instanceID:      instanceID,
		rdsInstanceID:   dbInstance.ID,
		rdsInstanceName: dbInstance.Name,
	}
	b.rdsInstances = append(b.rdsInstances, rdsInstance)

	fmt.Println("rdsInstance call create", rdsInstance)
	fmt.Println("rdsInstance call create", b.rdsInstances)

	return provisioningResponse, true, nil
}

func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
	b.logger.Debug("update", lager.Data{
		instanceIDLogKey:        instanceID,
		detailsLogKey:           details,
		acceptsIncompleteLogKey: acceptsIncomplete,
	})
	fmt.Println("########Call update#########")
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("Update: Prefixed instanceID:", instanceID)
	fmt.Println("Update: rdsInstanceID:", b.rdsInstances)
	var rdsInstanceID string
	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("Update: rdsInstance:", v)
			rdsInstanceID = v.rdsInstanceID
		}
	}

	fmt.Println("Update: rdsInstanceID:", rdsInstanceID)
	fmt.Println("Update: details.RawParameters:", details.Parameters)

	updateParameters := UpdateParameters{}
	if err := mapstructure.Decode(details.Parameters, &updateParameters); err != nil {
		return false, fmt.Errorf("Update: Getting updateParameters failed. Error: %s", err)
	}

	fmt.Println("Update: updateParameters:", updateParameters)
	fmt.Println("Update: updateParameters.VolumeSize:", updateParameters.VolumeSize)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("Deprovision:Creating RDS Client failed. Error:", err)
		return false, fmt.Errorf("Creating RDS Client failed. Error: %s", err)
	}

	// Find the rdsInstanceID from rds service
	dbList := rds.List(rdsClient)
	fmt.Println("Update: call list: ", dbList)
	dbInstanceList, err := dbList.Extract()
	if err != nil {
		fmt.Println("Update: Converting db dbList failed. Error:", err)
		return false, fmt.Errorf("Update: Converting db dbList failed. Error: %s", err)
	}
	fmt.Println("Update: dbInstanceList", dbInstanceList)

	if rdsInstanceID == "" {
		fmt.Println("Update: Getting rds InstanceID from rdsInstances failed.", b.rdsInstances)
		for _, dbinstance := range dbInstanceList {
			var instanceName string = fmt.Sprintf("%s-%s", instanceID, "PostgreSQL-1")
			if dbinstance.Name == instanceName {
				rdsInstanceID = dbinstance.ID
				fmt.Println("Update: range rdsInstanceID:", rdsInstanceID)
			}
		}
	}

	fmt.Println("Update: rdsInstanceID from rds:", rdsInstanceID)
	// Find the rdsInstanceID from rds service end

	updateResult := rds.UpdateVolumeSize(rdsClient, rds.UpdateOps{
		Volume: map[string]interface{}{"size": updateParameters.VolumeSize},
	}, rdsInstanceID)
	fmt.Println("Update: updateResult", updateResult)
	if updateResult.Err != nil {
		fmt.Println("Update:Updating dbInstance failed. Error:", updateResult.Err)
		return false, fmt.Errorf("Update: Updating dbInstance failed. Error: %s", updateResult.Err)
	}

	newgetResult := rds.Get(rdsClient, rdsInstanceID)
	if newgetResult.Err != nil {
		fmt.Println("Getting dbInstance failed. Error:", newgetResult.Err)
		return false, fmt.Errorf("Getting dbInstance failed. Error: %s", newgetResult.Err)
	}

	fmt.Println("Update: newgetResult", newgetResult)

	return true, nil
}

func (b *RDSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
	b.logger.Debug("deprovision", lager.Data{
		instanceIDLogKey:        instanceID,
		detailsLogKey:           details,
		acceptsIncompleteLogKey: acceptsIncomplete,
	})

	fmt.Println("########Call Deprovision#########")
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("Deprovision: Prefixed instanceID:", instanceID)
	fmt.Println("Deprovision: rdsInstanceID:", b.rdsInstances)
	var rdsInstanceID string
	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("Deprovision: rdsInstance:", v)
			rdsInstanceID = v.rdsInstanceID
		}
	}

	fmt.Println("Deprovision: rdsInstanceID:", rdsInstanceID)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("Deprovision:Creating RDS Client failed. Error:", err)
		return false, fmt.Errorf("Creating RDS Client failed. Error: %s", err)
	}

	// Find the rdsInstanceID from rds service
	dbList := rds.List(rdsClient)
	fmt.Println("Deprovision: call list", dbList)
	dbInstanceList, err := dbList.Extract()
	if err != nil {
		fmt.Println("Deprovision: Converting db dbList failed. Error:", err)
		return false, fmt.Errorf("Deprovision: Converting db dbList failed. Error: %s", err)
	}
	fmt.Println("Deprovision: dbInstanceList", dbInstanceList)

	if rdsInstanceID == "" {
		fmt.Println("Deprovision: Getting rds InstanceID from rdsInstances failed.", b.rdsInstances)
		for _, dbinstance := range dbInstanceList {
			var instanceName string = fmt.Sprintf("%s-%s", instanceID, "PostgreSQL-1")
			if dbinstance.Name == instanceName {
				rdsInstanceID = dbinstance.ID
				fmt.Println("Deprovision: range rdsInstanceID:", rdsInstanceID)
			}
		}
	}

	fmt.Println("Deprovision: rdsInstanceID from rds:", rdsInstanceID)
	// Find the rdsInstanceID from rds service end

	getResult := rds.Delete(rdsClient, rdsInstanceID)
	if getResult.Err != nil {
		fmt.Println("Deprovision:Deleting dbInstance failed. Error:", getResult.Err)
		return false, fmt.Errorf("Deprovision:Deleting dbInstance failed. Error: %s", getResult.Err)
	}

	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("Deprovision: rdsInstance:", v)
			v.instanceID = ""
			v.rdsInstanceID = ""
		}
	}

	return true, nil
}

func (b *RDSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
	b.logger.Debug("bind", lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})
	fmt.Println("########Call Bind#########")
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("Bind: Prefixed instanceID:", instanceID)

	bindingResponse := brokerapi.BindingResponse{}
	service, ok := b.catalog.FindService(details.ServiceID)
	if !ok {
		return bindingResponse, fmt.Errorf("Service %s not found", details.ServiceID)
	}

	if !service.Bindable {
		return bindingResponse, brokerapi.ErrInstanceNotBindable
	}

	servicePlan, ok := b.catalog.FindServicePlan(details.PlanID)
	if !ok {
		return bindingResponse, fmt.Errorf("Service Plan %s not found", details.PlanID)
	}

	fmt.Println("Bind: Service", servicePlan)

	rds_prop := servicePlan.RDSProperties
	fmt.Println("Bind:Getting Plan RDSProperties:", rds_prop)

	fmt.Println("rdsInstanceID:", b.rdsInstances)
	var rdsInstanceID string
	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("Bind: rdsInstance:", v)
			rdsInstanceID = v.rdsInstanceID
		}
	}

	fmt.Println("rdsInstanceID:", rdsInstanceID)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("Creating RDS Client failed. Error:", err)
		return bindingResponse, fmt.Errorf("Creating RDS Client failed. Error: %s", err)
	}

	// Find the rdsInstanceID from rds service
	dbList := rds.List(rdsClient)
	fmt.Println("Bind: call list", dbList)
	dbInstanceList, err := dbList.Extract()
	if err != nil {
		fmt.Println("Bind: Converting db dbList failed. Error:", err)
		return bindingResponse, fmt.Errorf("Bind: Converting db dbList failed. Error: %s", err)
	}
	fmt.Println("Bind: dbInstanceList", dbInstanceList)

	if rdsInstanceID == "" {
		fmt.Println("Bind: Getting rds InstanceID from rdsInstances failed.", b.rdsInstances)
		for _, dbinstance := range dbInstanceList {
			var instanceName string = fmt.Sprintf("%s-%s", instanceID, "PostgreSQL-1")
			if dbinstance.Name == instanceName {
				rdsInstanceID = dbinstance.ID
				fmt.Println("Bind: range rdsInstanceID:", rdsInstanceID)
			}
		}
	}

	fmt.Println("Bind: rdsInstanceID from rds:", rdsInstanceID)
	// Find the rdsInstanceID from rds service end

	getResult := rds.Get(rdsClient, rdsInstanceID)
	if getResult.Err != nil {
		fmt.Println("Getting dbInstance failed. Error:", getResult.Err)
		return bindingResponse, fmt.Errorf("Getting dbInstance failed. Error: %s", getResult.Err)
	}

	fmt.Println("Bind: getResult", getResult)

	dbInstance, err := getResult.Extract()
	if err != nil {
		fmt.Println("Converting db intersence failed. Error:", err)
		return bindingResponse, fmt.Errorf("Converting db intersence failed. Error: %s", err)
	}

	fmt.Println("Bind: dbInstance", dbInstance)
	fmt.Println("Bind: dbInstance status", dbInstance.Status)
	var dbAddress, dbName, dbUsername, dbPassword string
	var dbPort int64
	dbAddress = dbInstance.HostName
	dbName = rds_prop.DbName
	dbUsername = rds_prop.Dbusername
	dbPassword = rds_prop.Dbpassword
	dbPort = dbInstance.DbPort

	var jdbcurl string
	jdbcurl = b.JDBCURI(dbAddress, dbPort, dbName, dbUsername, dbPassword)
	fmt.Println("Bind: jdbcurl", jdbcurl)

	bindingResponse.Credentials = &brokerapi.CredentialsHash{
		Host:     dbAddress,
		Port:     dbPort,
		Name:     dbName,
		Username: dbUsername,
		Password: dbPassword,
		URI:      b.URI(dbAddress, dbPort, dbName, dbUsername, dbPassword),
	}
	return bindingResponse, nil
}

func (b *RDSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	b.logger.Debug("unbind", lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})

	fmt.Println("########Call Unbind#########")
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("Unbind: Prefixed instanceID:", instanceID)
	servicePlan, ok := b.catalog.FindServicePlan(details.PlanID)
	if !ok {
		return fmt.Errorf("Service Plan %s not found", details.PlanID)
	}

	fmt.Println("Unbind : Service", servicePlan)

	fmt.Println("Unbind: rdsInstanceID:", b.rdsInstances)
	var rdsInstanceID string
	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("Unbind: rdsInstance:", v)
			rdsInstanceID = v.rdsInstanceID
		}
	}

	fmt.Println("Unbind: rdsInstanceID:", rdsInstanceID)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("Unbind: Creating RDS Client failed. Error:", err)
		return fmt.Errorf("Unbind: Creating RDS Client failed. Error: %s", err)
	}

	dbList := rds.List(rdsClient)
	fmt.Println("Unbind: call list", dbList)
	dbInstanceList, err := dbList.Extract()
	if err != nil {
		fmt.Println("Unbind: Converting db dbList failed. Error:", err)
		return fmt.Errorf("Unbind: Converting db dbList failed. Error: %s", err)
	}
	fmt.Println("Unbind: dbInstanceList", dbInstanceList)

	if rdsInstanceID == "" {
		fmt.Println("Unbind: Getting rds InstanceID from rdsInstances failed.", b.rdsInstances)
		for _, dbinstance := range dbInstanceList {
			var instanceName string = fmt.Sprintf("%s-%s", instanceID, "PostgreSQL-1")
			if dbinstance.Name == instanceName {
				rdsInstanceID = dbinstance.ID
				fmt.Println("Unbind: range rdsInstanceID:", rdsInstanceID)
			}
		}
	}

	fmt.Println("Unbind: rdsInstanceID:", rdsInstanceID)

	return nil
}

func (b *RDSBroker) LastOperation(instanceID string) (brokerapi.LastOperationResponse, error) {
	b.logger.Debug("last-operation", lager.Data{
		instanceIDLogKey: instanceID,
	})

	fmt.Println("########Call LastOperation#########")
	instanceID = fmt.Sprintf("%s-%s", b.dbPrefix, instanceID)
	fmt.Println("LastOperation: Prefixed instanceID:", instanceID)
	fmt.Println("b.rdsInstances:", b.rdsInstances)
	lastOperationResponse := brokerapi.LastOperationResponse{State: brokerapi.LastOperationFailed}

	var rdsInstanceID string
	for _, v := range b.rdsInstances {
		if v.instanceID == instanceID {
			fmt.Println("LastOperation: rdsInstance:", v)
			rdsInstanceID = v.rdsInstanceID
		}
	}

	fmt.Println("rdsInstanceID:", rdsInstanceID)

	rdsClient, err := b.NewRDSClientFromConfig()
	if err != nil {
		fmt.Println("LastOperation: Creating RDS Client failed. Error:", err)
		return lastOperationResponse, nil
	}

	// Find the rdsInstanceID from rds service
	dbList := rds.List(rdsClient)
	fmt.Println("LastOperation: call list", dbList)
	dbInstanceList, err := dbList.Extract()
	if err != nil {
		fmt.Println("LastOperation: Converting db dbList failed. Error:", err)
		return lastOperationResponse, fmt.Errorf("LastOperation: Converting db dbList failed. Error: %s", err)
	}
	fmt.Println("LastOperation: dbInstanceList", dbInstanceList)

	if rdsInstanceID == "" {
		fmt.Println("LastOperation: Getting rds InstanceID from rdsInstances failed.", b.rdsInstances)
		for _, dbinstance := range dbInstanceList {
			var instanceName string = fmt.Sprintf("%s-%s", instanceID, "PostgreSQL-1")
			if dbinstance.Name == instanceName {
				rdsInstanceID = dbinstance.ID
				fmt.Println("LastOperation: range rdsInstanceID:", rdsInstanceID)
			}
		}
	}
	fmt.Println("LastOperation: rdsInstanceID from rds:", rdsInstanceID)
	// Find the rdsInstanceID from rds service end

	getResult := rds.Get(rdsClient, rdsInstanceID)
	if getResult.Err != nil {
		fmt.Println("LastOperation: Getting dbInstance failed. Error:", getResult.Err)
		lastOperationResponse.State = brokerapi.LastOperationSucceeded
		return lastOperationResponse, nil
	}

	fmt.Println("LastOperation: getResult", getResult)

	dbInstance, err := getResult.Extract()
	if err != nil {
		fmt.Println("LastOperation: Converting db intersence failed. Error:", err)
		return lastOperationResponse, nil
	}

	fmt.Println("LastOperation: dbInstance", dbInstance)
	fmt.Println("LastOperation: dbInstance status", dbInstance.Status)

	if dbInstance.Status == "BUILD" {
		lastOperationResponse.State = brokerapi.LastOperationInProgress
	}
	if dbInstance.Status == "ACTIVE" || dbInstance.Status == "DELETED" {
		lastOperationResponse.State = brokerapi.LastOperationSucceeded
	}

	return lastOperationResponse, nil
}

func (b *RDSBroker) NewRDSClient() (*gophercloud.ServiceClient, error) {
	keystoneEndpoint := "https://iam.eu-de.otc.t-systems.com/v3"
	pc, err := openstack.NewClient(keystoneEndpoint)
	if err != nil {
		fmt.Println("Creating OpenStack provider failed. Error:", err)
		return nil, err
	}

	roots, err := certutil.NewPool("/root/ca/ca.crt")
	if err != nil {
		fmt.Println("Creating roots Ca failed. Error:", err)
		return nil, err
	}

	config := &tls.Config{}
	config.RootCAs = roots
	pc.HTTPClient.Transport = netutil.SetOldTransportDefaults(&http.Transport{TLSClientConfig: config})

	eo := gophercloud.EndpointOpts{Region: "eu-de", Availability: gophercloud.AvailabilityPublic}
	opts := tokens3.AuthOptions{
		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v3",
		Username:         "xxxxxx",
		Password:         "xxxxxx",
		//DomainID:         "",
		DomainName:  "xxxxxx",
		Scope:       tokens3.Scope{ProjectName: "eu-de", DomainName: "xxxxxx"},
		AllowReauth: true,
	}
	authOptsExt := trusts.AuthOptsExt{
		TrustID:            "",
		AuthOptionsBuilder: &opts,
	}

	err = openstack.AuthenticateV3(pc, authOptsExt, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Creating Keystone Auth failed. Error:", err)
		return nil, err
	}

	sc, err := openstack.NewRdsServiceV1(pc, eo, "xxxxxx")
	if err != nil {
		fmt.Println("Creating RDS Client failed. Error:", err)
		return nil, err
	}
	return sc, nil
}

func (b *RDSBroker) NewRDSClientFromConfig() (*gophercloud.ServiceClient, error) {

	keystoneEndpoint := b.IdentityEndpoint
	pc, err := openstack.NewClient(keystoneEndpoint)
	if err != nil {
		fmt.Println("Creating OpenStack provider failed. Error:", err)
		return nil, err
	}

	if b.Ca != "" {
		roots, err := certutil.NewPool(b.Ca)
		if err != nil {
			fmt.Println("Creating roots Ca failed. Error:", err)
			return nil, err
		}
		config := &tls.Config{}
		config.RootCAs = roots
		pc.HTTPClient.Transport = netutil.SetOldTransportDefaults(&http.Transport{TLSClientConfig: config})
	}

	eo := gophercloud.EndpointOpts{Region: b.Region, Availability: gophercloud.AvailabilityPublic}
	opts := tokens3.AuthOptions{
		IdentityEndpoint: b.IdentityEndpoint,
		Username:         b.Username,
		Password:         b.Password,
		//DomainID:         "",
		DomainName:  b.DomainName,
		Scope:       tokens3.Scope{ProjectName: b.ProjectName, DomainName: b.DomainName},
		AllowReauth: true,
	}
	authOptsExt := trusts.AuthOptsExt{
		TrustID:            "",
		AuthOptionsBuilder: &opts,
	}

	err = openstack.AuthenticateV3(pc, authOptsExt, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Creating Keystone Auth failed. Error:", err)
		return nil, err
	}

	sc, err := openstack.NewRdsServiceV1(pc, eo, b.ProjectID)
	if err != nil {
		fmt.Println("Creating RDS Client failed. Error:'", err)
		return nil, err
	}
	return sc, nil
}

func (d *RDSBroker) URI(address string, port int64, dbname string, username string, password string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?reconnect=true", username, password, address, port, dbname)
}

func (d *RDSBroker) JDBCURI(address string, port int64, dbname string, username string, password string) string {
	return fmt.Sprintf("jdbc:postgresql://%s:%d/%s?user=%s&password=%s", address, port, dbname, username, password)
}
