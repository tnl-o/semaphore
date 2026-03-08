# feat: Semaphore 100% completion plan

## Изменения

- **WebSocket**: добавлена feature `ws` в axum, реализован `websocket_handler` с WebSocketManager
- **Task execution**: интеграция LocalJob в `execute_task`, `prepare_run` создаёт AnsibleApp/TerraformApp/ShellApp
- **Exporter**: убраны `unimplemented!`, реализованы `get_type_exporter`, `get_loaded_keys`, `get_loaded_keys_int`
- **API**: маршруты backup, invites, refs, integration aliases; OIDC flow; project_invite token
- **ShellApp**: кроссплатформенная команда-заглушка для тестов (cmd/sh)
- **apps**: уточнён `delete_app` для config-based приложений

## Тесты

Все 501 unit-тест проходят (`cargo test --lib`).
