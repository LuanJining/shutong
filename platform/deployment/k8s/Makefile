deploy:
	./fast-deploy.sh

uninstall:
	./fast-uninstall.sh

check:
	./fast-check.sh

init-db:
	./fast-init-db.sh

upgrade:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "请提供版本号，例如: make upgrade v1.0.2"; \
		exit 1; \
	fi
	@./fast-upgrade-image.sh $(filter-out $@,$(MAKECMDGOALS))

# 防止 make 将版本号当作目标
%:
	@: