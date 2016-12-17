##### 1. Install and run the `docker-compose.yml`
```bash
git clone https://gist.github.com/1975674c22ce8948c895.git taiga
cd taiga
# Update docker-compose.yml
# - Replace Hostname of taigaback and taigafront
# - Update or disable Email settings
docker-compose up -d
```

##### 2. Wait a few seconds / minutes (generating sample data takes some time)

##### 3. Access Taiga via Taiga-Front's hostname and login with User `admin` and Password `123123`