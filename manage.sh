#!/usr/bin/env bash
# ============================================================
#  占卜网站 - 统一管理脚本
#  用法: ./manage.sh [命令]
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
SERVER_DIR="$PROJECT_DIR/server"
CLIENT_DIR="$PROJECT_DIR/client"

# 打印函数
info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()   { echo -e "${RED}[ERROR]${NC} $1"; }
title() { echo -e "\n${CYAN}━━━ $1 ━━━${NC}\n"; }

# ============================================================
#  环境检查
# ============================================================
check_env() {
    title "环境检查"

    # 检查 .env 文件
    if [ ! -f "$PROJECT_DIR/.env" ]; then
        warn ".env 文件不存在，从模板创建..."
        cp "$PROJECT_DIR/.env.example" "$PROJECT_DIR/.env"
        warn "请编辑 .env 文件配置密码和密钥后再启动"
    fi

    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        err "未安装 Docker，请先安装: https://docs.docker.com/get-docker/"
        exit 1
    fi
    ok "Docker $(docker --version | awk '{print $3}' | tr -d ',')"

    # 检查 Docker Compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
        ok "Docker Compose $(docker compose version --short)"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
        ok "Docker Compose $(docker-compose --version | awk '{print $4}' | tr -d ',')"
    else
        err "未安装 Docker Compose"
        exit 1
    fi
}

# ============================================================
#  构建
# ============================================================
do_build() {
    title "构建镜像"
    cd "$PROJECT_DIR"
    check_env
    $COMPOSE_CMD build --no-cache
    ok "镜像构建完成"
}

# ============================================================
#  启动
# ============================================================
do_start() {
    title "启动服务"
    cd "$PROJECT_DIR"
    check_env

    info "构建并启动所有服务..."
    $COMPOSE_CMD up -d --build

    info "等待服务就绪..."
    sleep 5

    # 检查服务状态
    echo ""
    $COMPOSE_CMD ps
    echo ""

    # 检查健康状态
    local retries=12
    local count=0
    while [ $count -lt $retries ]; do
        if curl -sf http://localhost:18080/api/health > /dev/null 2>&1; then
            ok "后端 API 就绪: http://localhost:18080"
            break
        fi
        count=$((count + 1))
        sleep 2
    done

    if [ $count -eq $retries ]; then
        warn "后端 API 启动超时，请检查日志: ./manage.sh logs"
    fi

    echo ""
    ok "========================================="
    ok "  🔮 占卜网站已启动"
    ok "========================================="
    ok ""
    ok "  前端: http://localhost"
    ok "  API:  http://localhost:18080"
    ok ""
    ok "  查看日志: ./manage.sh logs"
    ok "  停止服务: ./manage.sh stop"
    ok "========================================="
}

# ============================================================
#  停止
# ============================================================
do_stop() {
    title "停止服务"
    cd "$PROJECT_DIR"
    check_env
    $COMPOSE_CMD down
    ok "服务已停止"
}

# ============================================================
#  重启
# ============================================================
do_restart() {
    title "重启服务"
    do_stop
    do_start
}

# ============================================================
#  查看日志
# ============================================================
do_logs() {
    cd "$PROJECT_DIR"
    check_env
    if [ -n "$2" ]; then
        $COMPOSE_CMD logs -f "$2"
    else
        $COMPOSE_CMD logs -f
    fi
}

# ============================================================
#  状态
# ============================================================
do_status() {
    title "服务状态"
    cd "$PROJECT_DIR"
    check_env
    $COMPOSE_CMD ps
    echo ""

    # API 健康检查
    if curl -sf http://localhost:18080/api/health > /dev/null 2>&1; then
        ok "后端 API: ✅ 运行中"
    else
        warn "后端 API: ❌ 未响应"
    fi

    # 前端检查
    if curl -sf http://localhost > /dev/null 2>&1; then
        ok "前端 Web: ✅ 运行中"
    else
        warn "前端 Web: ❌ 未响应"
    fi

    # 数据库检查
    if docker exec zhanbu-db pg_isready -U zhanbu > /dev/null 2>&1; then
        ok "数据库:  ✅ 运行中"
    else
        warn "数据库:  ❌ 未响应"
    fi

    # 数据库统计
    echo ""
    info "数据库统计:"
    docker exec zhanbu-db psql -U zhanbu -d zhanbu -t -c "
        SELECT '  用户数: ' || COUNT(*) FROM users
        UNION ALL
        SELECT '  占卜记录: ' || COUNT(*) FROM divination_records
        UNION ALL
        SELECT '  塔罗牌: ' || COUNT(*) FROM tarot_cards
        UNION ALL
        SELECT '  六十四卦: ' || COUNT(*) FROM hexagrams;
    " 2>/dev/null || warn "无法连接数据库"
}

