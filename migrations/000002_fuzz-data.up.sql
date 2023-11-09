INSERT INTO authors (id, name, created_at)
VALUES ('1b2849ec-887e-4193-bce0-e0da3b3e4577', 'Anonymous', NOW());

INSERT INTO wisdoms(id, value, author_id)
VALUES ('2868474d-1e23-4b3c-a981-df4a4c669f3b', 'If the wolf is silent, it''s better not to interrupt it.', '1b2849ec-887e-4193-bce0-e0da3b3e4577'),
       ('1817db00-24de-4fda-b793-31e6b732909d', 'The wolf is weaker than the lion and the tiger, but it doesn''t perform in the circus.', '1b2849ec-887e-4193-bce0-e0da3b3e4577');
