import { Command } from "commander";
import { exportSchema, clearRegistry } from "@mimsy/sdk";
import { writeFileSync, mkdirSync, existsSync } from "fs";
import { resolve, join } from "path";
import inquirer from "inquirer";

export interface ApplyOptions {
  description?: string;
  clear: boolean;
}

function sanitizeForFilename(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .substring(0, 20);
}

function formatDate(date: Date): string {
  const year = date.getFullYear().toString().slice(-2);
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  return `${year}${month}${day}`;
}

export async function applyAction(options: ApplyOptions): Promise<void> {
  try {
    let description = options.description;

    if (!description) {
      const answers = await inquirer.prompt([
        {
          type: "input",
          name: "description",
          message: "Please describe the changes being applied:",
          validate: (input: string) => {
            if (!input || input.trim().length === 0) {
              return "Description is required";
            }
            return true;
          },
        },
      ]);
      description = answers.description;
    }

    if (options.clear) {
      clearRegistry();
    }

    const collectionsPath = resolve(process.cwd(), "src/lib/collections.ts");

    if (!existsSync(collectionsPath)) {
      console.error(
        "❌ Error: No collections file found at src/lib/collections.ts",
      );
      console.error(
        "   Please ensure you have a collections file before applying schema changes.",
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

    const schemasDir = resolve(process.cwd(), ".mimsy", "schemas");
    mkdirSync(schemasDir, { recursive: true });

    const date = new Date();
    const dateStr = formatDate(date);
    const shortDesc = sanitizeForFilename(description!);
    const filename = `${dateStr}-${shortDesc}.jsonc`;
    const outputPath = join(schemasDir, filename);

    const header = [
      `// Description: ${description}`,
      `// Applied at: ${date.toISOString()}`,
      `// Version: @mimsy/cli@1.0.0`,
      "",
    ].join("\n");

    const jsonContent = JSON.stringify(schema, null, 2);
    const fileContent = header + jsonContent;

    writeFileSync(outputPath, fileContent, "utf8");

    console.log(`✅ Schema applied successfully`);
    console.log(`📁 Saved to: ${outputPath}`);
    console.log(`📊 Collections exported: ${schema.collections.length}`);
    console.log(`📝 Description: ${description}`);
  } catch (error) {
    console.error(
      "❌ Failed to apply schema:",
      error instanceof Error ? error.message : error,
    );
    process.exit(1);
  }
}

export function applyCommand(program: Command): Command {
  return program
    .command("apply")
    .description("Apply and save schema changes with a description")
    .option(
      "-d, --description <text>",
      "Description of the changes (will prompt if not provided)",
    )
    .option("--clear", "Clear the registry before importing collections", true)
    .action(applyAction);
}
