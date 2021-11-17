# Small simple cms

This project is made as an opportunity study.

## Planred

- pagination
- backup bucket
- restore bucket items
- upload bucket item
- file change history ( list, diff view)
- bucket item change history ( list, diff view)


## run files only
```
make cms
```
![tmpls](https://github.com/alexsuslov/cms/raw/main/cmd/server/static/images/run.jpg)

### Filemanager preview

![tmpls](https://github.com/alexsuslov/cms/raw/main/cmd/server/static/images/tmpls.jpg)


## run with Bolt DB example
```
make bolt
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

## Markdown Wiki
```
make wiki
```

