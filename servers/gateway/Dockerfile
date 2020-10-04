FROM alpine
COPY gateway /gateway
EXPOSE 443
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["/gateway"]