# ujiprog.com

- Uses [`syumai/workers`](https://github.com/syumai/workers) package to run an HTTP server.

## Requirements

- Node.js
- Go 1.24.0 or later

### Commands

```
npm start      # run dev server
# or
go run .       # run dev server without Wrangler (Cloudflare-related features are not available)
npm run build  # build Go Wasm binary
npm run deploy # deploy worker
```
