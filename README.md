#floppy 
##a personal file transfer system

I created this project so that I could easily send small and large files between my laptop, computers, and phone without having to follow the restirctions of cloud services such as Drive or Dropbox. 

Currently, the project contains a CLI and Daemon. This was inspired by Docker. The daemon runs on a host computer and recieves requests from clients who run the CLI locally. THe CLI has access to the daemons API which recieves and stores files locally.
Any client can upload/download/delete files on the daemon.

This project uses only the standard go library and cobra for the CLI. 
