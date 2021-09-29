
echo "Building..."
cd client
npm run build
echo "React Frontend built!"
cd ..
del server/views -Recurse
cp -r client/build server/views
cd server
go build
echo "Native Executable Built!"
cd ..
rm server
mv server/server server
