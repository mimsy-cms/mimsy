import { collection, fields, builtins } from "@mimsy-cms/sdk";

const Tags = collection("tags", {
  name: fields.shortString({
    description: "The name of the tag",
    constraints: {
      minLength: 2,
      maxLength: 50,
    },
  }),
  color: fields.shortString({
    description: "The color of the tag, in **hexadecimal** format",
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

console.log("Blog collections loaded and registered.");

export default { Tags, Posts };
