#!/usr/bin/env bash
BASEDIR=$(dirname "$0")
echo "Base directory is : $BASEDIR"
echo "Building..."
go build
echo "Built."
echo "Making executable..."
chmod +x $BASEDIR/podcatch
echo "Made."
echo "Copying to /usr/bin"
sudo cp $BASEDIR/podcatch /usr/bin/podcatch
echo "Copied."
echo "You should now be able to run podcatch as just 'podcatch'"
