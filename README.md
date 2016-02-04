# go-imgur
A basic api client for [imgur](https://api.imgur.com/). Currently I have only implemented the features of the imgur api that I use for other projects.

## Installation:
    go get github.com/dmashuda/go-imgur

## Example Usage:
#### Creating a Client:

    client := imgur.NewClient(clientID)
*clientID is issued by imgur on a per application basis*

#### Retrieving album information:

    aww, err := client.GetAlbum("/gallery/r/CorgiGifs", 0, 20)
