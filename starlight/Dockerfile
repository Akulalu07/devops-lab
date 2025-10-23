FROM rust:1.90-trixie

WORKDIR /starlight
COPY . .
RUN cargo build --release
EXPOSE 8090
CMD ["./target/release/starlight"]
