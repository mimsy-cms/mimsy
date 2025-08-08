# Mimsy CLI

A command-line interface for the Mimsy CMS SDK that provides utilities for working with collection schemas and more.

## Installation

```bash
# Install globally
pnpm install -g @mimsy/cli

# Or run directly with pnpm
pnpm dlx @mimsy/cli
```

## Commands

### `export-schema`

Export collection schemas to a JSON file. This command allows you to serialize your collection definitions into a portable JSON format.

#### Usage

```bash
msy export-schema [options]
```

#### Options

- `-o, --output <path>` - Output file path (default: "schema.json")
- `-i, --import <path>` - Import collections from a TypeScript or JavaScript file before exporting
- `--pretty` - Pretty print the JSON output (default: false)
- `--clear` - Clear the registry before importing (useful for testing) (default: false)
- `-h, --help` - Display help for command

#### Examples

**Basic usage:**
```bash
# Export current registry to schema.json
msy export-schema

# Export with pretty formatting
msy export-schema --pretty
```

**Import collections and export:**
```bash
# Import TypeScript collections and export
msy export-schema --import ./collections/blog.ts --pretty --output blog-schema.json

# Import JavaScript collections and export
msy export-schema --import ./collections/blog.js --output blog-schema.json

# Clear registry first, then import and export
msy export-schema --import ./collections/blog.ts --clear --pretty
```

#### Collection File Examples

**TypeScript Example (`collections/blog.ts`):**
```typescript
import { collection, fields, builtins, type Collection } from "@mimsy/sdk";

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
```

**JavaScript Example (`collections/blog.js`):**
```javascript
const { collection, fields, builtins } = require("@mimsy/sdk");

const Tags = collection("tags", {
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

const Posts = collection("posts", {
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

module.exports = { Tags, Posts };
```

#### Output Schema Format

The exported JSON schema has the following structure:

```json
{
  "collections": [
    {
      "name": "collection-name",
      "schema": {
        "fieldName": {
          "type": "field-type",
          "relatesTo": "related-collection-name",
          "options": {
            "description": "Field description",
            "constraints": {
              "required": true,
              "minLength": 5
            }
          }
        }
      }
    }
  ],
  "generatedAt": "2025-07-22T11:52:56.435Z"
}
```

## Development

### Building

```bash
pnpm build
```

### Testing

```bash
pnpm test
```

### Development Mode

```bash
pnpm dev
```

## TypeScript Support

The CLI supports importing TypeScript files directly thanks to [esbuild-register](https://www.npmjs.com/package/esbuild-register). No compilation step is required - just point to your `.ts` files and they'll be transpiled on the fly.

## Workspace Integration

This CLI is designed to work within the Mimsy monorepo workspace. Collections are automatically registered when created using the `collection()` function from the SDK, making schema export seamless.

## Contributing

Please see the main [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## License

See [LICENSE](../LICENCE) for more information.