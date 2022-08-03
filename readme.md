# What's My Schedule

## Setup

### Creating Google OAuth Credentials

What's My Schedule retrieves your event schedule from the Google Calendar API by verifying the binary using OAuth Credentials.
To set up WMS you will have to create and retrieve API credentials to authenticate an OAuth flow.

1. [Create a New Project](https://console.cloud.google.com/projectcreate) in the Google Developer Console
2. [Enable the Google Calendar API](https://console.cloud.google.com/apis/library/calendar-json.googleapis.com) for this project
3. [Create a new OAuth Consent Screen](https://console.cloud.google.com/apis/credentials/consent)
   1. Select type: `external`
   2. Enter the following required fields:
      1. App name e.g. `wms-cli`
      2. User support email e.g. `ex@gmail.com`
      3. Developer contact email address e.g. `ex@gmail.com`
   3. Save and Continue to `Scopes`
   4. Add the following scopes:
      1. `/auth/calendar.calendarlist.readonly` non-sensitive scope
      2. `/auth/calendar.readonly` sensitive scope
      3. `/auth/calendar.events.readonly` sensitive scope
   5. Save and Continue to `Summary`
   6. Add your email under `Test users`
4. [Create OAuth Client Credentials](https://console.cloud.google.com/apis/credentials)
   1. Application type: `Desktop App`
   2. Save your `Client ID` and `Client Secret`

### Setting up WMS

1. Download the binary for your machine from the `release` section
2. Add the `wms` binary to your path, I recommend `/usr/local/bin`
3. Run `wms init <client-id> <client-secret>`

You're ready to go! Run `wms` to get today's calendar!
