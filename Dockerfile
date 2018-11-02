# Based on centos
FROM centos:7.4.1708
LABEL maintainers="Kubernetes Authors"
LABEL description="Huawei Cloud Service Broker"

# Copy from build directory
COPY huaweicloud-service-broker /huaweicloud-service-broker

# Update
RUN yum -y update

# Define default command
ENTRYPOINT ["/huaweicloud-service-broker" ]
