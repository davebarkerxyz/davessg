# TODO

- Extract frontmatter from content files
- Load all files in templates/ (excl static) to map[string][]byte
  - Allow 'template: blah.html' in frontmatter
- Partials folder (instead of k/v site.conf)
  - Regex to find {{ partial mypartial.html }} and replace with content from files in partials folder
