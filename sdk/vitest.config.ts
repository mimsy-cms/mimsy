import swc from "unplugin-swc";
import { loadEnv } from "vite";
import path from "path";
import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    root: "./",
    globals: true,
    isolate: false,
    passWithNoTests: true,
    include: ["tests/**/*.test.ts"],
    env: loadEnv("test", process.cwd(), ""),
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      reportsDirectory: "coverage/unit",
      include: ["src/**/*.ts"],
    },
    typecheck: {
      enabled: true,
      tsconfig: "./tsconfig.json",
    },
  },
  resolve: {
    alias: {
      $src: path.resolve(__dirname, "./src"),
      $test: path.resolve(__dirname, "./test"),
    },
  },
  plugins: [swc.vite({ module: { type: "es6" } })],
});
