FROM golang:1.7

# User iCE
RUN useradd -m ice

# Install testing dependencies
RUN curl -o /usr/local/bin/jq -L https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
  chmod +x /usr/local/bin/jq
RUN cd /opt && \
  git clone https://github.com/sstephenson/bats.git && \
  cd bats && \
  ./install.sh /usr/local
