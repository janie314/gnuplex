import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: "/home/",
  server: {
    port: 50000,
    strictPort: true,
  },
  preview: {
    port: 50000,
    strictPort: true,
  },
  build: {
    outDir: "out",
  },
});
