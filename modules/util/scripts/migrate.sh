#!/usr/bin/sh



rm -d */;                                       #removes all empty directories
rm -rf .git;                                    # removes version control from folder
rm .gitmodules;                                 #removes the .gitmodules file creates by superproject
for i in */;                                    #start loop
    do cd $i;                                       #move into each directory
    echo "Checking $PWD";                           # log whats being done 
        if test -f .git                                 #check to see if the .git FILE is in the directory (normal repos have a directory not a file)
        then                                            # if the directory is in fact a submodule then migrate it. else skip it. 
        echo "submodule found. Migrating $i";
        rm -r .git --force;                             # removes version control from folder
        git init -b deleteme;                           # creates a dummy branch that will disappear once a real branch is check out
        git remote add origin http://192.168.1.47:3000/websites/$i;   #adds origin based on folder name
        git fetch origin main;                          #pulls existing branches
        git checkout main -f;                           #creates main branch and checks it out
    else echo "no module found in $i";
        git fetch;                                      #This fires if the migration was successfull or if it was interrupted, 
        git checkout main -f                            #and then forces checkout of the main branch. 
    fi;                                                 #This catches cases where migration did not complete because of disconnection.
    cd $OLDPWD;
done;


## I tested this by alternating it with the reinit.sh script
# I tested disconnection by interrupting the gitea process during migration then running migration a second time.