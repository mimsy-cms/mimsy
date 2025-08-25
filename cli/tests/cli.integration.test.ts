import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { createProgram } from "../src/index";
import { version } from "../src/version";
import { Command } from "commander";

vi.mock("$src/utils/locate");

describe("CLI integration tests", () => {
  let consoleSpy: ReturnType<typeof vi.spyOn>;
  let processExitSpy: any;
  let program: Command;

  beforeEach(() => {
    consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    processExitSpy = vi.spyOn(process, "exit").mockImplementation(() => {
      throw new Error("process.exit called");
    });
    program = createProgram();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("should have correct program metadata", () => {
    expect(program.name()).toBe("msy");
    expect(program.description()).toBe(
      "A CLI tool for mimsy, the simple SvelteKit CMS"
    );
    expect(program.version()).toBe(version);
  });

  it("should have all commands registered", () => {
    const commandNames = program.commands.map((cmd) => cmd.name());
    expect(commandNames).toContain("update");
    expect(commandNames).toContain("init");
    expect(commandNames).toHaveLength(2);
  });

  it("should show error when collections file is missing", async () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    const { locateProject } = await import("$src/utils/locate");
    const mockLocateProject = vi.mocked(locateProject);
    
    mockLocateProject.mockReturnValue("/test/project");
    program.exitOverride();

    try {
      await program.parseAsync(["update"], { from: "user" });
    } catch (e) {
      // Expected to throw due to process.exit
    }

    expect(errorSpy).toHaveBeenCalledWith(
      expect.stringContaining("No collections file found")
    );
  });

  it("should show help when --help flag is used", () => {
    const writeSpy = vi
      .spyOn(process.stdout, "write")
      .mockImplementation(() => true);
    const helpProgram = createProgram();
    helpProgram.exitOverride();

    try {
      helpProgram.parse(["node", "msy", "--help"], { from: "user" });
    } catch (error) {
      // Expected to throw due to exitOverride
    }

    const output = writeSpy.mock.calls.map((call) => call[0]).join("");
    expect(output).toContain("A CLI tool for mimsy");
    expect(output).toContain("update");
    expect(output).toContain("Update the mimsy.schema.json");
  });

  it("should show version with --version flag", () => {
    const writeSpy = vi
      .spyOn(process.stdout, "write")
      .mockImplementation(() => true);
    const versionProgram = createProgram();
    versionProgram.exitOverride();

    try {
      versionProgram.parse(["node", "msy", "--version"], { from: "user" });
    } catch (error) {
      // Expected to throw due to exitOverride
    }

    expect(writeSpy).toHaveBeenCalledWith(expect.stringContaining(version));
  });

  it("should handle unknown commands gracefully", () => {
    const errorSpy = vi
      .spyOn(process.stderr, "write")
      .mockImplementation(() => true);
    const unknownProgram = createProgram();
    unknownProgram.exitOverride();

    expect(() => {
      unknownProgram.parse(["node", "msy", "unknown-command"], {
        from: "user",
      });
    }).toThrow();

    const errorOutput = errorSpy.mock.calls.map((call) => call[0]).join("");
    expect(errorOutput).toContain("unknown command");
  });
});
