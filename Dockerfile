FROM golang:1.18 AS server-build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o server main.go

FROM node:18 as client-build
WORKDIR /app
COPY client /app
RUN npm install && npm run build

FROM gcr.io/distroless/static-debian11 as final
COPY --from=server-build /app/server /
COPY --from=client-build /app/build /client/build
EXPOSE 8080
ENTRYPOINT [ "/server" ]
