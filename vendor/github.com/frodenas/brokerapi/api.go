package brokerapi

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/frodenas/brokerapi/auth"
)

const provisionLogKey = "provision"
const updateLogKey = "update"
const deprovisionLogKey = "deprovision"
const bindLogKey = "bind"
const unbindLogKey = "unbind"
const lastOperationLogKey = "last-operation"

const instanceIDLogKey = "instance-id"
const bindingIDLogKey = "binding-id"
const provisionDetailsLogKey = "provision-details"
const updateDetailsLogKey = "update-details"
const deprovisionDetailsLogKey = "deprovision-details"
const bindDetailsLogKey = "bind-details"
const unbindDetailsLogKey = "unbind-details"

const invalidProvisionDetailsErrorKey = "invalid-provision-details"
const invalidUpdateDetailsErrorKey = "invalid-update-details"
const invalidBindDetailsErrorKey = "invalid-bind-details"

const instanceAlreadyExistsErrorKey = "instance-already-exists"
const instanceMissingErrorKey = "instance-missing"
const instanceLimitReachedErrorKey = "instance-limit-reached"
const instanceAsyncRequiredErrorKey = "instance-async-required"
const instanceNotUpdateableErrorKey = "instance-not-updateable"
const instanceNotBindableErrorKey = "instance-not-bindable"
const bindingAlreadyExistsErrorKey = "binding-already-exists"
const bindingMissingErrorKey = "binding-missing"
const bindingAppGUIDRequiredErrorKey = "binding-app-guid-required"
const unknownErrorKey = "unknown-error"

const statusUnprocessableEntity = 422

type BrokerCredentials struct {
	Username string
	Password string
}

func New(serviceBroker ServiceBroker, logger lager.Logger, brokerCredentials BrokerCredentials) http.Handler {
	router := newHTTPRouter()

	router.Get("/v2/catalog", catalog(serviceBroker, router, logger))

	router.Put("/v2/service_instances/{instance_id}", provision(serviceBroker, router, logger))
	router.Patch("/v2/service_instances/{instance_id}", update(serviceBroker, router, logger))
	router.Delete("/v2/service_instances/{instance_id}", deprovision(serviceBroker, router, logger))

	router.Put("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", bind(serviceBroker, router, logger))
	router.Delete("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", unbind(serviceBroker, router, logger))

	router.Get("/v2/service_instances/{instance_id}/last_operation", lastOperation(serviceBroker, router, logger))

	return wrapAuth(router, brokerCredentials)
}

func wrapAuth(router httpRouter, credentials BrokerCredentials) http.Handler {
	return auth.NewWrapper(credentials.Username, credentials.Password).Wrap(router)
}

func catalog(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		catalogResponse := serviceBroker.Services()

		respond(w, http.StatusOK, catalogResponse)
	}
}

