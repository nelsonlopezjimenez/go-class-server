# Website management 2.0

## API request
- install route (clones given website)
- update route (force pulls given website)
- delete route (force deletes website folder on disk)

## Gitea Information

```sh
curl -X 'GET' \
  'http://{giteaURL}/api/v1/repos/search?uid={websites userid #}' \
  -H 'accept: application/json'
```

returns (omitted irrelevant fields)

```json
{
    "ok":boolean,
    "data":[
          {
            "id": 23,
            "owner": {
                "id": 1,
                "username": "MonkeyGoblin"
            },
            "name": "cis243-cis245",
            "full_name": "MonkeyGoblin/cis243-cis245",
            "description": "AM Class",
            "size": 13731,
            "clone_url": "http://localhost:3000/MonkeyGoblin/cis243-cis245.git",
            "default_branch": "main",
            "created_at": "2025-07-28T07:58:07-07:00",
            "updated_at": "2025-08-21T05:43:35-07:00",
            "topics": [],
        },

    ]
}
```
In the `Meta` object sent to the frontend, I would like to have the following fields from each website repo:
- name
- description
- size
- topics

## Local file structure

```sh
root/websites #(websites folder - no-git)
    |   websites.json # contains gitea data
    |
    +---website1 #(cloned repo)
    |       index.html
    |
    +---website2 #(cloned repo)
    |       index.html
    |
    \---website3 #(cloned repo)
            index.html


```

## Get and compare websites metadata

Conducting online check would negate the need to pull gitea data from remote if offline.
In the `offline links` object currently returned the `updated` value should be changed to `mayUpdate` and should default to `false` unless the user currently has a gitea connection and either `size` or `updated_at` differ from the local copy of _websites.json_.

> When considering how to update the local websites.json, I would prefer to automatically update all fields except for `size` and `updated_at`. These two fields should be excluded from automatic update because they are used to determine if the local copy of the website is out of date. They should only be updated when git commands are run (both initial clone and any subsequent pull request). This allows the user to decide if they want to run the update against a previously downloaded website.
> This method is more complicated than a single source of truth, however it allows the data like description to be updated without initiating a pull request and maintains the design patterns of allowing users to control which websites they want to download.

## some pseudocode

> _try_ `file from storage` => `local`

> _try_ `JSON from gitea` => `remote`

> _if_ `local` && `remote` then
>
> - `compare()`

> _if_ `local` && `!remote` then
>
> - `return local`

> _if_ `!local` && `remote` then
>
> - `write local`
> - `read local`
> - `return local`

> _if_ `!local` && `!remote` then
>
> - handle accordingly

## Migration script
```sh
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
        git remote add origin http://localhost:3000/MonkeyGoblin/$i;   #adds origin based on folder name
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

```