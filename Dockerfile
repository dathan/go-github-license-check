# this is a multi-stage build
# the 1st FROM builds the application
FROM golang:1.14-alpine3.11 AS baseGo
  LABEL stage=build
  #private repos are not used
  #ARG GITHUB_SSH_PRIV_KEY
  #RUN test -n "$GITHUB_SSH_PRIV_KEY"

  ENV CGO_ENABLED 0
  RUN apk add --no-cache git bash openssh

  # Grab Private Repo credentials
  #RUN mkdir /root/.ssh/
  #RUN echo "${GITHUB_SSH_PRIV_KEY}" > /root/.ssh/id_rsa
  #RUN chmod 400 /root/.ssh/id_rsa
  #RUN ssh-keyscan -H github.com >> ~/.ssh/known_hosts
  RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
  RUN mkdir /root/gocode
  COPY . /root/gocode
  WORKDIR /root/gocode

  ENV CGO_ENABLED 0
  RUN apk --no-cache add git bzr mercurial make
  ENV GO111MODULE on
  RUN go version
  RUN make build
  # turn off make test until we can mock google services  
  #RUN make test

#
# this is a dependancy for our container to have CA certs to talk to vault
#
FROM alpine:latest as alpineCerts
  LABEL stage=alpineCerts
  RUN apk add -U --no-cache ca-certificates

#
# the 2nd FROM says copy the binary from baseGo and put it here using scratch as its base
# - note using alpine because need to run a shell command wrapper
#FROM scratch AS release
FROM alpine:3.11.3 AS release
  LABEL stage=release

  # Pull CA Certificates to allow for TLS validation a CA
  COPY --from=alpineCerts /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
  COPY --from=baseGo /root/gocode/bin /app
  COPY --from=baseGo /root/gocode/scripts /app
  ENTRYPOINT ["/app/entrypoint.sh"]


#
# add another base image from scratch and add meta data called stage=mock
#
FROM scratch AS mock
  LABEL stage=mock

#
# Only ship the release layer
FROM release
