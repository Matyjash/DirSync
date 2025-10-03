# DirSync

A simple tool for synchronizing files and directories between two locations.

## Usage

DirSync copies files from source directory to destination directory if the files do not exist in the destination path or are not the same.

```sh
dirsync /path/to/source_dir /path/to/destination_dir 
```

### Source based synchronization with deletion

To enable an additional feature of deleting files and directories in the destination that do not exist in the source, use the `--delete-missing` flag at the start of the command:

```sh
dirsync --delete-missing /path/to/source_dir /path/to/destination_dir 
```

## Building

To build the project, ensure you have Go installed and run:

```sh
go build ./cmd/main.go
```

## Testing

To run the tests, use:

```sh
go test ./...
```
