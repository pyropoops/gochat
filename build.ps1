
echo "Building..."
cd client
npm run build
echo "React Frontend built!"
cd ..
del server/views -Recurse
copy -r client/build server/views
cd server
go build
echo "Native Executable Built!"
cd ..
del server.exe
move server/server.exe server.exe