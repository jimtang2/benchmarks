ARG TARGET=axum

FROM rust:1.91-alpine AS builder
ARG TARGET
RUN apk add --no-cache musl-dev openssl-dev pkgconfig clang lld make
WORKDIR /app
COPY rust/Cargo.toml rust/Cargo.lock ./
RUN cargo fetch
COPY rust/src ./src
RUN cargo build --release --features ${TARGET} --bin ${TARGET}

FROM alpine:3.20
ARG TARGET
RUN apk add --no-cache ca-certificates libgcc
COPY --from=builder /app/target/release/${TARGET} /usr/local/bin/rust-bench
EXPOSE 8080
CMD ["rust-bench"]