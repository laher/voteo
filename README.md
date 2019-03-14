voteo
-----

voteo is a toy app for watching videos together.

	vote + video = voteo

---------------------

WIP
---

voteo is a work in progress …

Technically, I'm trying a back-to-basics approach:

 * Minmise dependencies front and back end

## UI

 * Vanilla JS - no frameworks, limited libraries
 * No transpiler; deliver libraries via CDNs
 * Embrace ES6, because browser support

## Auth / Identity

 * Okta Widget
 * Plans to add Google/FB integration

-----------------------

Running
-------

Make a config.json with `{"address":"localhost:3000"}` and run the app.

At the moment it depends on okta auth - I'll build out the fallback-auth option soon.
