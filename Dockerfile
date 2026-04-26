# Официальный образ для сборки из Docker Hub
# Используется аргумент сборки для кросс-компиляции, хотя
# конкретно здесь он не нужен
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.26-alpine AS build

# Изоляция контекста сборки от файлов рабочей директории
# образа сборки
WORKDIR /app

# Скачивание зависимостей до копирования исходного кода,
# чтобы не инвалидировать кеш зависимостей
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Копирование исходного кода
COPY . .

# Аргументы для кросс-компиляции, которые BuildKit
# подставит автоматически
ARG TARGETOS
ARG TARGETARCH

# Сборка с использованием кеширования зависимостей
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /bin/app ./cmd/main.go

# Безопасный образ
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

# Копирование готового бинарника
COPY --from=build --chown=nonroot:nonroot /bin/app /bin/app

# Явное переключение на непривилегированного пользователя
USER nonroot

# Объявление стандартного для приложения TCP порта
EXPOSE 3000/tcp

# Указание точки входа
ENTRYPOINT [ "/bin/app" ]
