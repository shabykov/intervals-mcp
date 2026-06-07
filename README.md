# intervals-icu MCP server (локальный, stdio)

MCP-сервер на Go, обёртка над API intervals.icu. Запускается локально и
подключается к **Claude Desktop** по stdio. Даёт Claude инструменты для
выгрузки тренировок, wellness и профиля — анализ делает уже сам Claude.

## Инструменты

- `list_activities` (oldest, newest) — список заездов с метриками (нагрузка, мощность, пульс, CTL/ATL)
- `get_activity` (id, intervals) — полные данные заезда, опционально с разбивкой по интервалам
- `get_activity_streams` (id, types) — посекундные ряды (watts/heartrate/cadence/…)
- `get_wellness` (oldest, newest) — CTL/ATL/форма (TSB), HRV, пульс покоя, вес
- `get_athlete` — профиль: FTP, зоны, настройки

## Требования

- Go 1.23+ (go.mod закреплён на SDK v1.3.0). Если у тебя Go 1.25+ — можешь
  обновиться на свежий SDK: `go get -u github.com/modelcontextprotocol/go-sdk`.
- API-ключ и Athlete ID из intervals.icu → Settings → Developer Settings.

## Сборка

```bash
go build -o intervals-mcp .
```

Получишь бинарь `intervals-mcp` в текущей папке. Запиши его абсолютный путь.

## Подключение к Claude Desktop

> **Важно:** подключи Wahoo/Garmin напрямую к intervals.icu. Strava не отдаёт
> данные о мощности через intervals.

### Где взять API-ключ и Athlete ID

Оба значения лежат в одном месте на intervals.icu:

1. Залогинься и открой **Settings** (прямая ссылка — `intervals.icu/settings`).
2. Прокрути в самый низ — секция **Developer Settings**.
3. Скопируй API-ключ → это `INTERVALS_API_KEY`. Если поле пустое — рядом
   кнопка сгенерировать.
4. Там же виден **Athlete ID**. Начинается с буквы `i` (например, `i12345`) —
   это `INTERVALS_ATHLETE_ID`. Его же видно в адресной строке на своей
   странице: `intervals.icu/athlete/i12345/...`.

### Конфиг Claude Desktop

`Settings → Developer → Edit Config` (откроет `claude_desktop_config.json`),
добавь:

```json
{
  "mcpServers": {
    "intervals-icu": {
      "command": "/АБСОЛЮТНЫЙ/ПУТЬ/К/intervals-mcp",
      "env": {
        "INTERVALS_API_KEY": "твой_ключ",
        "INTERVALS_ATHLETE_ID": "твой_athlete_id"
      }
    }
  }
}
```

Перезапусти Claude Desktop. Сервер появится в меню «+» → Connectors.
Важно: Claude Desktop запускает конфиг с минимальным PATH — указывай
полный путь к бинарю, не короткое имя.

## Заметки

- Сервер пишет в stdout только MCP-сообщения (ошибки идут в stderr) —
  это обязательное условие stdio-транспорта, не ломай его логами в stdout.
- Ограничение Strava: активности, синканутые в intervals из Strava, через
  /activity и /activities могут отдаваться урезанно. Garmin напрямую — ок.
  get_wellness это не затрагивает.

## Локальная проверка (без Claude)

```bash
INTERVALS_API_KEY=fake INTERVALS_ATHLETE_ID=i123 \
printf '%s\n' \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}' \
| ./intervals-mcp
```

Должен прийти JSON с serverInfo intervals-icu.