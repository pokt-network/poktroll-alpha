load('ext://restart_process', 'docker_build_with_restart')

# A list of directories where changes trigger a hot-reload of the sequencer
hot_reload_dirs = ['app', 'cmd']

# Run celestia node
k8s_yaml('localnet/kubernetes/celestia-rollkit.yaml')

# Build sequencer
local_resource('hot-reload: poktrolld', 'GOOS=linux ignite chain build --output=./bin --debug -v', deps=hot_reload_dirs)

# Build an image with a sequencer
docker_build_with_restart(
    "poktrolld",
    '.',
    dockerfile_contents="""FROM debian:bullseye
RUN apt-get -q update && apt-get install -qyy curl jq
COPY bin/poktrolld /usr/local/bin/poktrolld
WORKDIR /
""",
    only=["./bin/poktrolld"],
    entrypoint=[
        "/bin/sh", "/etc/config/poktroll.sh"
    ],
    live_update=[sync("bin/poktrolld", "/usr/local/bin/poktrolld")],
)

# Run poktrolld
k8s_yaml('localnet/kubernetes/poktrolld.yaml')
