# CIS Class Server - v2.2.3

## Overview

This Go server provides a platform for delivering educational content, primarily focused on lesson materials, interactive exercises and management of offline clones of internet websites. It includes a web server, a command-line interface, and a Git integration module. This server is designed for educational use and prioritizes security and controlled access tailored to limited network access and absence of internet connectivity ( i.e. corrections ).

## Prerequisites

- **Go:** Go 1.18 or later must be installed on your system. You can download it from [https://go.dev/dl/](https://go.dev/dl/).
- **Dependencies:** The server relies on several Go packages. All Go modules are vendored for consistency across developers.
- **Environment Variables:** The server utilizes optional environment variables for configuration.

## Installation

1.  **Clone the Repository:**

    ```sh
    git clone --recursive http://192.168.1.28:3000/Go_DEV/go-class-server.git

    cd go-class-server
    ```

2.  **Build the Frontend**

    ```sh
    cd go-class-server-frontend
    npm run extract && npm run build && npm run move
    ```

3.  **Run the Server**
    ```sh
    cd ..
    go run .
    ```

## Configuration
### Environment Variables

Environment variables will override the hardcoded defaults if they are present. This feature was included to allow for development on a mock network. To change the defaults for these variables they should be edited in the `environment.go` file in the `util` package.

As of version 2.2.3 there is significant overlap between these variables and this has been identified as a target for code cleanup.

- RELEASE_VERSION: The current version number
- RELEASE_DATE: The date of current release
- WEBSITES_GITEA_ADDR: The query address and parameters for websites repo information
- UPDATE_IP: The root address for course content
- GIT_INSTALL_ADDR: The address for the user or organization hosting offline website repos
- WEBSITES_REPO_ADDR: The address to be queried for a single website meta
- CIS_CLASS_SERVER_PORT: The port the server will listen on (default is 22022)

## Running the Server

1.  **Start the Server:**

    ```sh
    go run main.go [options]
    ```

## Usage

### Web Interface

The server provides a simple web interface accessible at `http://localhost:<port>`. This interface provides the core functionality, offering access to lesson materials and interactive exercises.

## Key Features

- **Lesson Delivery:** Provides a structured platform for accessing educational content.
- **Git Integration:** Enables automatic updates and management of lesson materials via Git.
- **Command-Line Interface:** Offers a flexible way to interact with the server, useful for testing and debugging.
- **Development Mode:** Facilitates debugging and development.

## Security Considerations

- **Restricted Access:** The server is designed to be accessed only by the localhost. Contents of this server cannot be accessed by others on the local network preventing the transfer of data between incarcerated students. This can only be disabled by the developer and cannot be changed by any means once the code has been compiled.
- **Git Access Control:** Strictly control access to the Git repository to prevent unauthorized modifications. The application uses Git repositories as source of truth for all materials (websites, and course content), if changes are made to content on the local machine, once the update routine runs
  it will be overwritten with a forced pull request. This ensures that content is as secure as the policies used to manage the content repo.

## Contributing

We welcome contributions to this project. Please follow these guidelines:

1.  Fork the repository.
2.  Create a new branch for your changes.
3.  Commit your changes.
4.  Include a description of the changes and the intent behind them in a `CONTRIBUTION.MD` file committed to your branch.
5.  Submit a pull request.

## How to push out new release of the class server

To launch a new release, make sure you change the version number in the verson.json file located in the bin repo of the class server. Also ensure the version change has been reflected in version.ts in the frontend, environment.go in the backend, and the class server build.sh script (we should streamline this somehow).

Run the deploy script in the frontend from the go-class-server-frontend directory:

```sh
npm run deploy
```

This will build the frontend, move the built files to the proper spot in the backend to be embedded, and then call the build script for the backend. The backend build script will then build the class server with the frontend files embedded into the executable, and then automatically push the new changes to the repo on Gitea 28. 

*Note:* the version should be updated **before** executing the deploy script. It is the version number the launcher checks to initiate the update process.


## API and Routes

### `GET "/api/information/[filename]"`

- Sends the contents for the requested informational files located in `classServer/ClassServerResources/information`. e.g about.md

### `GET "/api/lessons"`

- Sends a json object of all available lessons to the frontend for processing.

### `GET "/api/data/[subfolder]/[filename]"`

- Sends the markdown for the desired lesson to the frontend to be rendered to HTML.

### `GET "/api/links"`

- Sends a json object in the following shape:

```json
{
  "name": "api-docs.deepseek.com", // name of offline site
  "indexPath": "C:/websites/api-docs.deepseek.com/index.html", // path on local system,
  "state": "installed", // state (not_installed, installed, need_update, or orphaned)
  "meta": {
    // metadata pulled from Gitea 47
    "size": 1044,
    "topics": [],
    "description": "The DeepSeek API uses an API format compatible with OpenAI. By modifying the configuration, you can use the OpenAI SDK or softwares compatible with the OpenAI API to access the DeepSeek API.",
    "created_at": "2025-10-28T12:54:13-07:00",
    "updated_at": "2025-11-03T14:16:38-08:00"
  }
}
```

### `GET "/"`

- Root route. Serves the index file for the front end.

### `GET "/api/git/[git command]/[repo]"`

- This route allows the running of three commands related to webpage maintenance:
  - `install`: Runs git clone to download the desired website.
  - `update`: Runs git push to update the website if needed.
  - `delete`: runs rm -rf to delete the desired website

### `GET "/websites/[filepath]"`

- This route will serve the index file supplied by the `indexPath` property of the links json.

### `POST "/websites/try.w3schools.com/[path]"`
  - This route allows the code examples to work on the offline version of W3S for TS, Perl, Python, and Go. This feature does require the download of a browser extension.

