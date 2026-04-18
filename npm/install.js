const { existsSync, copyFileSync, chmodSync } = require('fs');
const { join, dirname } = require('path');

const PLATFORMS = {
  'darwin-arm64': '@nanoloop/cli-darwin-arm64',
  'darwin-x64': '@nanoloop/cli-darwin-x64',
  'linux-x64': '@nanoloop/cli-linux-x64',
  'linux-arm64': '@nanoloop/cli-linux-arm64',
  'win32-x64': '@nanoloop/cli-win32-x64',
};

const platform = `${process.platform}-${process.arch}`;
const pkg = PLATFORMS[platform];

if (!pkg) {
  console.warn(`nanoloop: Unsupported platform ${platform}, using JS wrapper fallback`);
  process.exit(0);
}

let srcPath;
try {
  const ext = process.platform === 'win32' ? '.exe' : '';
  srcPath = require.resolve(`${pkg}/bin/nanoloop${ext}`);
} catch {
  console.warn(`nanoloop: Could not find binary for ${platform}`);
  process.exit(0);
}

const destPath = join(__dirname, 'bin', 'nanoloop' + (process.platform === 'win32' ? '.exe' : ''));

try {
  copyFileSync(srcPath, destPath);
  chmodSync(destPath, 0o755);
} catch (err) {
  console.warn(`nanoloop: Failed to copy binary: ${err.message}`);
}
