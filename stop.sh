#!/bin/bash

# ============================================================================
# Скрипт остановки Semaphore UI
# ============================================================================
# Останавливает Docker контейнеры с возможностью очистки volumes
#
# Использование: ./stop.sh [OPTIONS]
#
# Опции:
#   --clean, -c      Очистить volumes (удалить данные БД)
#   --all, -a        Полная очистка (остановка + volumes + сети)
#   --dry-run        Показать, что будет сделано (без выполнения)
#   --force, -f      Не запрашивать подтверждение
#   --help, -h       Показать эту справку
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Флаги
CLEAN=false
REMOVE_ALL=false
DRY_RUN=false
FORCE=false

# Функция для вывода справки
show_help() {
    head -22 "$0" | tail -19
    exit 0
}

# Функция запроса подтверждения
confirm() {
    if [ "$FORCE" = true ]; then
        return 0
    fi
    
    local prompt="$1"
    echo -e "${YELLOW}⚠️  $prompt${NC}"
    read -p "Продолжить? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Отменено пользователем${NC}"
        exit 0
    fi
}

# Парсинг аргументов
while [[ $# -gt 0 ]]; do
    case $1 in
        --clean|-c)
            CLEAN=true
            shift
            ;;
        --all|-a)
            REMOVE_ALL=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --force|-f)
            FORCE=true
            shift
            ;;
        --help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}❌ Неизвестный параметр: $1${NC}"
            echo "Используйте --help для справки"
            exit 1
            ;;
    esac
done

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker не установлен. Установите Docker.${NC}"
    exit 1
fi

# Определение команды docker-compose
if docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Semaphore UI - Остановка сервисов                  ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Режим dry-run
if [ "$DRY_RUN" = true ]; then
    echo -e "${CYAN}📋 РЕЖИМ ПРОВЕРКИ (ничего не удаляется)${NC}"
    echo ""
fi

# Остановка контейнеров
echo -e "${YELLOW}⏹️  Остановка Docker контейнеров...${NC}"
if [ "$DRY_RUN" = true ]; then
    echo "   [DRY-RUN] $COMPOSE_CMD -f \"$COMPOSE_FILE\" down"
else
    $COMPOSE_CMD -f "$COMPOSE_FILE" down
fi
echo -e "${GREEN}✓ Контейнеры остановлены${NC}"
echo ""

# Очистка volumes
if [ "$CLEAN" = true ] || [ "$REMOVE_ALL" = true ]; then
    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}📋 Будут удалены volumes:${NC}"
        echo "   [DRY-RUN] $COMPOSE_CMD -f \"$COMPOSE_FILE\" down -v"
        docker volume ls --filter name=semaphore 2>/dev/null || true
    else
        echo -e "${YELLOW}🗑️  Очистка volumes (данные БД будут потеряны)...${NC}"
        confirm "Вы уверены, что хотите удалить volumes? Все данные БД будут потеряны!"
        $COMPOSE_CMD -f "$COMPOSE_FILE" down -v
        echo -e "${GREEN}✓ Volumes очищены${NC}"
    fi
    echo ""
fi

# Удаление сетей (только с --all)
if [ "$REMOVE_ALL" = true ]; then
    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}📋 Будут удалены сети:${NC}"
        docker network ls --filter name=semaphore 2>/dev/null || true
    else
        echo -e "${YELLOW}🗑️  Удаление сетей...${NC}"
        $COMPOSE_CMD -f "$COMPOSE_FILE" down --remove-orphans 2>/dev/null || true
        for network in $(docker network ls --filter name=semaphore --format "{{.Name}}" 2>/dev/null); do
            docker network rm "$network" 2>/dev/null || true
        done
        echo -e "${GREEN}✓ Сети удалены${NC}"
    fi
    echo ""
fi

# Итоговый отчет
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║              Остановка завершена!                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ "$DRY_RUN" = false ]; then
    echo -e "${CYAN}📊 Текущее состояние:${NC}"
    echo ""
    
    # Контейнеры
    echo "Контейнеры:"
    docker ps -a --filter name=semaphore --format "  {{.Names}} - {{.Status}}" 2>/dev/null || echo "  Нет контейнеров semaphore"
    echo ""
    
    # Volumes
    echo "Volumes:"
    docker volume ls --filter name=semaphore --format "  {{.Name}}" 2>/dev/null || echo "  Нет volumes semaphore"
    echo ""
fi

echo -e "${GREEN}✅ Остановка завершена!${NC}"
echo ""
echo -e "${YELLOW}💡 Для повторного запуска выполните:${NC}"
echo -e "   ${CYAN}./start.sh${NC}"
echo ""
