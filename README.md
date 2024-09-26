# wget in go

wget is a command-line program written in Go that seeks to replicate the functionality of the free, open source tool, [`wget`](https://www.gnu.org/software/wget/manual/wget.html), used to retrieve content from web servers. It supports HTTP and HTTPS protocols.

## Features
To be more specific, the features implemented include:
- Downloading a file from a given url
- Allowing the user to set the download speed
- Downloading in the background
- Multiple file download
- Mirroring a website

## Prerequisites

To be able to run the program, you should have Go installed. You can download and install Go from the official [Go website](https://go.dev/dl/).

## Usage

To clone the projects:
```bash
    git clone https://learn.zone01kisumu.ke/git/nichotieno/wget.git
    cd wget
```

The structure of the command for using the program is:
```bash
go run ./cmd/app [flags] URL
```
Where flags (which are optional) can be any of:

 1. `-B` downloads a file immediately to the background with the output redirected to a log file. When a command with this flag is executed, it immediately logs to the terminal `Output will be written to "wget-log"`.
    ```bash
    $ go run . -B https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
    Output will be written to "wget-log".
    ```

 2. `-O` followed by the name you want to name the file. For example:
  
    ```bash
    $ go run . -O=meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
    ```
 3. `-P` followed by the path to where you want to save the file. For example:
    ```bash
    $ go run . -P=~/Downloads/ -O=meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
    ```
 4. `--rate-limit` followed by the speed you want to throttle your download to. For example:
    ```bash
    go run . --rate-limit=400k https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
    ```
    Where the rate limit can be specified in kilobytes, `k`, or in megabytes, `M`
 5. `-i` flag followed by a file name that will contain all links that are to be downloaded, where you want to download multiple files asynchronously. For example:
     ```bash
    $ ls
    download.txt   main.go
    $ cat download.txt
    http://ipv4.download.thinkbroadband.com/20MB.zip
    http://ipv4.download.thinkbroadband.com/10MB.zip
    $ go run . -i=download.txt
    ```
    **Note**: In this case, the URL will not be passed as the second argument.
 6. The `--mirror` falg can be used when you want to download a websites resources to be able to use parts of website offline. Some optional flags will go with --mirror. The basic syntax will be:
    ```bash
    go run . --mirror [mirror flags] https://example.com
    ```
    The optional `--mirror` flags include:
      - Directory-Based Limits  (`--reject` short hand `-R`). Tthis flag will have a list of file suffixes that the program will avoid downloading during the retrieval.
        ```bash
        go run . --mirror -R=jpg,gif https://example.com
        ```
     - Directory-Based Limits (`--exclude` short hand `-X`). This flag will have a list of paths that the program will avoid to follow and retrieve. So if the URL is https://example.com and the directories are /js, /css and /assets you can avoid any path by using -X=/js,/assets. The fs will now just have /css.
        ```bash
        $ go run . --mirror -X=/assets,/css https://example.com
        ```
    - Convert Links for Offline Viewing (`--convert-links`).  This flag will convert the links in the downloaded files so that they can be viewed offline, changing them to point to the locally downloaded resources instead of the original URLs.
        ```bash
        $ go run . --mirror --convert-links https://example.com
        ```
## Implementation

The main entry point of the program is located in `main.go`, which parses command-line arguments and sets values in an input struct to determine the desired operations. The program features several packages in `/internal` that contains functions to handle various functionalities. Highligted are some of the primary functions in each package

####  donwloader package

 - Single File Download: downloader.OneDownload(file, url, rateLimit, path) manages downloading a file with specified options and provides a progress bar with feedback on the download status.
 - Multiple File Downloads: The downloader.DownloadMultipleFiles(sourcefile, file, rateLimit, path) function reads URLs from a specified file and downloads them concurrently.
 - 
####  background package

 - Background Downloads: The background.DownloadInBackground(file, url, rateLimit) function allows users to download files in the background, logging output to a designated file. This includes capturing the start and finish time of the download, response status, and content size.
####  mirror package
 - Website Mirroring: The mirror.DownloadPage(url, flagInput) function retrieves the entire website, parsing HTML to find linked resources while following specified rules like excluding certain file types and directories.

#### fileManager package
 - Logging: The fileManager.Logger(file, url, rateLimit) function logs detailed information about the download process when running in background mode. This includes timestamps, request statuses, content sizes, and file paths, providing a comprehensive audit trail for all download activities.

## Contributors
This project was a collaboration of  three apprentices from [z01Kisumu](https://www.zone01kisumu.ke/). 
1. Nicholas Otieno
2. [Hillary Okello](https://github.com/HilaryOkello)
3. [Raymond Caleb](https://github.com/Raymond9734)
4. [Andrew Osindo](https://github.com/andyosyndoh) 


### License

This project is licensed under  [MIT License](./LICENSE)


