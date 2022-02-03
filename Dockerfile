FROM node:lts-alpine3.13 as ui-build
RUN apk add -U curl
COPY . /src
WORKDIR /src/ui
RUN npm install
RUN npm run build

FROM golang:1.17-alpine as app-build
RUN apk add -U make git
COPY . /src
WORKDIR /src
RUN make

FROM scratch as package
COPY --from=ui-build /src/ui/dist /public
COPY --from=app-build /src/bin/ /

FROM alpine:3.15 as server
COPY --from=app-build /src/bin/fynca-manager /usr/local/bin/fynca-manager
COPY --from=ui-build /src/ui/dist /var/fynca/public
WORKDIR /var/fynca
ENTRYPOINT ["/usr/local/bin/fynca-manager"]
