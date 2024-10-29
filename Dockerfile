FROM openeuler/openeuler:24.03-lts as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=direct,https://proxy.golang.org

# build binary
WORKDIR /opt/source
COPY . .
RUN go build -a -o robot-universal-hook-dispatcher -buildmode=pie -ldflags "-s -linkmode 'external' -extldflags '- pie,-Wl,-z,now,relro,noexecstack'" .

# copy binary config and utils
FROM openeuler/openeuler:24.03-lts
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 robot && \
    useradd -u 1000 -g robot -s /bin/bash -m robot

USER robot

COPY --chown=robot --from=BUILDER /opt/source/robot-universal-hook-dispatcher /opt/app/robot-universal-hook-dispatcher

ENTRYPOINT /opt/app/robot-universal-hook-dispatcher
