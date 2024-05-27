# Floppy

## A Personal File Transfer System

Floppy is a personal file transfer system that allows you to easily send small and large files between your laptop, computers, and phone without the restrictions of cloud services such as Google Drive or Dropbox. 

The name is inspired by the floppy disk where you had your own files that were on a transferable extrnal disk!

### Overview

Floppy consists of two main components:
1. **CLI (Command Line Interface)**: A client application that allows you to interact with the Floppy daemon.
2. **Daemon**: A server application that runs on a host computer, receives requests from CLI clients, and manages file storage and transfer.

This project was inspired by Docker, where the daemon runs on a host computer and receives requests from clients running the CLI locally. The CLI has access to the daemon's API, which receives and stores files locally. Any client can upload, download, or delete files on the daemon.

### Features

- **File Upload**: Upload files from any client to the daemon.
- **File Download**: Download files from the daemon to any client.
- **File Deletion**: Delete files stored on the daemon.

### Technologies Used

- **Go**: The project is implemented using the Go programming language.
- **Cobra**: Used for building the CLI application.

This project is still in progress!
