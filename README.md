#shape-sql-loader
shape-sql-loader is intended to be a sysadmin tool to upload a set of shp files to the new NSIv2.6 database. Upload environment must be accessible to the PostGIS database instance. Compiled executable requires ogr2ogr and assets/metadataBase.xlsx

Usage:

    TODO replace 'go run .' with compiled executable name
    go run . --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp
    go run . --mode upload --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp --xlsPath /workspaces/shape-sql-loader/metadata.xls --sqlConn "host=host.docker.internal port=25432 user=admin password=notPassword database=gis"
    go run . --mode access  --datasetId randomguid --group nsi --role admin

Database setup and cleanup SQL scripts are stored in scripts/sql/. All tables must be created inside the 'nsi' database schema

    1. Generate metadata template using "--mode prep"
        go run . --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp

    2. Fill in metadata xls file and upload with "--mode upload"

    3. Add access groups

Bonus VIM config: Delve can be used to start a headless debug server inside a container and connected from the local environment using vimspector. The attached .vimspector config can "docker exec" into a container and start the dlv server automatically. Alternatively, the debugger can be start and attached manually:

    dlv debug --listen=0.0.0.0:2345 --api-version=2 --log --log-output=dap --headless -- --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp
