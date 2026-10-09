#!/usr/bin/env node
"use strict";

// Runs the webex binary from the platform package npm installed alongside this
// one (an optionalDependency selected by its os/cpu fields).

const { spawn } = require("node:child_process");
const fs = require("node:fs");

const SUPPORTED = new Set([
  "darwin-arm64",
  "darwin-x64",
  "linux-arm64",
  "linux-x64",
  "win32-arm64",
  "win32-x64",
]);

function fail(message) {
  process.stderr.write(`webex-cli: ${message}\n`);
  process.exit(1);
}

function binaryPath() {
  const key = `${process.platform}-${process.arch}`;
  if (!SUPPORTED.has(key)) {
    fail(`no build for ${key}. See https://github.com/Cloverhound/webex-cli#install for other install methods.`);
  }
  const pkg = `@cloverhound/webex-cli-${key}`;
  const exe = process.platform === "win32" ? "webex.exe" : "webex";
  let bin;
  try {
    bin = require.resolve(`${pkg}/bin/${exe}`);
  } catch {
    fail(`${pkg} is not installed. Reinstall without --omit=optional or --no-optional.`);
  }
  if (process.platform !== "win32") {
    try {
      fs.accessSync(bin, fs.constants.X_OK);
    } catch {
      // Some package managers drop the executable bit when unpacking.
      try {
        fs.chmodSync(bin, 0o755);
      } catch {}
    }
  }
  return bin;
}

const child = spawn(binaryPath(), process.argv.slice(2), { stdio: "inherit" });

// MCP hosts stop servers with a signal; pass it on so the binary is not orphaned.
for (const signal of ["SIGINT", "SIGTERM", "SIGHUP"]) {
  process.on(signal, () => {
    if (!child.killed) child.kill(signal);
  });
}

child.on("error", (err) => fail(err.message));
child.on("exit", (code, signal) => {
  if (signal) {
    process.removeAllListeners(signal);
    process.kill(process.pid, signal);
  } else {
    process.exit(code ?? 1);
  }
});
