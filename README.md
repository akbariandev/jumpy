# Jumpy ( DAG Based Chain)

How to run one node:

```sh
go run cli/main.go run -p {port}
```

To use as an application you can run below commands in terminal: 

```shell
    log // print chain blocks
    transaction:{data} // add transaction to last block
    commit //commit last block and push to chain
```