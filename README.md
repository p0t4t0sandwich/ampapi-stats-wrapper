# ampapi-stats-wrapper

[![License](https://img.shields.io/github/license/p0t4t0sandwich/ampapi-stats-wrapper?color=blue)](https://img.shields.io/github/downloads/p0t4t0sandwich/ampapi-stats-wrapper/LICENSE)
[![Github](https://img.shields.io/github/stars/p0t4t0sandwich/ampapi-stats-wrapper)](https://github.com/p0t4t0sandwich/ampapi-stats-wrapper)
[![Github Issues](https://img.shields.io/github/issues/p0t4t0sandwich/ampapi-stats-wrapper?label=Issues)](https://github.com/p0t4t0sandwich/ampapi-stats-wrapper/issues)
[![Discord](https://img.shields.io/discord/1067482396246683708?color=7289da&logo=discord&logoColor=white)](https://discord.neuralnexus.dev)
[![wakatime](https://wakatime.com/badge/github/p0t4t0sandwich/ampapi-stats-wrapper.svg)](https://wakatime.com/badge/github/p0t4t0sandwich/ampapi-stats-wrapper)

[![Github Releases](https://img.shields.io/github/downloads/p0t4t0sandwich/ampapi-stats-wrapper/total?label=Github&logo=github&color=181717)](https://github.com/p0t4t0sandwich/ampapi-stats-wrapper/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/p0t4t0sandwich/ampapi-stats-wrapper?label=Docker&logo=docker&color=2496ed)](https://hub.docker.com/r/p0t4t0sandwich/ampapi-stats-wrapper)

A simple wrapper built on the AMP API to expose API endpoints that return status responses

Ready to be used with [UpTimeRobot](https://uptimerobot.com/) or [Uptime Kuma](https://github.com/louislam/uptime-kuma)

Support:

- Ping `@thepotatoking3452` in the `#development` channel of the [AMP Discord](https://discord.gg/cubecoders)
- My own [development Discord](https://discord.neuralnexus.dev/)

## Installation

Download the latest release from the [releases page](https://github.com/p0t4t0sandwich/ampapi-stats-wrapper/releases) and extract it to a folder of your choice.

## Usage

### NOTE: The first time you ping an instance, it needs to authenticate with the main ADS, with Controller -> Target setups this can take 15 seconds or more, subsequent pings are much faster

### CLI

Run `ampapi-stats-wrapper.exe` (./ampapi-stats-wrapper on Linux) to start the wrapper. You can then access the API endpoints at `http://localhost:3021/`.

### Docker

#### Docker CLI

```bash
docker run -d \
    -p 3021:3021 \
    -e AMP_API_URL=http://localhost:8080 \
    -e AMP_API_USERNAME=admin \
    -e AMP_API_PASSWORD=myfancypassword123 \
    --name ampapi-stats-wrapper \
    p0t4t0sandwich/ampapi-stats-wrapper:latest
```

#### Docker Compose

```yaml
---
version: "3.8"
services:
  ampapi-stats-wrapper:
    image: p0t4t0sandwich/ampapi-stats-wrapper:latest
    container_name: ampapi-stats-wrapper
    environment:
      - TZ=UTC
      - IP_ADDRESS=0.0.0.0
      - PORT=3021
      - AMP_API_URL=http://localhost:8080
      - AMP_API_USERNAME=admin
      - AMP_API_PASSWORD=myfancypassword123
    ports:
      - "0.0.0.0:3021:3021"
    restart: unless-stopped
```

## Permissions

Within AMP, the user you use only needs the `manage` permission for each instance/target you want to monitor.
After changing permissions you'll need to restart the container/program for the changes to take effect.

## Configuration

The wrapper can be configured using the `settings.json` file. The following options are available:

| Option | Description | Default |
| --- | --- | --- |
| `ADDRESS` | The IP address and port to bind to | `0.0.0.0:3021` |
| `USE_UDS` | Whether to use a Unix Domain Socket (UDS) instead of a TCP socket | `false` |
| `AMP_API_URL` | The URL of the AMP API | `http://localhost:8080` |
| `AMP_API_USERNAME` | The username to use when authenticating with the AMP API | `admin` |
| `AMP_API_PASSWORD` | The password to use when authenticating with the AMP API | `myfancypassword123` |

## API Endpoints

You can read the OpenAPI documentation for the API endpoints at `http://0.0.0.0:3021/docs`

### GET /target/status/{TargetName}

Returns the status of the target defined by the schema `APICoreGetStatus`

### GET /instance/status/simple/{InstanceName}

Returns the simple status response of the instance: `Running` or `Offline`

### GET /server/status/{InstanceName}

Returns the status of the server defined by the schema `APICoreGetStatus`

### GET /server/status/simple/{InstanceName}

Returns the simple status response of the server defined by the schema `APIStateEnum`

- `Running` means the game server is active, the other possible responses can be found [here](https://github.com/p0t4t0sandwich/ampapi-spec/blob/main/useful_info/amp_status_enum.txt)
