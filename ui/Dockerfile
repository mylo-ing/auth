#######################################################################
#  Build stage – compile the Next.js App Router app                   #
#######################################################################
FROM node:24-alpine AS builder
WORKDIR /app

# 1) install dependencies
COPY ui/package.json ui/yarn.lock ./
RUN yarn install --frozen-lockfile

# 2) copy source and build
COPY ui .
RUN yarn build

#######################################################################
#  Runtime stage – tiny image that runs the compiled output          #
#######################################################################
FROM node:24-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

# copy only what we need
COPY --from=builder /app/.next/standalone .
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
