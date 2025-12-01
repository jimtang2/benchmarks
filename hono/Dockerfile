FROM oven/bun:1-alpine
WORKDIR /app
COPY hono/package.json .
RUN bun install --production 
COPY hono/server.ts .
EXPOSE 8080
CMD ["bun", "run", "server.ts"]