FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21 as builder

RUN apk add --no-cache -U \
    libc-dev curl nodejs npm git gcc zip unzip tar

WORKDIR /usr/local
# hadolint ignore=DL4006
RUN curl -sL https://taskfile.dev/install.sh | sh

WORKDIR /go/src/semaphore
COPY . /go/src/semaphore

# Конвертируем server-wrapper в builder stage, чтобы избежать проблем с CRLF
RUN if [ -f /go/src/semaphore/deployment/docker/server/server-wrapper ]; then \
        python3 -c "import sys; data = open('/go/src/semaphore/deployment/docker/server/server-wrapper', 'rb').read(); open('/go/src/semaphore/deployment/docker/server/server-wrapper', 'wb').write(data.replace(b'\r\n', b'\n').replace(b'\r', b'\n'))" && \
        chmod +x /go/src/semaphore/deployment/docker/server/server-wrapper; \
    fi

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG APP_BUILD_TYPE
ARG TARGETOS
ARG TARGETARCH
ARG GH_TOKEN

RUN if [ -n "$APP_BUILD_TYPE" ]; then \
      git clone https://${GH_TOKEN}@github.com/semaphoreui/semaphorepro-module.git pro_impl && \
      go work init . ./pro_impl; \
    fi

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task deps APP_BUILD_TYPE=${APP_BUILD_TYPE} && \
    task build GOOS=${TARGETOS} GOARCH=${TARGETARCH} APP_BUILD_TYPE=${APP_BUILD_TYPE}


ENV OPENTOFU_VERSION="1.9.0"
ENV TERRAFORM_VERSION="1.11.3"
ENV TERRAGRUNT_VERSION="0.78.0"
#ENV PULUMI_VERSION="3.116.1"
#ENV POWERSHELL_VERSION="3.116.1"

RUN wget https://github.com/opentofu/opentofu/releases/download/v${OPENTOFU_VERSION}/tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz && \
    tar xf tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz -C /tmp && \
    rm tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz

RUN wget -O terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip || \
    (echo "Warning: Terraform ${TERRAFORM_VERSION} not found, skipping..." && touch /tmp/terraform) && \
    (test -f terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip && unzip terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip -d /tmp && rm terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip || true)

RUN wget -O /tmp/terragrunt https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_${TARGETARCH} && \
    chmod +x /tmp/terragrunt

FROM alpine:3.21

ARG TARGETARCH="amd64"
# renovate: datasource=pypi depName=ansible
ARG ANSIBLE_VERSION=11.1.0
ENV ANSIBLE_VERSION=${ANSIBLE_VERSION}
ARG ANSIBLE_VENV_PATH=/opt/semaphore/apps/ansible/${ANSIBLE_VERSION}/venv

RUN apk add --no-cache -U \
    bash curl git gnupg mysql-client openssh-client-default python3 py3-pip rsync sshpass tar tini tzdata unzip wget zip jq dos2unix && \
    rm -rf /var/cache/apk/* && \
    adduser -D -u 1001 -G root semaphore && \
    mkdir -p /tmp/semaphore && \
    mkdir -p /etc/semaphore && \
    mkdir -p /var/lib/semaphore && \
    mkdir -p /opt/semaphore && \
    chown -R semaphore:0 /tmp/semaphore && \
    chown -R semaphore:0 /etc/semaphore && \
    chown -R semaphore:0 /var/lib/semaphore && \
    chown -R semaphore:0 /opt/semaphore && \
    find /usr/lib/python* -iname __pycache__ | xargs rm -rf

RUN echo $'Host *\n  StrictHostKeyChecking no\n  UserKnownHostsFile /dev/null' > /etc/ssh/ssh_config.d/semaphore.conf

COPY --chown=1001:0 ./deployment/docker/server/ansible.cfg /etc/ansible/ansible.cfg
# Копируем server-wrapper из builder stage, где он уже конвертирован
COPY --from=builder --chown=1001:0 --chmod=755 /go/src/semaphore/deployment/docker/server/server-wrapper /usr/local/bin/server-wrapper
COPY --from=builder --chmod=755 /go/src/semaphore/bin/semaphore /usr/local/bin/semaphore
COPY --from=builder --chmod=755 /tmp/tofu /usr/local/bin/tofu
COPY --from=builder --chmod=755 /tmp/terraform /usr/local/bin/terraform
COPY --from=builder --chmod=755 /tmp/terragrunt /usr/local/bin/terragrunt

RUN chown semaphore:root /usr/local/bin/server-wrapper && \
    chown semaphore:root /usr/local/bin/semaphore && \
    chmod +x /usr/local/bin/server-wrapper && \
    test -f /usr/local/bin/server-wrapper && \
    test -x /usr/local/bin/server-wrapper && \
    head -1 /usr/local/bin/server-wrapper | grep -q '#!/bin/sh' || (echo "ERROR: Invalid shebang" && head -1 /usr/local/bin/server-wrapper && exit 1)

WORKDIR /home/semaphore

RUN apk add --no-cache -U python3-dev build-base openssl-dev libffi-dev cargo && \
     mkdir -p ${ANSIBLE_VENV_PATH} && \
     python3 -m venv ${ANSIBLE_VENV_PATH} --system-site-packages && \
     source ${ANSIBLE_VENV_PATH}/bin/activate && \
     pip3 install --upgrade pip ansible==${ANSIBLE_VERSION} boto3 botocore requests pywinrm passlib && \
     apk del python3-dev build-base openssl-dev libffi-dev cargo && \
     rm -rf /var/cache/apk/* && \
     find ${ANSIBLE_VENV_PATH} -iname __pycache__ | xargs rm -rf && \
     chown -R semaphore:0 /opt/semaphore

USER 1001
EXPOSE 3000

ENV VIRTUAL_ENV="$ANSIBLE_VENV_PATH"
ENV PATH="$ANSIBLE_VENV_PATH/bin:$PATH"

# Preventing ansible zombie processes. Tini kills zombies.
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/bin/sh", "/usr/local/bin/server-wrapper"]
