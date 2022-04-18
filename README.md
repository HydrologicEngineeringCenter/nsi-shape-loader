# shape-sql-loader

shape-sql-loader is intended to be a sysadmin tool to upload a set of shp files
to the new NSIv2.6 database. PostGIS database instance must be accessible by the
upload environment. The tool requires ogr2ogr and assets/metadataBase.xlsx.

Database setup and cleanup SQL scripts are stored in scripts/sql/. All tables
must be created inside the 'nsi' database schema.

```golang
    1. Generate metadata template using "--mode prep"
        go run . --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp

    2. Fill in metadata xls file and upload with "--mode upload"
        go run . --mode upload --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15003.shp --xlsPath /workspaces/shape-sql-loader/metadatatest.xlsx --sqlConn "host=host.docker.internal port=25432 user=admin password=notPassword database=gis"

    3. Add access groups via "--mode access"
        go run . --mode access --group nsidev --role admin --user user_id --sqlConn "host=host.docker.internal port=25432 user=admin password=notPassword database=gis"
```

Bonus VIM config: Delve can be used to start a headless debug server inside a
container and connected from the local environment using vimspector. The
attached .vimspector config can "docker exec" into a container and start the
delve server automatically. Alternatively, the debugger can be start and
attached manually via the "remoteConnect" option:

```
    dlv debug --listen=0.0.0.0:2345 --api-version=2 --log --log-output=dap --headless -- --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp
```
