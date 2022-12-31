FROM ruby:2.7
ENV LANG C.UTF-8
WORKDIR /veracode-cli
COPY docker/base/Gemfile Gemfile
COPY docker/base/Gemfile.lock Gemfile.lock
COPY docker/base/.gemrc /root/.gemrc

# gem update because bundler 2 requires RubyGems 3
RUN apt update && \
    apt-get install -y golang-go && \
    export GOPATH=$HOME/go && \
    apt-get install -y docker.io && \
    go get github.com/mitchellh/gox && \
    export PATH=$PATH:$GOPATH/bin && \
    gox -h && \
    gem install bundler && \
    gem update --system && \
    BUNDLER_VERSION=2.0.1 bundle install && \
    bundle binstubs cucumber