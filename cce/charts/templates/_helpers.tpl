{{/* vim: set filetype=mustache: */}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- printf "%s-%s" .Release.Namespace .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "config" -}}
{{- printf (.Files.Get "files/config.json") .Values.brokerConfig.logLevel .Values.brokerConfig.username .Values.brokerConfig.password .Values.backDatabase.databaseType .Values.backDatabase.databaseHost .Values.backDatabase.databasePort .Values.backDatabase.databaseName .Values.backDatabase.databaseUsername .Values.backDatabase.databasePassword .Values.cloudCredentials.authUrl .Values.cloudCredentials.username .Values.cloudCredentials.password .Values.cloudCredentials.domainName .Values.cloudCredentials.tenantName .Values.cloudCredentials.region .Values.cloudCredentials.accessKey .Values.cloudCredentials.secretKey .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone .Values.catalog.vpcID .Values.catalog.subnetID .Values.catalog.securityGroupID .Values.catalog.primaryAvailabilityZone | b64enc -}}
{{- end -}}
