const fs = require('fs');
const path = require('path');
const https = require('https');

const pkg = require('../package.json');
const version = pkg.version;

const platformMap = {
  'darwin-arm64': 'darwin-arm64',
  'darwin-x64': 'darwin-amd64',
  'linux-x64': 'linux-amd64',
  'linux-arm64': 'linux-arm64',
  'win32-x64': 'windows-amd64',
};

const platform = `${process.platform}-${process.arch}`;
const binaryName = platformMap[platform];

if (!binaryName) {
  console.error(`Unsupported platform: ${platform}`);
  process.exit(1);
}

const isWindows = process.platform === 'win32';
const ext = isWindows ? '.exe' : '';
const filename = `nanoloop-${binaryName}${ext}`;

const binDir = path.join(__dirname, '..', 'bin');
const binPath = path.join(binDir, isWindows ? 'nanoloop.exe' : 'nanoloop');

fs.mkdirSync(binDir, { recursive: true });

const baseUrl = `https://github.com/nanoloophq/cli/releases/download/v${version}`;
const url = `${baseUrl}/${filename}`;

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        download(res.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      if (res.statusCode !== 200) {
        reject(new Error(`Download failed: ${res.statusCode} ${url}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', reject);
  });
}

async function main() {
  console.log(`Downloading nanoloop ${version} for ${platform}...`);

  try {
    await download(url, binPath);
    fs.chmodSync(binPath, 0o755);
    console.log('Done.');
  } catch (err) {
    console.error(`Failed to download binary: ${err.message}`);
    console.error(`URL: ${url}`);
    process.exit(1);
  }
}

main();
