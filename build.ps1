if ($args.count -eq 0) 
{
    Write-Output "./build.ps1 <accessible IP address>"
    exit
}

$address = $args[0] + ":8000"

echo "Building..."
cd client
($env:REACT_APP_BACKEND = "ws://$address/ws") -and (npm run build)
echo "React Frontend built for: $address"
cd ..
del server/views -Recurse
copy -r client/build server/views
cd server
go build
echo "Native Executable Built!"
cd ..
del server.exe
move server/server.exe server.exe