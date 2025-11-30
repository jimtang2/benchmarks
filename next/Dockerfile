FROM node:22-alpine AS base
WORKDIR /app
COPY next/package.json ./
RUN npm install --omit=dev
COPY next .
RUN npm run build
EXPOSE 8080
ENV PORT=8080
ENV NODE_ENV=production
CMD ["node", ".next/standalone/server.js"]