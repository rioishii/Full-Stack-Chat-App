FROM alpine
COPY gateway /gateway
RUN apk add --no-cache ca-certificates
EXPOSE 80
ENTRYPOINT ["/gateway"]
RUN update-ca-certificates
