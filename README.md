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

## Working Mechanism

We're using `udev` to run a script that directly disables the USB in the kernel
if the competition is ongoing. This is done in `udev/icpc-usb-check`.

However, because of [systemd-udevd's sandboxing](https://sandboxdb.org/service/systemd-udevd.service.html),
we cannot use networking in the udev script. To work around this, we use a
one-shot systemd service that calls `wget`, then call `systemctl start --now`
in the udev script. This avoids the network call entirely in the udev script.
It only works because `systemctl` will block until the one-shot service is
done.

### Diagram

```
┌─ ICPC competition machine ─────────────────┐
│                                            │
│  ┌─────────┐    ┌─────────┐                │
│  │   USB   ├────►  Linux  │                │
│  └─────────┘    └────┬──▲─┘                │
│                      │  │                  │
│       ┌──────────────┘  │                  │
│       │                 │                  │
│  ┌────▼────┐    ┌───────┴────────┐         │
│  │  udev   ├────► icpc-usb-check │         │
│  └─────────┘    └───────┬────────┘         │
│                         │                  │
│       ┌─────────────────┘                  │
│       │                                    │
│  ┌────▼────┐    ┌───────────────────────┐  │
│  │ systemd ├────► icpc-usb-check-server │  │
│  └─────────┘    └───────────┬───────────┘  │
│                             │              │
│                    ┌────────┘              │
│                    │                       │
│                 ┌──▼───┐                   │
│                 │ wget │                   │
│                 └──┬───┘                   │
│                    │                       │
└─────────────────── │ ──────────────────────┘
                     │
┌─ ICPC server ───── │ ──────────────────────┐
│                    │                       │
│  ┌─────────────────▼────────────────────┐  │
│  │ Competition Server                   │  │
│  │   - code submission                  │  │
│  │   - competition status report        │  │
│  └──────────────────────────────────────┘  │
│                                            │
└────────────────────────────────────────────┘
```
