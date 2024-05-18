# vk-oauth-proxy
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)
](https://hub.docker.com/r/afansv/vk-oauth-proxy)

VK Generic OAuth Proxy


VK использует немного нестандартные практики OAuth:

- `token_type` не возвращается при получении `access_token`
- `email` возвращается только при получении `access_token`, получить его через `users.get` уже не получится
- Ответы от `api.vk.com` оборачиваются в объект `response`, иногда это может быть неудобно

**vk-oauth-proxy** решает эти проблемы:

- В ответ на запрос `access_token` вместе с ним придёт и поле `token_type` со значением `Bearer`
- `email`, который приходит с `access_token` запоминается и ассоциируется с ID пользователя, поле `email` с этим
  значением будет добавлено в ответ `users.get` (хранится не на диске, просто кэш в рантайме с TTL – для большинства
  сценариев подойдёт)
- Поля объекта `response` из ответов `api.vk.com` будут разворачиваться в корень ответа (только для метода `users.get`)

Протестировано с версией VK API **5.199**

### Зачем?

Я не мог нормально подключить [ZITADEL](https://zitadel.com) к OAuth VK.
Потенциально может пригодиться для других SSO-сервисов и OAuth-клиентов.

### Конфигурация

Конфигурация происходит при помощи переменных окружения среды

| ENV                        | default                | Описание                       |
|----------------------------|------------------------|--------------------------------|
| `VOP_USER_EMAIL_STORE_TTL` | `1m`                   | Время хранения  email          |
| `VOP_OAUTH_UPSTREAM_HOST`  | `https://oauth.vk.com` | URL VK OAuth                   |
| `VOP_API_UPSTREAM_HOST`    | `https://api.vk.com`   | URL VK API                     |
| `VOP_OAUTH_PROXY_ADDR`     | `:9090`                | Адрес для запуска прокси OAuth |
| `VOP_API_PROXY_ADDR`       | `:9091`                | Адрес для запуска прокси API   |

### Docker
```shell
docker run -p 9090:9090 -p 9091:9091 afansv/vk-oauth-proxy:latest
```

### Полезные ссылки

- stackoverflow: ["OAuth2, Spring Авторизация через vk - ошибка invalid_token_response tokenType cannot be null"](https://ru.stackoverflow.com/questions/991083/oauth2-spring-Авторизация-через-vk-ошибка-invalid-token-response-tokentype-ca) 