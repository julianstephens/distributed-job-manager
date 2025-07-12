FROM node:24-alpine as base

ARG APP_HOME

FROM base as development

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . $APP_HOME

WORKDIR $APP_HOME

CMD ["pnpm", "dev"]
