
rm -rf ./release

go build || exit 1

echo "Build successful"

mkdir ./release
mkdir ./release/FactoCord3

version=`git describe --tags`

if echo "${version}" | grep -q "\-g"; then
  echo "Error: no tag for this version"
  exit 2
fi

version+=" ("
version+=`git rev-parse --short HEAD`
version+=")"

echo "Version: ${version}"
echo "${version}" > ./release/FactoCord3/.version

cp config-example.json control.lua FactoCord-3.0 INSTALL.md LICENSE README.md SECURITY.md ./release/FactoCord3

zip -q ./release/FactoCord3.zip -r ./release/FactoCord3 || exit 3
echo "Created .zip archive"
tar -czf ./release/FactoCord3.tar.gz ./release/FactoCord3 || exit 4
echo "Created .tar archive"

