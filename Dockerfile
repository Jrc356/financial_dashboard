# Backend
FROM golang:1.18 AS go
WORKDIR /app
COPY server/go.mod server/go.sum /app/
RUN go mod download
COPY server /app

FROM go AS server-build
RUN CGO_ENABLED=0 go build -o server main.go

FROM go AS test-go
ENTRYPOINT [ "make", "test" ]

# Frontend
FROM node:18 as node
WORKDIR /app
COPY client/package.json client/package-lock.json /app/
RUN npm ci --include=dev
COPY client /app

FROM node as client-build
RUN npm run build # TODO: should remove devDependencies

FROM node as test-node
ENV CI=true
ENTRYPOINT [ "npm", "run", "test:coverage" ]

# Production
FROM gcr.io/distroless/static-debian11:nonroot as production
COPY --from=server-build /app/server /
COPY --from=client-build /app/build /client/build
EXPOSE 8080
ENTRYPOINT [ "/server" ]
