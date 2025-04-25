run:
	pnpm build
	go run main.go

lint:
	pnpm dlx eslint src/**/*.js
