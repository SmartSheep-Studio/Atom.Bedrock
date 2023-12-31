import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";
import VueI18nPlugin from "@intlify/unplugin-vue-i18n/vite";
import { VitePWA } from 'vite-plugin-pwa'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      imports: [
        "vue",
        {
          "naive-ui": ["useDialog", "useMessage", "useNotification", "useLoadingBar"]
        }
      ]
    }),
    Components({
      resolvers: [NaiveUiResolver()]
    }),
    VueI18nPlugin({ runtimeOnly: false }),
    VitePWA({
      manifest: {
        name: "Atom",
        description: "Committed to improving the Internet experience",
        theme_color: "#ca4d4d",
        icons: [
          {
            src: "/favicon.png",
            sizes: "256x256",
            type: "image/png",
          },
        ]
      },
      registerType: "autoUpdate",
      devOptions: {
        enabled: true
      }
    })
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url))
    }
  },
  build: {
    sourcemap: true
  },
  server: {
    proxy: {
      "/api": {
        target: "http://127.0.0.1:9443",
        changeOrigin: true
      },
      "/cgi": {
        target: "http://127.0.0.1:9443",
        changeOrigin: true
      },
      "/srv": {
        target: "http://127.0.0.1:9443",
        changeOrigin: true
      }
    }
  }
});
