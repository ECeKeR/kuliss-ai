import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  build: {
    outDir: "../frontend/dist",
    emptyOutDir: true,
  },
  server: {
    port: 9245,
    strictPort: true,
    host: "0.0.0.0"
  }
});
