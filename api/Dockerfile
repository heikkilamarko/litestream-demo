FROM golang AS build

COPY . .

ENV GOPATH=""
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Statically compiling Go programs:
# https://www.arp242.net/static-go.html
RUN go build -ldflags '-s -w -extldflags "-static"' -tags osusergo,netgo,sqlite_omit_load_extension -o /usr/local/bin/api ./cmd/api

# Download the litestream binary:
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.9/litestream-v0.3.9-linux-amd64-static.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz

FROM busybox

COPY --from=build /usr/local/bin/litestream /usr/local/bin/
COPY --from=build /usr/local/bin/api /usr/local/bin/
COPY litestream.yml /etc/
COPY entrypoint.sh /

# SQLite database:
COPY items.db /var/api/data/

ENTRYPOINT ["/entrypoint.sh"]
