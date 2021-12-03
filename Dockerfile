# Compile gcr binary on golang docker
FROM golang:latest
RUN git clone https://github.com/marekaf/gcr-lifecycle-policy.git /gcr
WORKDIR /gcr
RUN make dep \
    && make build

# Create container with kubectl and gcr for minimal footprint
FROM ubuntu:latest
COPY --from=0 /gcr/bin/gcr /usr/local/bin/gcr
RUN apt-get update \
    && apt-get install -y apt-transport-https gnupg2 curl \
    && curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
    && echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list \
    && apt-get update \
    && apt-get install -y kubectl
WORKDIR /root
ENV GOOGLE_APPLICATION_CREDENTIALS=/root/creds/serviceaccount.json
ENTRYPOINT ["gcr"]
CMD ["--help"]
