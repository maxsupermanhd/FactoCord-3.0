
rm -rf ./release

mkdir ./release
mkdir ./release/FactoCord3

cp config-example.json control.lua FactoCord3 INSTALL.md LICENSE README.md SECURITY.md ./release/FactoCord3

pushd ./release > /dev/null

pushd ./FactoCord3 > /dev/null
chmod 664 **
chmod +x ./FactoCord3
popd > /dev/null

zip -q ./FactoCord3.zip -r ./FactoCord3 || exit 3
echo "Created .zip archive"
tar -czf ./FactoCord3.tar.gz ./FactoCord3 || exit 4
echo "Created .tar archive"
popd > /dev/null
