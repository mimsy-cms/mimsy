FROM node:24-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable pnpm

COPY web /app/web/
WORKDIR /app/web

COPY pnpm-lock.yaml pnpm-workspace.yaml /app/

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS builder
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM base AS runner
COPY --from=prod-deps /app/node_modules /app/node_modules
COPY --from=prod-deps /app/web/node_modules /app/web/node_modules
COPY --from=builder /app/web/build /app/web/build
EXPOSE 3000

CMD ["node", "/app/build/index.js"]
