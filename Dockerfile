FROM golang:1.8

# Set proxy for go get
#ENV http_proxy=http://pkg.proxy.prod.jp.local:10080
#ENV https_proxy=http://pkg.proxy.prod.jp.local:10080

#RUN groupadd -g CI_GID ci && useradd -u CI_UID -g CI_GID ci

#RUN chown ci /go