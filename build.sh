ECHO "checking for binary repo"
if checkForBinRepo != 0; then echo "Failed to build binary"; exit 1;
ECHO "Building the latest binary now";
go build -v -o bin/classServer.exe;
cd bin;
ECHO $PWD;
ECHO "Preparing to commit";
git add classServer.exe;
git commit -m "v2.3.1"
ECHO "Preparing to push new binary to remote repo"
ECHO "Pushing to origin"
git push -f --progress origin development:main;

function checkForBinRepo() {
if ! test -d bin; 
then 
echo "WARN: No bin directory detected";
echo "WARN: Please make sure you are connected to the network";
echo "WARN: If there is no network, build will fail without bin directory";
git clone http://192.168.1.28:3000/ClassroomResources/classServer.git bin;
if [ $? != 0 ]; then echo "Cloning repo failed"; exit 1; fi;
fi;

}