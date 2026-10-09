# Webex CLI documentation site

Astro/Starlight source for the published Webex CLI documentation. Area reference pages live in `src/content/docs/`; update them alongside README and the handwritten sections of `../skill/*/SKILL.md` when commands change. `make codegen` from the repository root regenerates skill command tables, but does not regenerate these site pages.

```bash
npm ci
npm run dev
npm run build
```

Build output is `dist/`. The repository Pages workflow publishes the site. The Postman refresh inventory and command migration map are maintained in `../codegen/INVENTORY.md` and `../codegen/NAMING.md`.
