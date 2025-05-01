# TestCloud - Балансировщик нагрузки с Rate Limiting

Проект реализует:
- HTTP балансировщик нагрузки с алгоритмом **Round Robin**
- Проверку состояния (health-check) бекендов
- **Rate limiting** по IP или API-ключу на основе Token Bucket
- Конфигурацию через YAML
- Интеграцию с **Prometheus**
- Подключение к базе данных Postgres для хранения лимитов
- Docker-развёртыванием и миграциями
Быстрый старт

### 1. Склонируйте проект
[Смотреть демонстрацию на YouTube](https://www.youtube.com/watch?v=dQw4w9WgXcQ)

 ```bash
git clone https://github.com/yourname/TestCloud.git
make upDocker
make migrate_up
git clone https://github.com/yourname/TestCloud.git
make upDocker
make migrate_up
 ```


 ```bash
internal/
  service/
    backend/    # Реализация backend и reverse proxy
    balancer/   # Round-robin логика
    limiter/    # Token Bucket + Middleware
    proxy/      # Входная точка + маршрутизация
    health/     # Проверка доступности backend'ов
  config/       # Загрузка конфигурации
  server/       # HTTP сервер + инициализация
pkg/
  gp/           # Работа с БД (pgx pool)
  logger/          # Конфигурация zap-логгера
  metrics/      # Prometheus метрики
 ```
