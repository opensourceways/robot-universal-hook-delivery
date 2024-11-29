FROM openeuler/openeuler:24.03-lts as BUILDER
RUN dnf install -y golang && \
    dnf -y upgrade-minimal

ARG USER
ARG PASS
RUN echo "machine github.com login $USER password $PASS" > ~/.netrc

# build binary
WORKDIR /opt/source
COPY . .
RUN go build -a -o robot-universal-hook-delivery -buildmode=pie -ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'" .

# copy binary config and utils
FROM openeuler/openeuler:24.03-lts
RUN dnf -y upgrade-minimal && \
    dnf in -y shadow && \
    groupadd -g 1000 robot && \
    useradd -u 1000 -g robot -s /bin/bash -m robot

USER robot

COPY --chown=robot --from=BUILDER /opt/source/robot-universal-hook-delivery  /opt/app/robot-universal-hook-delivery

ENTRYPOINT ["/opt/app/robot-universal-hook-delivery", "--hmac-secret-file=/vault/secrets/gitcode-secret", "--enable-debug=true", "--handle-path=gitcode-hook", "--config-file=/vault/secrets/config"]
