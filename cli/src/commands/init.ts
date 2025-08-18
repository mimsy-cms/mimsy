import { Command } from "commander";
import { writeFileSync, existsSync, mkdirSync } from "fs";
import { resolve, dirname } from "path";
import inquirer from "inquirer";

export interface InitAnswers {
  schemaJsonPath: string;
  schemaCollectionsPath: string;
  confirmPaths: boolean;
}

const DEFAULT_SCHEMA_JSON = "mimsy.schema.json";
const DEFAULT_COLLECTIONS_TS = "src/lib/collections.ts";

const INITIAL_COLLECTIONS_CONTENT = `
{
  "collections": [
    {
      "name": "tags",
      "schema": {
        "name": {
          "type": "string",
          "options": {
            "description": "The name of the tag",
            "constraints": { "minLength": 2, "maxLength": 50 }
          }
        },
        "color": {
          "type": "string",
          "options": {
            "description": "The color of the tag, in **hexadecimal** format",
            "constraints": { "minLength": 6, "maxLength": 6 }
          }
        }
      }
    },
    {
      "name": "posts",
      "schema": {
        "title": {
          "type": "string",
          "options": {
            "description": "The title of the post",
            "constraints": { "minLength": 5, "maxLength": 100 }
          }
        },
        "author": {
          "type": "relation",
          "relatesTo": "User",
          "options": {
            "description": "The author of the post",
            "constraints": { "required": true }
          }
        },
        "tags": {
          "type": "relation",
          "relatesTo": "tags",
          "options": {
            "description": "The tags associated with the post",
            "constraints": { "required": true }
          }
        },
        "coverImage": {
          "type": "relation",
          "relatesTo": "Media",
          "options": {
            "description": "The cover image of the post",
            "constraints": { "required": true }
          }
        }
      }
    }
  ],
  "generatedAt": "2025-08-11T07:55:00.170Z"
}
`;

const INITIAL_SCHEMA_TS_CONTENT = `
import { collection, fields, builtins, type Collection } from "@mimsy-cms/sdk";

export const Tags: Collection<any> = collection("tags", {
  name: fields.shortString({
    description: "The name of the tag",
    constraints: {
      minLength: 2,
      maxLength: 50,
    },
  }),
  color: fields.shortString({
    description: "The color of the tag, in hexadecimal format",
    constraints: {
      minLength: 6,
      maxLength: 6,
    },
  }),
});

export const Posts: Collection<any> = collection("posts", {
  title: fields.shortString({
    description: "The title of the post",
    constraints: {
      minLength: 5,
      maxLength: 100,
    },
  }),
  author: fields.relation({
    description: "The author of the post",
    relatesTo: builtins.User,
    constraints: {
      required: true,
    },
  }),
  tags: fields.multiRelation({
    description: "The tags associated with the post",
    relatesTo: Tags,
    constraints: {
      required: true,
    },
  }),
  coverImage: fields.media({
    description: "The cover image of the post",
    constraints: {
      required: true,
    },
  }),
});
`;

export async function initAction(): Promise<void> {
  try {
    const basePath = process.cwd();

    console.log("🚀 Initializing mimsy project...");
    console.log(`📁 Base directory: ${basePath}`);

    const configPath = resolve(basePath, "mimsy.config.json");
    const defaultJsonPath = resolve(basePath, DEFAULT_SCHEMA_JSON);
    const defaultTsPath = resolve(basePath, DEFAULT_COLLECTIONS_TS);

    const answers = await inquirer.prompt<InitAnswers>([
      {
        type: "input",
        name: "schemaJsonPath",
        message: "Where should the mimsy.schema.json file be placed?",
        default: defaultJsonPath,
        validate: (input: string) => {
          if (!input.trim()) {
            return "Path cannot be empty";
          }
          if (!input.endsWith(".json")) {
            return "File must have a .json extension";
          }
          return true;
        },
      },
      {
        type: "input",
        name: "schemaCollectionsPath",
        message: "Where should the collections.ts file be placed?",
        default: defaultTsPath,
        validate: (input: string) => {
          if (!input.trim()) {
            return "Path cannot be empty";
          }
          if (!input.endsWith(".ts")) {
            return "File must have a .ts extension";
          }
          return true;
        },
      },
      {
        type: "confirm",
        name: "confirmPaths",
        message: (answers) =>
          `📝 Confirm creation of:\n` +
          `   JSON: ${configPath}\n` +
          `   JSON: ${answers.schemaJsonPath}\n` +
          `   TS:   ${answers.schemaCollectionsPath}\n` +
          `   Continue?`,
        default: true,
      },
    ]);

    if (!answers.confirmPaths) {
      console.log("❌ Initialization cancelled.");
      return;
    }

    const jsonPath = resolve(answers.schemaJsonPath);
    const collectionsPath = resolve(answers.schemaCollectionsPath);

    // Create directories if they don't exist
    const jsonDir = dirname(jsonPath);
    const tsDir = dirname(collectionsPath);

    if (!existsSync(jsonDir)) {
      mkdirSync(jsonDir, { recursive: true });
      console.log(`📁 Created directory: ${jsonDir}`);
    }

    if (!existsSync(tsDir)) {
      mkdirSync(tsDir, { recursive: true });
      console.log(`📁 Created directory: ${tsDir}`);
    }

    const configContent = JSON.stringify(
      {
        basePath,
      },
      null,
      2
    );
    writeFileSync(configPath, configContent, "utf8");
    console.log(`✅ Created: ${configPath}`);

    // TODO: We could improve this by using the applyAction command,
    // but it would require installing the @mimsy/sdk package which is a bit tricky.
    // For now, we just add the content manually.
    writeFileSync(jsonPath, INITIAL_COLLECTIONS_CONTENT, "utf8");
    console.log(`✅ Created: ${jsonPath}`);

    writeFileSync(collectionsPath, INITIAL_SCHEMA_TS_CONTENT, "utf8");
    console.log(`✅ Created: ${collectionsPath}`);

    console.log();
    console.log("🎉 Mimsy project initialized successfully!");
    console.log();
    console.log("📋 Next steps:");
    console.log("1. Define your collections in the typescript schema file");
    console.log("2. Use 'msy update' to update the mimsy json schema");
  } catch (error) {
    console.error(
      "❌ Failed to initialize project:",
      error instanceof Error ? error.message : error
    );
    process.exit(1);
  }
}

export function initCommand(program: Command): Command {
  return program
    .command("init")
    .description("Initialize a new mimsy project with schema files")
    .action(initAction);
}
