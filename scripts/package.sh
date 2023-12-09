# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

evilarc_version=$(./bin/evilarc -v)
echo "build version: $evilarc_version"

# cross_compiles
make -f ./scripts/Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release || exit

for os in $os_all; do
    for arch in $arch_all; do
      evilarc_dir_name="${evilarc_version}_${os}_${arch}"
      evilarc_path="./packages/${evilarc_version}_${os}_${arch}"

      if [ "x${os}" = x"windows" ]; then
          if [ ! -f "./evilarc_${os}_${arch}.exe" ]; then
              continue
          fi
          mkdir ${evilarc_path}
          mv ./evilarc_${os}_${arch}.exe ${evilarc_path}/evilarc.exe
      else
          if [ ! -f "./evilarc_${os}_${arch}" ]; then
              continue
          fi
          mkdir ${evilarc_path}
          mv "./evilarc_${os}_${arch}" "${evilarc_path}/evilarc"
      fi

      # packages
      cd ./packages || exit
      if [ "x${os}" = x"windows" ]; then
          zip -rq ${evilarc_dir_name}.zip ${evilarc_dir_name}
      else
          tar -zcf ${evilarc_dir_name}.tar.gz ${evilarc_dir_name}
      fi
      cd ..
      rm -rf "${evilarc_path}"
    done
done

cd - || exit