<h1 align="center">ISOGo</h1>
<p align="center">
ISOGo is an small utility to automate the download of OS ISO images for safekeeping on a server.
</p>

<hr> 

## Configuration

The config file [config/isogo.yml](https://github.com/kamushadenes/isogo/blob/master/config/isogo.yml) should be self-explanatory.

## Usage
 ```
 Usage of ./isogo:
  -config string
        the YAML config file (default "isogo.yml")
  -download
        download ISOs
  -keep
        run keep jobs
 ```

## Building

```
go build -o isogo
```

Or download one of the auto-built [releases](https://github.com/kamushadenes/isogo/releases).

### SystemD

To run it periodically through systemd, place the `isogo` binary you built on `/usr/local/bin` and run:

```
sudo mkdir /etc/isogo

sudo cp config/isogo.yml /etc/isogo

sudo cp systemd/isogo.service /etc/systemd/system/isogo.service
sudo cp systemd/isogo.timer /etc/systemd/system/isogo.timer

sudo systemctl enable isogo.timer
```

Make sure to edit [isogo.yml](https://github.com/kamushadenes/isogo/blob/master/config/isogo.yml) and [isogo.timer](https://github.com/kamushadenes/isogo/blob/master/systemd/isogo.timer) to fit your needs.

## Show your support

Give a ⭐️ if this project helped you!

<a href="https://www.buymeacoffee.com/kamushadenes"> <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" height="32" alt="kamushadenes" /></a>
