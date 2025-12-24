# PetsProject

Мини-сервис для управления питомцами (PetFeed):
- добавление питомцев
- получение питомца по ID
- получение списка всех питомцев

Проект поднимает всё окружение автоматически: Postgres, Redis и сервис PetFeed.

Поднять Postgres + Redis через Docker
```bash
make docker-up
```

Собрать проект
```bash
make build
```

Запустить сервис (переменные окружения берутся из docker-compose)
```bash
make run
```

Прогнать тесты

```bash
make test
```


Остановить окружение Docker

```bash 
make docker-down
```