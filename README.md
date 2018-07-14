# jstn demo website

This repo contains the JSTN demonstration.

## Deploying

- `npm run build:prod`
- `git add dist && git commit -m "Rebuild production site"
- `git subtree split --branch gh-pages --prefix dist/`

And then push the `gh-pages` and `gh-pages-source` branches.
