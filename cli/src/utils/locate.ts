import { existsSync, readFileSync } from "fs";
import { resolve, dirname, join } from "path";
/**
 * locateProject returns the location of the project the user is currently in.
 * Here's how this works:
 * 1. If, in the current directory, a mimsy.config.json exists, use the path defined inside of the basePath property.
 * 2. If, in the current directory, a mimsy.schema.json exists, use the current path.
 * 3. Go up a directory, repeat 1. and 2.
 * 4. If at any point in 3. we encounter a .git folder (indicating we're at the root of a git repository), we stop the search, and return a not found error.
 */
export function locateProject(): string {
  let currentDir = process.cwd();

  while (true) {
    // 1. Check for mimsy.config.json and use basePath from it
    const configPath = join(currentDir, "mimsy.config.json");
    if (existsSync(configPath)) {
      try {
        const configContent = readFileSync(configPath, "utf8");
        const config = JSON.parse(configContent);

        if (config.basePath && typeof config.basePath === "string") {
          // Resolve the basePath relative to the config file's directory
          return resolve(currentDir, config.basePath);
        }
      } catch (error) {
        // If config file is malformed, continue searching
        console.warn(
          `Warning: Invalid mimsy.config.json at ${configPath}: ${error instanceof Error ? error.message : error}`,
        );
      }
    }

    // 2. Check for mimsy.schema.json and use current path
    const schemaPath = join(currentDir, "mimsy.schema.json");
    if (existsSync(schemaPath)) {
      return currentDir;
    }

    // Before moving up directories, check for .git folder first (stop condition)
    const gitPath = join(currentDir, ".git");
    if (existsSync(gitPath)) {
      throw new Error("No mimsy project found (reached git repository root)");
    }

    // 3. Go up one directory
    const parentDir = dirname(currentDir);

    // If we've reached the filesystem root, stop searching
    if (parentDir === currentDir) {
      throw new Error("No mimsy project found (reached filesystem root)");
    }

    currentDir = parentDir;
  }
}
