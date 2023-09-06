import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    //    port: 3000,
    proxy: {
      // Using the proxy instance
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        // configure: (proxy, options) => {
        //   // proxy will be an instance of 'http-proxy'
        // }
      },
    }
  }
})
