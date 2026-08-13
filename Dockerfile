# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

# ---- Stage 1: build the frontend ----------------------------------------
# Produces ui/dist, which is embedded into the Go binary via go:embed.
FROM node:24-bookworm-slim@sha256:235600a8101ab264e117b1768e925532262668dc9b581ef1dd7d96ced463b8e7 AS ui
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
WORKDIR /ui

# Install deps first so this layer caches unless the lockfile changes.
# .yarnrc.yml is required here, not optional: it sets nodeLinker: node-modules,
# without which yarn 4 defaults to PnP and vite/tsc resolve differently than
# they do locally. The yarn version comes from package.json's packageManager
# field, so corepack activates whatever is pinned there.
COPY ui/package.json ui/yarn.lock ui/.yarnrc.yml ./
RUN corepack enable && yarn install --immutable

# Build the SPA (tsc -b && vite build) -> /ui/dist
COPY ui/ ./
RUN yarn build


# ---- Stage 2: build the Go binary ---------------------------------------
# CGO is disabled so the result is a fully static binary for a scratch/distroless
# runtime. Version metadata is injected via -ldflags.
FROM golang:1.26-bookworm@sha256:1ecb7edf62a0408027bd5729dfd6b1b8766e578e8df93995b225dfd0944eb651 AS build
WORKDIR /src

# Download modules first for layer caching.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Application source.
COPY . .

# Drop in the freshly built UI so go:embed picks up the latest assets.
COPY --from=ui /ui/dist ./ui/dist

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath \
      -ldflags "-s -w \
        -X github.com/glueops/waggle/internal/buildinfo.Version=${VERSION} \
        -X github.com/glueops/waggle/internal/buildinfo.Commit=${COMMIT} \
        -X github.com/glueops/waggle/internal/buildinfo.Date=${DATE}" \
      -o /out/waggle .


# ---- Stage 3: minimal runtime image -------------------------------------
# Alpine gives a shell so you can `docker exec <container> waggle migrate up`
# (and run worker/encrypt/etc.) interactively. ca-certificates is required for
# outbound TLS (Proxmox, SMTP); the binary is static (CGO disabled) so it runs
# fine on musl.
FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS runtime

# postgresql-client provides `psql`, which the `waggle psql` subcommand execs.
RUN apk add --no-cache ca-certificates tzdata postgresql-client \
    && addgroup -S waggle && adduser -S -G waggle waggle

# Bind to all interfaces and serve the embedded UI out of the box.
ENV BIND_HOST=0.0.0.0 \
    BIND_PORT=8080 \
    FRONTEND_MODE=embed

COPY --from=build /out/waggle /usr/local/bin/waggle

EXPOSE 8080
USER waggle

# waggle is on PATH, so `docker exec <container> waggle migrate up` works.
ENTRYPOINT ["/usr/local/bin/waggle"]
# Override with `worker`, `migrate`, etc. at `docker run`.
CMD ["serve"]
