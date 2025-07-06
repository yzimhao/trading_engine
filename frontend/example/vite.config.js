import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 5174,
    proxy: {
        "/api":{
            target:"http://localhost:8080",
            changeOrigin: true
        },
        "/ws":{
            target:"ws://localhost:8080",
            ws: true,
            changeOrigin: true
        }

    }
  },
  plugins: [
    uni(),
  ],
})
