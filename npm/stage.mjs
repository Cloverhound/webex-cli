#!/usr/bin/env node
// Builds the npm packages from a goreleaser dist directory:
//   one @cloverhound/webex-cli-<os>-<cpu> package per binary, and the
//   @cloverhound/webex-cli wrapper that depends on all of them.
//
// Usage: node npm/stage.mjs [--dist dist] [--out npm/build] [--version X.Y.Z]

import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { parseArgs } from "node:util";

const here = path.dirname(fileURLToPath(import.meta.url));
const { values: args } = parseArgs({
  options: {
    dist: { type: "string", default: "dist" },
    out: { type: "string", default: path.join(here, "build") },
    version: { type: "string" },
  },
});

const OS = { darwin: "darwin", linux: "linux", windows: "win32" };
const CPU = { amd64: "x64", arm64: "arm64" };
const EXPECTED = Object.keys(OS).length * Object.keys(CPU).length;

const readJSON = (p) => JSON.parse(fs.readFileSync(p, "utf8"));
const writeJSON = (p, v) => fs.writeFileSync(p, JSON.stringify(v, null, 2) + "\n");

const version = (args.version ?? readJSON(path.join(args.dist, "metadata.json")).version).replace(/^v/, "");
const template = readJSON(path.join(here, "webex-cli", "package.json"));

const binaries = readJSON(path.join(args.dist, "artifacts.json")).filter(
  (a) => a.type === "Binary" && OS[a.goos] && CPU[a.goarch],
);
if (binaries.length !== EXPECTED) {
  throw new Error(`expected ${EXPECTED} binaries in ${args.dist}/artifacts.json, found ${binaries.length}`);
}

fs.rmSync(args.out, { recursive: true, force: true });

const platformPackages = {};
for (const a of binaries) {
  const os = OS[a.goos];
  const cpu = CPU[a.goarch];
  const name = `${template.name}-${os}-${cpu}`;
  const dir = path.join(args.out, `webex-cli-${os}-${cpu}`);
  const exe = a.goos === "windows" ? "webex.exe" : "webex";

  fs.mkdirSync(path.join(dir, "bin"), { recursive: true });
  fs.copyFileSync(a.path, path.join(dir, "bin", exe));
  fs.chmodSync(path.join(dir, "bin", exe), 0o755);
  writeJSON(path.join(dir, "package.json"), {
    name,
    version,
    description: `The ${os}-${cpu} binary for ${template.name}`,
    license: template.license,
    homepage: template.homepage,
    repository: template.repository,
    os: [os],
    cpu: [cpu],
    files: ["bin"],
    preferUnplugged: true,
  });
  fs.writeFileSync(
    path.join(dir, "README.md"),
    `# ${name}\n\nThe ${os}-${cpu} binary for [${template.name}](https://www.npmjs.com/package/${template.name}). Install that package instead.\n`,
  );
  platformPackages[name] = version;
}

const mainDir = path.join(args.out, "webex-cli");
fs.cpSync(path.join(here, "webex-cli"), mainDir, { recursive: true });
fs.chmodSync(path.join(mainDir, "bin", "webex.js"), 0o755);
writeJSON(path.join(mainDir, "package.json"), {
  ...template,
  version,
  optionalDependencies: Object.fromEntries(Object.entries(platformPackages).sort()),
});

console.log(`Staged ${template.name}@${version} and ${binaries.length} platform packages in ${args.out}`);
