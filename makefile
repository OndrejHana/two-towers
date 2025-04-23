run:
	pnpm dlx eslint src/**/*.js
	pnpm build
	go run main.go
