    go run . -u admin -p notPassword -n gis -t host.docker.internal -o 25432 -a nsi_v2022 -d /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp
    ./upload -s "host=host.docker.internal port=25432 user=admin password=notPassword database=gis" -d . -c nsi -t table
    ogrinfo -so -al test/nsi/NSI_V2_Archives/V2022/15001.shp

    go run . --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp
    go run . --mode upload --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp --xlsPath /workspaces/shape-sql-loader/metadata.xls --sqlConn "host=host.docker.internal port=25432 user=admin password=notPassword dbname=gis"
    go run . --mode access  --datasetId randomguid --group nsi --role admin

To start debugger inside container + listen on mapped port

    dlv debug --listen=0.0.0.0:2345 --api-version=2 --log --log-output=dap --headless -- --mode prep --shpPath /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp

    https://github.com/golang/vscode-go/blob/bc8b63e62e6dcbfbd60e6e9f71eec7aff9ae367a/src/goDebugConfiguration.ts#L31-L114
