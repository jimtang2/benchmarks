FROM node:22-alpine AS base
ARG SOURCE_PATH=express/node-postgres
WORKDIR /app
COPY express/package.json .
RUN npm install --omit=dev
COPY ${SOURCE_PATH}/server.js ./server.js
EXPOSE 8080
CMD ["npm", "start"]