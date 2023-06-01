# --- Base ----
FROM golang:1.18-buster AS base
WORKDIR $GOPATH/src/gitlab.neoway.com.br/sebrae-sp-ddm-datacapture-v2

# ---- Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.neoway.com.br
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN	go mod tidy
RUN	go mod vendor

# ---- Test ----
FROM dependencies AS test
COPY . .
RUN go install github.com/axw/gocov/gocov@latest
RUN go install github.com/matm/gocov-html/cmd/gocov-html@latest
RUN go get -d golang.org/x/oauth2/internal@v0.0.0-20211104180415-d3ed0bb246c8
RUN go get -d google.golang.org/grpc/internal/transport@v1.43.0
RUN go get -d google.golang.org/grpc@v1.43.0
ARG POSTGRES_URL
ARG PUB_SUB_EMULATOR_HOST
RUN POSTGRES_URL=$POSTGRES_URL PUB_SUB_EMULATOR_HOST=$PUB_SUB_EMULATOR_HOST go test -v -cpu 1 -failfast -coverprofile=coverage.out -covermode=set ./...
RUN gocov convert coverage.out | gocov-html > /index.html
RUN grep -v "_mock" coverage.out >> filtered_coverage.out
RUN go tool cover -func filtered_coverage.out

# ---- Lint ----
FROM dependencies AS lint
COPY . .
ENV GOLANG_CI_LINT_VERSION=v1.49.0
ENV GOFLAGS="-mod=mod"
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANG_CI_LINT_VERSION} && golangci-lint run -c ./.golangci.yml

# ---- audit ----
FROM dependencies AS audit
ENV NANCY_VERSION=v1.0.15
COPY go.mod .
RUN wget https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-linux-amd64 -qO /bin/nancy  && chmod +x /bin/nancy
RUN go list -m all | nancy sleuth

# ---- Build ----
FROM dependencies AS build
COPY . .
ARG VERSION
ARG BUILD
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/sebrae-sp-ddm-datacapture-v2  -ldflags "-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" ./cmd

# --- Release ----
FROM debian:stable AS image
# Downloading tools for install gcloud package
RUN apt-get update && \
    apt-get -y --no-install-recommends install curl \
        ca-certificates \
        python3
#Fix for google-cloud-sdk 427.0.0 python2 removal
ENV CLOUDSDK_PYTHON /usr/bin/python3
# Downloading gcloud package
RUN curl https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz > /tmp/google-cloud-sdk.tar.gz
# Installing the package
RUN mkdir -p /usr/local/gcloud \
  && tar -C /usr/local/gcloud -xvf /tmp/google-cloud-sdk.tar.gz \
  && /usr/local/gcloud/google-cloud-sdk/install.sh
# Adding the package path to local
ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin
# Copying Golang Build
COPY --from=build /go/bin/sebrae-sp-ddm-datacapture-v2 /sebrae-sp-ddm-datacapture-v2
# Copying Scripts files
COPY ./scripts/queries/*.sql /scripts/queries/
ENTRYPOINT ["/sebrae-sp-ddm-datacapture-v2"]