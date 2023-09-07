[![Go](https://github.com/TurningOfTables/planespotter/actions/workflows/go.yml/badge.svg)](https://github.com/TurningOfTables/planespotter/actions/workflows/go.yml)

# planespotter v0.1

## Description

One of many 'planespotter' apps, which tells you when a new plane enters the area around a given location according to ADS-B data and pops up a notification. Also keeps a running total of how many planes you have ever spotted, and their callsigns.

I used this as practice for consuming different API structures, and it tied in nicely with my love of aviation :)

Uses [The OpenSky Network](https://opensky-network.org)

![Screenshot](readme_screenshot.png)

Written in Go using [Beeep](https://github.com/gen2brain/beeep) for notifications and [godotenv](https://github.com/joho/godotenv) for configuration from env vars.

Lots of things that need tidying up a bit from this initial version, and improvements (see below)

## Parts that are still in progress

* Handling the lack of a provided OpenSky auth details by using the (more limited) anonymous access allowed by the OpenSky API
* Streamlining of lots of the code - in particular the saving, the fact we maintain an in-memory short term 'seen' list which was rendered pointless by the save file system added later, refactoring much of the setup out of main()
* Moving more out of the source into the .env configuration, like the frequency of checks and the search area size
* Improving and adding to the currently basic test coverage, particularly for the helper functions

## Configuration

1. Create an OpenSky account at [The OpenSky Network](https://opensky-network.org/)'s site
2. Create a copy of the `.env.dist` file at the root of the project, and name it `.env`
3. Add a decimal longitude and latitude under `LONGITUDE=` and `LATITUDE=`
4. Add your OpenSky details under `OPENSKY_USERNAME=` and `OPENSKY_PASSWORD=` (no, it's not base64 encoded and the API just does basic auth, so remember to *keep your details secure*, don't commit them to source control and perhaps remove them from the .env file afterwards)

After configuration your `.env` file might look like this if you wanted to see planes near Paris:

```
LATITUDE=48.864716
LONGITUDE=2.349014
OPENSKY_USERNAME=YourUserName
OPENSKY_PASSWORD=YourGreatPassword
```

## Usage

Run the correct binary for your platform from `/bin`, e.g ```./planespotter_win_amd64.exe```
Or build your own from the source with `go build` :)

To the best of my knowledge this should work but I've only been able to test Windows binaries and notifications so far.

The console will log out a few lines on each check: 

```
2023/09/07 11:51:14 Updating data...
2023/09/07 11:51:14 Remaining API requests today: 3911
2023/09/07 11:51:16 Waiting...
```

Your progress will be saved in `./save.json`

```
{"SeenCount":3,"Callsigns":["AHY007","SHT8Q","EIN182"]}
```

And you should start seeing notifications about flights you haven't previously seen

![Screenshot](readme_screenshot.png)

## Attributions

* All data is from [The OpenSky Network](https://opensky-network.org)
* Notifications use [Beeep](https://github.com/gen2brain/beeep)
* Configuration uses [godotenv](https://github.com/joho/godotenv)
* Icon uses [Airplane icons created by Mihimihi - Flaticon](https://www.flaticon.com/free-icons/airplane)
