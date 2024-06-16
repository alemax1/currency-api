# Currencies rate API

### 1. Requirements
+ Docker installed (to run docker-compose)
+ Make installed to run commands from Makefile

## 2. Getting started

### 2.1 Create environment:
    - Create .env file copying .example.env
    - Add custom 'CURRENCIES_API_KEY' to access forex API 

    ENV naming rules:
        m - minute
        s - second
        h - hour
        l - millisecond

### 2.2 Generate swagger docs:
```
make swagger-gen
```

### 2.3 Start App:
```
make compose
```