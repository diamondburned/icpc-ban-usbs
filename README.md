# icpc-ban-usbs

Hodgepodge to ban USB mass storage devices from working unless an upstream HTTP
server reports otherwise. This is a prototype/proof-of-concept and should not
be used in production!

## Setting Up

First, copy this entire repository (minus the `icpc-server` folder) into
`/etc/icpc`.  Then, install the needed files:

```sh
install -m640 systemd/icpc-usb-check-server.service /etc/systemd/system/
install -m640 udev/rules.d/99-icpc.rules            /etc/udev/rules.d/
```

> [!WARNING]
> After installation, USBs will NOT WORK BY DEFAULT! You will need to run a web
> server that the service can understand. See the *Testing* section below.

## Testing

To test, you can host `icpc-server`. To do this, run:

```sh
(cd icpc-server && go run .)
```

The server will expose a single endpoint, `/competition-status`, which returns
the current competition status. By default, `ongoing` is returned. To change
this, edit `icpc-server/main.go`'s `currentCompetitionStatus` to
`CompetitionFinished`. This will allow USBs to be mounted.
