# Run image with this command
# docker run --rm --name bedrock --net host -v $(pwd)/config.toml:/app/config.toml -v $(pwd)/resources:/resources bedrock

# Building Frontend
FROM node:18-alpine as bedrock-web
WORKDIR /source
COPY . .
WORKDIR /source/packages/bedrock-web
RUN rm -rf dist node_modules
RUN --mount=type=cache,target=/source/packages/bedrock-web/node_modules,id=bedrock_web_modules_cache,sharing=locked \
    --mount=type=cache,target=/root/.npm,id=bedrock_web_node_cache \
    yarn install
RUN --mount=type=cache,target=/source/packages/bedrock-web/node_modules,id=bedrock_web_modules_cache,sharing=locked \
    yarn run build-only
RUN mv /source/packages/bedrock-web/dist /dist

# Building Backend
FROM golang:alpine as bedrock-server

WORKDIR /source
COPY . .
COPY --from=bedrock-web /dist /source/packages/bedrock-web/dist
RUN mkdir /dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/server .

# Runtime
FROM golang:alpine

COPY --from=bedrock-server /dist/server /app/server

EXPOSE 9443

CMD ["/app/server"]