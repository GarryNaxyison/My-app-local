import { spawn } from "node:child_process";

const commandForPlatform = (command) => (process.platform === "win32" ? `${command}.cmd` : command);

function run(command, args, env = process.env) {
  return new Promise((resolve, reject) => {
    const isWindows = process.platform === "win32";
    const childCommand = isWindows ? [commandForPlatform(command), ...args].join(" ") : commandForPlatform(command);
    const childArgs = isWindows ? [] : args;
    const child = spawn(childCommand, childArgs, {
      cwd: new URL("..", import.meta.url),
      env,
      shell: isWindows,
      stdio: "inherit",
    });

    child.on("error", reject);
    child.on("exit", (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`${command} ${args.join(" ")} exited with code ${code}`));
      }
    });
  });
}

const cloudflareEnv = {
  ...process.env,
  PUBLIC_SITE_OUT_DIR: "dist",
  PUBLIC_SITE_EMPTY_OUT_DIR: "true",
};

await run("npm", ["run", "seo:generate"], cloudflareEnv);
await run("npm", ["run", "og:generate"], cloudflareEnv);
await run("npx", ["tsc", "-b"], cloudflareEnv);
await run("npx", ["vite", "build"], cloudflareEnv);
