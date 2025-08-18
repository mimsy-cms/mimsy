import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { updateAction, updateCommand } from "../../src/commands/update";
import { Command } from "commander";
import fs from "fs";
import path from "path";
import * as sdk from "@mimsy-cms/sdk";
import Module from "module";
import { version } from "$src/version";

vi.mock("fs");
vi.mock("@mimsy-cms/sdk");

describe("update command", () => {
  let consoleSpy: ReturnType<typeof vi.spyOn>;
  let consoleErrorSpy: ReturnType<typeof vi.spyOn>;
  let processExitSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    processExitSpy = vi.spyOn(process, "exit").mockImplementation((code) => {
      throw new Error(`process.exit(${code})`);
    });
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("updateCommand", () => {
    it("should create update command with correct configuration", () => {
      const program = new Command();
      const command = updateCommand(program);

      expect(command.name()).toBe("update");
      expect(command.description()).toBe(
        "Update the mimsy.schema.json file in the root of the package",
      );

      const options = command.options;
      const clearOption = options.find((opt) => opt.long === "--clear");
      expect(clearOption).toBeDefined();
      expect(clearOption?.description).toBe(
        "Clear the registry before importing collections",
      );
    });
  });

  describe("updateAction", () => {
    it("should clear registry when clear option is true", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockWriteFileSync = vi.mocked(fs.writeFileSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);
      const mockClearRegistry = vi.mocked(sdk.clearRegistry);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockReturnValue({
        collections: [{ name: "test" }],
      } as any);

      await updateAction({ clear: true });

      expect(mockClearRegistry).toHaveBeenCalled();
      expect(mockWriteFileSync).toHaveBeenCalled();
      expect(consoleSpy).toHaveBeenCalledWith(
        expect.stringContaining("✅ Schema updated successfully"),
      );

      Module.prototype.require = originalRequire;
    });

    it("should not clear registry when clear option is false", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockWriteFileSync = vi.mocked(fs.writeFileSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);
      const mockClearRegistry = vi.mocked(sdk.clearRegistry);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockReturnValue({
        collections: [{ name: "test" }],
      } as any);

      await updateAction({ clear: false });

      expect(mockClearRegistry).not.toHaveBeenCalled();
      expect(mockWriteFileSync).toHaveBeenCalled();

      Module.prototype.require = originalRequire;
    });

    it("should exit with error when collections file does not exist", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);

      mockExistsSync.mockReturnValue(false);

      await expect(updateAction({ clear: false })).rejects.toThrow(
        "process.exit(1)",
      );

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        "❌ Error: No collections file found at src/lib/collections.ts",
      );
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        "   Please ensure you have a collections file before updating schema.",
      );
    });

    it("should write schema with correct format", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockWriteFileSync = vi.mocked(fs.writeFileSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);

      const testSchema = {
        collections: [
          { name: "posts", fields: [] },
          { name: "pages", fields: [] },
        ],
      };

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockReturnValue(testSchema as any);

      await updateAction({ clear: false });

      expect(mockWriteFileSync).toHaveBeenCalledWith(
        expect.stringContaining("mimsy.schema.json"),
        expect.stringContaining("// Updated at:"),
        "utf8",
      );

      expect(mockWriteFileSync).toHaveBeenCalledWith(
        expect.anything(),
        expect.stringContaining(`// Version: @mimsy-cms/cli@${version}`),
        "utf8",
      );

      expect(consoleSpy).toHaveBeenCalledWith("📊 Collections exported: 2");

      Module.prototype.require = originalRequire;
    });

    it("should handle import errors gracefully", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register") {
          return {};
        }
        if (id.includes("collections.ts")) {
          throw new Error("Import failed: Invalid syntax");
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);

      await expect(updateAction({ clear: false })).rejects.toThrow(
        "process.exit(1)",
      );

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        "❌ Failed to import collections:",
        "Import failed: Invalid syntax",
      );

      Module.prototype.require = originalRequire;
    });

    it("should handle schema export errors", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockImplementation(() => {
        throw new Error("Schema export failed");
      });

      await expect(updateAction({ clear: false })).rejects.toThrow(
        "process.exit(1)",
      );

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        "❌ Failed to update schema:",
        "Schema export failed",
      );

      Module.prototype.require = originalRequire;
    });

    it("should handle file write errors", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockWriteFileSync = vi.mocked(fs.writeFileSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockReturnValue({
        collections: [],
      } as any);
      mockWriteFileSync.mockImplementation(() => {
        throw new Error("Permission denied");
      });

      await expect(updateAction({ clear: false })).rejects.toThrow(
        "process.exit(1)",
      );

      expect(consoleErrorSpy).toHaveBeenCalledWith(
        "❌ Failed to update schema:",
        "Permission denied",
      );

      Module.prototype.require = originalRequire;
    });

    it("should log correct success messages", async () => {
      const mockExistsSync = vi.mocked(fs.existsSync);
      const mockWriteFileSync = vi.mocked(fs.writeFileSync);
      const mockExportSchema = vi.mocked(sdk.exportSchema);

      const originalRequire = Module.prototype.require;
      Module.prototype.require = vi.fn((id) => {
        if (id === "esbuild-register" || id.includes("collections.ts")) {
          return {};
        }
        return originalRequire.call(this, id);
      });

      mockExistsSync.mockReturnValue(true);
      mockExportSchema.mockReturnValue({
        collections: [{ name: "test" }],
      } as any);

      await updateAction({ clear: false });

      expect(consoleSpy).toHaveBeenCalledWith(
        expect.stringContaining("📥 Importing collections from:"),
      );
      expect(consoleSpy).toHaveBeenCalledWith(
        "✅ Collections imported successfully",
      );
      expect(consoleSpy).toHaveBeenCalledWith("✅ Schema updated successfully");
      expect(consoleSpy).toHaveBeenCalledWith(
        expect.stringContaining("📁 Updated:"),
      );
      expect(consoleSpy).toHaveBeenCalledWith("📊 Collections exported: 1");

      Module.prototype.require = originalRequire;
    });
  });
});
