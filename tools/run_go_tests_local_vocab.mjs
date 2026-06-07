import { spawnSync } from "node:child_process";

const env = {
  ...process.env,
  VOCABULARY_DIR: process.env.TEST_VOCABULARY_DIR || "data/vocabulary",
};

const result = spawnSync("go", ["test", "./..."], {
  cwd: process.cwd(),
  env,
  stdio: "inherit",
  shell: process.platform === "win32",
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 1);
