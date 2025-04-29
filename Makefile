.PHONY: RunTestServers

RunTestServers:
	@echo "Starting test servers..."
	@echo "Server 1 (port 9001) - http://localhost:9001"
	@echo "Server 2 (port 9002) - http://localhost:9002"
	@echo "Health checks:"
	@echo "  curl http://localhost:9001/health"
	@echo "  curl http://localhost:9002/health"
	
	# Запуск серверов в фоновом режиме с перенаправлением логов
	@go run server1.go > server1.log 2>&1 & \
	echo $$! > server1.pid; \
	go run server2.go > server2.log 2>&1 & \
	echo $$! > server2.pid; \
	echo "Servers started with PIDs: $$(cat server1.pid) and $$(cat server2.pid)"

stop-servers:
	@echo "Stopping test servers..."
	@-kill $$(cat server1.pid) 2>/dev/null || true
	@-kill $$(cat server2.pid) 2>/dev/null || true
	@-rm *.pid *.log 2>/dev/null || true
	@echo "Servers stopped"

clean:
	@rm -f *.log *.pid

check-servers:
	@echo "Checking servers status..."
	@curl -s http://localhost:9001/health || echo "Server 1 is DOWN"
	@curl -s http://localhost:9002/health || echo "Server 2 is DOWN"