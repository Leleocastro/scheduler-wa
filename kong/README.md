# Endereços no Kong

- http://localhost:8000/ - gateway de entrada das apis
- http://localhost:8001/ - admin API
- http://localhost:8002/ - admin GUI

$ KONG_DATABASE=postgres docker compose --profile database up -d