# ============================================================
#  数据库操作
# ============================================================
do_db() {
    cd "$PROJECT_DIR"
    check_env

    case "${2:-}" in
        shell)
            info "进入 PostgreSQL Shell..."
            docker exec -it zhanbu-db psql -U zhanbu -d zhanbu
            ;;
        backup)
            local backup_file="backup_$(date +%Y%m%d_%H%M%S).sql"
            info "备份数据库到 $backup_file ..."
            docker exec zhanbu-db pg_dump -U zhanbu -d zhanbu > "$PROJECT_DIR/$backup_file"
            ok "备份完成: $backup_file"
            ;;
        restore)
            if [ -z "$3" ]; then
                err "用法: ./manage.sh db restore <备份文件>"
                exit 1
            fi
            info "恢复数据库: $3 ..."
            docker exec -i zhanbu-db psql -U zhanbu -d zhanbu < "$3"
            ok "恢复完成"
            ;;
        reset)
            warn "⚠️  这将清除所有数据！"
            read -p "确认？(y/N): " confirm
            if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                info "重置数据库..."
                $COMPOSE_CMD down -v
                $COMPOSE_CMD up -d db
                sleep 5
                $COMPOSE_CMD up -d api
                ok "数据库已重置"
            else
                info "已取消"
            fi
            ;;
        *)
            echo "用法: ./manage.sh db [命令]"
            echo ""
            echo "命令:"
            echo "  shell    进入 PostgreSQL Shell"
            echo "  backup   备份数据库"
            echo "  restore  恢复数据库"
            echo "  reset    重置数据库（清空所有数据）"
            ;;
    esac
}

# ============================================================
#  本地开发模式（不用 Docker）
# ============================================================
do_dev() {
    title "本地开发模式"
    info "启动本地开发服务..."

    # 检查 Go
    if ! command -v go &> /dev/null; then
        err "未安装 Go"
        exit 1
    fi

    # 检查 Node
    if ! command -v node &> /dev/null; then
        err "未安装 Node.js"
        exit 1
    fi

    # 启动后端
    info "启动后端 (Go)..."
    cd "$SERVER_DIR"
    go run main.go &
    BACKEND_PID=$!

    sleep 3

    # 启动前端
    info "启动前端 (Vite)..."
    cd "$CLIENT_DIR"
    npm run dev &
    FRONTEND_PID=$!

    echo ""
    ok "本地开发服务已启动"
    ok "  前端: http://localhost:5173"
    ok "  API:  http://localhost:18080"
    ok ""
    ok "按 Ctrl+C 停止所有服务"

    # 捕获退出信号
    trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT TERM
    wait
}

# ============================================================
#  清理
# ============================================================
do_clean() {
    title "清理"
    cd "$PROJECT_DIR"
    check_env

    warn "清理所有容器、镜像和数据卷..."
    read -p "确认？(y/N): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        $COMPOSE_CMD down -v --rmi all
        ok "清理完成"
    else
        info "已取消"
    fi
}

# ============================================================
#  帮助
# ============================================================
show_help() {
    echo -e "${CYAN}🔮 占卜网站 - 管理脚本${NC}"
    echo ""
    echo "用法: ./manage.sh [命令]"
    echo ""
    echo "命令:"
    echo -e "  ${GREEN}start${NC}       启动所有服务（Docker）"
    echo -e "  ${GREEN}stop${NC}        停止所有服务"
    echo -e "  ${GREEN}restart${NC}     重启所有服务"
    echo -e "  ${GREEN}status${NC}      查看服务状态和数据库统计"
    echo -e "  ${GREEN}logs${NC}        查看日志（可指定服务: logs api/db/web）"
    echo -e "  ${GREEN}build${NC}       重新构建镜像"
    echo -e "  ${GREEN}dev${NC}         本地开发模式（不用 Docker）"
    echo ""
    echo -e "  ${GREEN}db shell${NC}    进入 PostgreSQL Shell"
    echo -e "  ${GREEN}db backup${NC}   备份数据库"
    echo -e "  ${GREEN}db restore${NC}  恢复数据库"
    echo -e "  ${GREEN}db reset${NC}    重置数据库（清空数据）"
    echo ""
    echo -e "  ${GREEN}clean${NC}       清理所有容器、镜像和数据卷"
    echo -e "  ${GREEN}help${NC}        显示此帮助"
    echo ""
    echo "示例:"
    echo "  ./manage.sh start          # 启动"
    echo "  ./manage.sh logs api       # 查看后端日志"
    echo "  ./manage.sh db backup      # 备份数据库"
    echo "  ./manage.sh dev            # 本地开发"
}

# ============================================================
#  主入口
# ============================================================
case "${1:-help}" in
    start)      do_start ;;
    stop)       do_stop ;;
    restart)    do_restart ;;
    status)     do_status ;;
    logs)       do_logs "$@" ;;
    build)      do_build ;;
    dev)        do_dev ;;
    db)         do_db "$@" ;;
    clean)      do_clean ;;
    help|*)     show_help ;;
esac
