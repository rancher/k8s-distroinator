# rke

Rancher Kubernetes Engine, an extremely simple, lightning fast Kubernetes installer that works everywhere.

## Download

Please check the [releases](https://github.com/rancher/rke/releases/) page.

## Requirements

Please review the [OS requirements](https://rancher.com/docs/rke/latest/en/os/) for each node in your Kubernetes cluster.

## Getting Started

Please refer to our [RKE docs](https://rancher.com/docs/rke/latest/en/) for information on how to get started!
For cluster config examples, refer to [RKE cluster.yml examples](https://rancher.com/docs/rke/latest/en/example-yamls/)

## Installing Rancher HA using rke

Please use [High Availability (HA) Install](https://rancher.com/docs/rancher/v2.x/en/installation/ha/) to install Rancher in a high-availability configuration.

## Building

RKE can be built using the `make` command, and will use the scripts in the `scripts` directory as subcommands. The default subcommand is `ci` and will use `scripts/ci`. Cross compiling can be enabled by setting the environment variable `CROSS=1`. The compiled binaries can be found in the `build/bin` directory. Dependencies are managed by Go modules and can be found in [go.mod](https://github.com/rancher/rke/blob/master/go.mod).

## License

Copyright (c) 2019 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
