# kata-ops

`kata-ops` is a tool that is used by ops teams for gathering information related to kata pods running in a k8s cluster.

## Building and publishing

Users can specify the following environment variables to the build:
* `KATA_OPS_REGISTRY` - The container registry to be used
- `KATA_OPS_TAG`      - A tag to the be used for the image
                          default: `$(git rev-parse HEAD)-$(uname -a)`
