import { Command } from "commander";
import { writeFileSync } from "fs";
import { resolve } from "path";
import { exportSchema, clearRegistry } from "@mimsy/sdk";

export interface ExportSchemaOptions {
  output: string;
  import?: string;
  pretty: boolean;
  clear: boolean;
}

export async function exportSchemaAction(options: ExportSchemaOptions): Promise<void> {
  try {
    // Clear registry if requested
    if (options.clear) {
      clearRegistry();
    }

    // Import collections from file if specified
    if (options.import) {
      const importPath = resolve(options.import);
      console.log(`📥 Importing collections from: ${importPath}`);

      try {
        if (importPath.endsWith(".ts")) {
          // Register esbuild for TypeScript support
          require("esbuild-register");
          // Use require for TypeScript files after registration
          require(importPath);
        } else {
          // Use import for JavaScript files
          await import(importPath);
        }
        console.log("✅ Collections imported successfully");
      } catch (importError) {
        console.error(
          "❌ Failed to import collections:",
          importError instanceof Error ? importError.message : importError,
        );
        process.exit(1);
      }
    }

    const schema = exportSchema();
    const outputPath = resolve(options.output);

    const jsonContent = options.pretty
      ? JSON.stringify(schema, null, 2)
      : JSON.stringify(schema);

    writeFileSync(outputPath, jsonContent, "utf8");

    console.log(`✅ Schema exported successfully to: ${outputPath}`);
    console.log(`📊 Collections exported: ${schema.collections.length}`);
    console.log(`⏰ Generated at: ${schema.generatedAt}`);
  } catch (error) {
    console.error(
      "❌ Failed to export schema:",
      error instanceof Error ? error.message : error,
    );
    process.exit(1);
  }
}

export function exportSchemaCommand(program: Command): Command {
  return program
    .command("export-schema")
    .description("Export collection schemas to a JSON file")
    .option("-o, --output <path>", "Output file path", "schema.json")
    .option(
      "-i, --import <path>",
      "Import collections from a TypeScript file before exporting",
    )
    .option("--pretty", "Pretty print the JSON output", false)
    .option(
      "--clear",
      "Clear the registry before importing (useful for testing)",
      false,
    )
    .action(exportSchemaAction);
}