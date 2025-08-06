INSERT INTO collections (slug, name, fields)
VALUES ("posts", "Post", "TODO"), ("tags", "Tag", "TODO");


CREATE table posts (
  id BIGINT PRIMARY KEY,
  slug VARCHAR(256),
  cover_image BIGINT,
  contents JSONB,
  author BIGINT,
  FOREIGN KEY (cover_image) REFERENCES media(id),
  FOREIGN KEY (author) REFERENCES users(id)
);

CREATE TABLE tags (
  id BIGINT PRIMARY KEY,
  slug VARCHAR(256),
  color VARCHAR(256),
  name VARCHAR(256)
);

CREATE TABLE posts_tags_relation (
    post_id BIGINT,
    tags_slug VARCHAR(256),
    -- Foreign keys
    PRIMARY KEY (post_id, tags_slug),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tags_slug) REFERENCES tags(slug)
);
