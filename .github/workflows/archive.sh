
rm -rf ./release

mkdir ./release
mkdir ./release/FactoCord3

version=`git describe --tags`
version+=" ("
version+=`git rev-parse --short HEAD`
version+=")"

echo "Version: ${version}"
echo "${version}" > ./release/FactoCord3/.version

cp config-example.json control.lua FactoCord-3.0 INSTALL.md LICENSE README.md SECURITY.md ./release/FactoCord3

pushd ./release > /dev/null

pushd ./FactoCord3 > /dev/null
chmod 664 **
chmod +x ./FactoCord-3.0
popd > /dev/null

zip -q ./FactoCord3.zip -r ./FactoCord3 || exit 3
echo "Created .zip archive"
tar -czf ./FactoCord3.tar.gz ./FactoCord3 || exit 4
echo "Created .tar archive"
popd > /dev/null
