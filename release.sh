
rm -rf ./release

go build || exit 1

echo "Build successful"

mkdir ./release
mkdir ./release/FactoCord3

git describe --tags > ./release/.version

if grep -q "-g" ./release/.version; then
  echo "No tag"
  exit 2
fi

cp config-example.json control.lua FactoCord-3.0 INSTALL.md LICENSE README.md SECURITY.md ./release/FactoCord3

zip ./release/FactoCord3.zip -r ./release/FactoCord3
tar -czvf ./release/FactoCord3.tar.gz ./release/FactoCord3

