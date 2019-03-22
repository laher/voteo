voteo
-----

voteo is a toy app for watching videos together.

	vote + video = voteo

---------------------

WIP
---

voteo is a work in progress … technically, I'm trying a back-to-basics approach.

I'm keeping the functionality to a minimum, so that I can experiment with different approaches.


## Premise

> web dev has become complex. _So many layers OMG_

What does it look like when we remove some layers?

 * Minmise dependencies front and back end
 * What can we do now that ES6 has landed? (No jquery please)
 * What do we still need from our backend? Keep it minimal 
 * Question: what is the overhead of an SPA?
	 * Before SPAs, we had html pages with smatterings of dynamic content
	 * With SPAs, everything is dynamic. The whole page is rendered at once.
	 * In practice, I think the structure of many html pages are mostly static.
		 * Most elements are static 
		 * Some elements are shown/hidden (style.display) - the HTML isn't recreated each time.
		 * For some, only content is updated (text, input values).
		 * Often the only 'DOM generation' is for a list of items.
	 * What would it look like without React/Angular/etc
	 * If it's better - what's the threshold for creating a new app? Is there a happy medium?

## UI

 * Vanilla JS - no frameworks, limited libraries
 * No transpiler; serve libraries from CDNs
 * Embrace ES6, because browser support

## Backend concerns

 * DB, crud
 * Auth/Identity - validate tokens, protect endpoints, identity. 3rd party service …
 * Proxying
 * SSL

## Auth / Identity

 * Okta Widget
 * Plans to add Google/FB integration


-----------------------

Running
-------

Create a config.json as below and run the app.

At the moment it depends on okta auth - I'll build out the fallback-auth option soon.

```
	 {
		"auth": {
			"type": "okta",
			"okta": {
							"baseUrl": "https://XXXXXX.okta.com",
							"clientId": "XXXXXXXXXXXXXX",
							"redirectUri": "http://localhost:3000/",
							"authParams": {
							"issuer": "default",
								"responseType": ["id_token", "token"]
						}
			}
		},
		"ssl": false,
		"address": "localhost:3000"
	 }
```
