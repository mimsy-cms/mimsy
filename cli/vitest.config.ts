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
      provider: "istanbul",
      reporter: ["text", "json", "html"],
      reportsDirectory: "coverage/unit",
      include: ["src/**/*.ts"],
    },
    setupFiles: ['./tests/setup.ts'],
  },
  resolve: {
    alias: {
      $src: path.resolve(__dirname, "./src"),
      $test: path.resolve(__dirname, "./test"),
      'esbuild-register': path.resolve(__dirname, './tests/__mocks__/esbuild-register.ts'),
    },
  },
  plugins: [swc.vite({ module: { type: "es6" } })],
});
