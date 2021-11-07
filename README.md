# Small simple cms

This project is made as an opportunity study.

## Planred

- sorting
- filtering
- pagination
- backup backet
- restore backet items
- upload backet item
- file change history ( list, diff view)
- backet item change history ( list, diff view)


## run files only
```
(cd cmd/server; go run server.go)
```
![tmpls](https://github.com/alexsuslov/cms/raw/main/cmd/server/static/images/run.jpg)


## run with BoltDB
```
(cd cmd/serverdb; go run server.go)
```

### set default admin user 

cmd/serverdb/.env
```
ADMIN_USER=admin
ADMIN_USER_PASS=admin
```

### create new bucket item
just open http://localhost:8080/admin/buckets/test/test
![tmpls](https://github.com/alexsuslov/cms/raw/main/cmd/serverdb/static/images/new.jpg)

add 
```
{"name":"test"}
```
then press cmd+s(ctr+s) to save

## Filemanager preview

![tmpls](https://github.com/alexsuslov/cms/raw/main/cmd/server/static/images/tmpls.jpg)
