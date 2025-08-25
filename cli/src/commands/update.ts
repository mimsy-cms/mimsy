import { Command } from "commander";
import { exportSchema, clearRegistry } from "@mimsy-cms/sdk";
import { writeFileSync, existsSync } from "fs";
import { resolve } from "path";
import { version } from "../version";
import { locateProject } from "$src/utils/locate";

export interface UpdateOptions {
  clear: boolean;
}

export async function updateAction(options: UpdateOptions): Promise<void> {
  try {
    if (options.clear) {
      clearRegistry();
    }

    let projectPath;
    try {
      projectPath = locateProject();
    } catch (error) {
      console.error("❌ Could not find a mimsy project: ", error);
      console.info(
        "   If you want to create a new project, use the `init` command.",
      );
      process.exit(1);
    }

    const collectionsPath = resolve(projectPath, "src/lib/collections.ts");

    if (!existsSync(collectionsPath)) {
      console.error(
        "❌ Error: No collections file found at src/lib/collections.ts",
      );
      console.error(
        "   Please ensure you have a collections file before updating schema.",
      );
      process.exit(1);
    }

    console.log(`📥 Importing collections from: ${collectionsPath}`);

    try {
      require("esbuild-register");
      require(collectionsPath);
      console.log("✅ Collections imported successfully");
    } catch (importError) {
      console.error(
        "❌ Failed to import collections:",
        importError instanceof Error ? importError.message : importError,
      );
      process.exit(1);
    }

    const schema = exportSchema();
    const outputPath = resolve(projectPath, "mimsy.schema.json");

    const jsonContent = JSON.stringify(schema, null, 2);
    const fileContent = jsonContent;

    writeFileSync(outputPath, fileContent, "utf8");

    console.log(`✅ Schema updated successfully`);
    console.log(`📁 Updated: ${outputPath}`);
    console.log(`📊 Collections exported: ${schema.collections.length}`);
  } catch (error) {
    console.error(
      "❌ Failed to update schema:",
      error instanceof Error ? error.message : error,
    );
    process.exit(1);
  }
}

export function updateCommand(program: Command): Command {
  return program
    .command("update")
    .description("Update the mimsy.schema.json file in the root of the package")
    .option("--clear", "Clear the registry before importing collections", true)
    .action(updateAction);
}
