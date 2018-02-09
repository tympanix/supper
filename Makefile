REPO := github.com/tympanix/supper

LOCAL := $(shell git rev-parse @)
REMOTE := $(shell git rev-parse @{u})

deploy:
ifeq ($(LOCAL), $(REMOTE))
	@echo "No updates since last deploy"
else
	make update
	make build
endif

update:
ifdef TOKEN
	@/usr/bin/env git pull https://$(TOKEN)@$(REPO)
else
	@/usr/bin/env git pull
endif

build:
	/usr/bin/env npm install
	/usr/bin/env npm run build