func provision(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]
		acceptsIncomplete := false
		if req.URL.Query().Get("accepts_incomplete") == "true" {
			acceptsIncomplete = true
		}

		logger := logger.Session(provisionLogKey, lager.Data{
			instanceIDLogKey: instanceID,
		})

		var details ProvisionDetails
		if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
			logger.Error(invalidProvisionDetailsErrorKey, err)
			respond(w, http.StatusBadRequest, ErrorResponse{
				Description: err.Error(),
			})
			return
		}

		logger = logger.WithData(lager.Data{
			provisionDetailsLogKey: details,
		})

		provisioningResponse, asynch, err := serviceBroker.Provision(instanceID, details, acceptsIncomplete)
		if err != nil {
			switch err {
			case ErrInstanceAlreadyExists:
				logger.Error(instanceAlreadyExistsErrorKey, err)
				respond(w, http.StatusConflict, EmptyResponse{})
			case ErrInstanceLimitMet:
				logger.Error(instanceLimitReachedErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			case ErrAsyncRequired:
				logger.Error(instanceAsyncRequiredErrorKey, err)
				respond(w, statusUnprocessableEntity, ErrorResponse{
					Error:       "AsyncRequired",
					Description: err.Error(),
				})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		if asynch {
			respond(w, http.StatusAccepted, provisioningResponse)
			return
		}

		respond(w, http.StatusCreated, provisioningResponse)
	}
}

func update(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]
		acceptsIncomplete := false
		if req.URL.Query().Get("accepts_incomplete") == "true" {
			acceptsIncomplete = true
		}

		logger := logger.Session(updateLogKey, lager.Data{
			instanceIDLogKey: instanceID,
		})

		var details UpdateDetails
		if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
			logger.Error(invalidUpdateDetailsErrorKey, err)
			respond(w, http.StatusBadRequest, ErrorResponse{
				Description: err.Error(),
			})
			return
		}

		logger = logger.WithData(lager.Data{
			updateDetailsLogKey: details,
		})

		asynch, err := serviceBroker.Update(instanceID, details, acceptsIncomplete)
		if err != nil {
			switch err {
			case ErrInstanceDoesNotExist:
				logger.Error(instanceMissingErrorKey, err)
				respond(w, http.StatusInternalServerError, EmptyResponse{})
			case ErrAsyncRequired:
				logger.Error(instanceAsyncRequiredErrorKey, err)
				respond(w, statusUnprocessableEntity, ErrorResponse{
					Error:       "AsyncRequired",
					Description: err.Error(),
				})
			case ErrInstanceNotUpdateable:
				logger.Error(instanceNotUpdateableErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		if asynch {
			respond(w, http.StatusAccepted, EmptyResponse{})
			return
		}

		respond(w, http.StatusOK, EmptyResponse{})
	}
}

func deprovision(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]
		acceptsIncomplete := false
		if req.URL.Query().Get("accepts_incomplete") == "true" {
			acceptsIncomplete = true
		}

		logger := logger.Session(deprovisionLogKey, lager.Data{
			instanceIDLogKey: instanceID,
		})

		details := DeprovisionDetails{
			ServiceID: req.FormValue("service_id"),
			PlanID:    req.FormValue("plan_id"),
		}

		logger = logger.WithData(lager.Data{
			deprovisionDetailsLogKey: details,
		})

		asynch, err := serviceBroker.Deprovision(instanceID, details, acceptsIncomplete)
		if err != nil {
			switch err {
			case ErrInstanceDoesNotExist:
				logger.Error(instanceMissingErrorKey, err)
				respond(w, http.StatusGone, EmptyResponse{})
			case ErrAsyncRequired:
				logger.Error(instanceAsyncRequiredErrorKey, err)
				respond(w, statusUnprocessableEntity, ErrorResponse{
					Error:       "AsyncRequired",
					Description: err.Error(),
				})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		if asynch {
			respond(w, http.StatusAccepted, EmptyResponse{})
			return
		}

		respond(w, http.StatusOK, EmptyResponse{})
	}
}

func bind(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]
		bindingID := vars["binding_id"]

		logger := logger.Session(bindLogKey, lager.Data{
			instanceIDLogKey: instanceID,
			bindingIDLogKey:  bindingID,
		})

		var details BindDetails
		if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
			logger.Error(invalidBindDetailsErrorKey, err)
			respond(w, http.StatusBadRequest, ErrorResponse{
				Description: err.Error(),
			})
			return
		}

		logger = logger.WithData(lager.Data{
			bindDetailsLogKey: details,
		})

		bindingResponse, err := serviceBroker.Bind(instanceID, bindingID, details)
		if err != nil {
			switch err {
			case ErrInstanceDoesNotExist:
				logger.Error(instanceMissingErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			case ErrBindingAlreadyExists:
				logger.Error(bindingAlreadyExistsErrorKey, err)
				respond(w, http.StatusConflict, ErrorResponse{
					Description: err.Error(),
				})
			case ErrAppGUIDRequired:
				logger.Error(bindingAppGUIDRequiredErrorKey, err)
				respond(w, statusUnprocessableEntity, ErrorResponse{
					Error:       "RequiresApp",
					Description: err.Error(),
				})
			case ErrInstanceNotBindable:
				logger.Error(instanceNotBindableErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		respond(w, http.StatusCreated, bindingResponse)
	}
}

func unbind(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]
		bindingID := vars["binding_id"]

		logger := logger.Session(unbindLogKey, lager.Data{
			instanceIDLogKey: instanceID,
			bindingIDLogKey:  bindingID,
		})

		details := UnbindDetails{
			ServiceID: req.FormValue("service_id"),
			PlanID:    req.FormValue("plan_id"),
		}

		logger = logger.WithData(lager.Data{
			unbindDetailsLogKey: details,
		})

		if err := serviceBroker.Unbind(instanceID, bindingID, details); err != nil {
			switch err {
			case ErrInstanceDoesNotExist:
				logger.Error(instanceMissingErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			case ErrBindingDoesNotExist:
				logger.Error(bindingMissingErrorKey, err)
				respond(w, http.StatusGone, EmptyResponse{})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		respond(w, http.StatusOK, EmptyResponse{})
	}
}

func lastOperation(serviceBroker ServiceBroker, router httpRouter, logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := router.Vars(req)
		instanceID := vars["instance_id"]

		logger := logger.Session(lastOperationLogKey, lager.Data{
			instanceIDLogKey: instanceID,
		})

		lastOperationResponse, err := serviceBroker.LastOperation(instanceID)
		if err != nil {
			switch err {
			case ErrInstanceDoesNotExist:
				logger.Error(instanceMissingErrorKey, err)
				respond(w, http.StatusGone, EmptyResponse{})
			default:
				logger.Error(unknownErrorKey, err)
				respond(w, http.StatusInternalServerError, ErrorResponse{
					Description: err.Error(),
				})
			}
			return
		}

		respond(w, http.StatusOK, lastOperationResponse)
	}
}

func respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}
