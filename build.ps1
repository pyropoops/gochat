cd client
($env:REACT_APP_BACKEND = "ws://127.0.0.1:8000/ws") -and (npm run build)
echo "Front-end views built!"
cd ..
del server/views -Recurse
copy -r client/build server/views
cd server
go build
echo "Final native executable built!"
cd ..
del server.exe
move server/server.exe server.exe
echo "Done!"
