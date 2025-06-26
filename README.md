# Simple IP Geolocation Service for OpenText Access Manager
This project provides a minimal, self-hosted Go application for IP-based geolocation. It is intended as a drop-in
replacement for the default geolocation provider used by OpenText Access Manager.

The application is designed for a containerized environment. As such, parameters such as the host, port, database
location, and the certificate and private key file names are not configurable by default. Feel free to open a pull
request if you'd like to contribute support for making these configurable.

## Examples
### Build app
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -tags 'netgo timetzdata' -o build/go-am-geolocation-provider cmd/app.go
```

### Build image
```shell
podman build --platform linux/amd64 -t localhost/go-am-geolocation-provider:latest ./build
```

### docker-compose.yml
```yaml
services:
  geoipupdate:
    container_name: geoipupdate
    image: ghcr.io/maxmind/geoipupdate:latest
    mem_limit: 256m
    userns_mode: auto
    environment:
      - GEOIPUPDATE_ACCOUNT_ID=
      - GEOIPUPDATE_LICENSE_KEY=
      - GEOIPUPDATE_EDITION_IDS=GeoLite2-Country
      - GEOIPUPDATE_FREQUENCY=72
    volumes:
      - ./app/data:/usr/share/GeoIP:z,U
  am-geolocation-provider:
    container_name: am-geolocation-provider
    image: localhost/go-am-geolocation-provider:latest
    depends_on:
      - geoipupdate
    mem_limit: 256m
    userns_mode: container:geoipupdate
    volumes:
      - ./app/tls:/app/tls:Z,U
      - ./app/data:/usr/share/GeoIP:z,U
    ports:
      - 5443:8443
```

Make sure that the certificate (cert.pem) and private key (key.pem) files are placed in the ./app/tls directory and the ./app/data directory exists.

## Access Manager configuration
To integrate this service with Access Manager, follow these steps in the Admin Console (5.1.x):

1. Navigate to **Risk-based Policies**.
2. Click the configuration icon and select **Geolocation Provider**.
3. For the Provider Type, choose the first provider.
4. Enter placeholder values for API Key and API Secret (these are not used).
5. Set the Web Service URL (including /\<ip\>) to the endpoint where the application is running.

## Disclaimer
All product names and brands mentioned in this project are the property of their respective owners. Use of this
software is entirely at your own risk. The authors and contributors accept no liability for any damage or loss resulting
from its use.