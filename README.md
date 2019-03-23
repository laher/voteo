voteo
-----

voteo is a toy app for watching videos together.

	vote + video = voteo

---------------------

WIP
---

voteo is a work in progress … technically, I'm trying a back-to-basics approach.

I'm keeping the functionality to a minimum, so that I can experiment with different approaches.

## Technical Premise

A modern web app, _without_ the SPA (Single Page App).

### Background

What does a non-SPA web app look like nowadays?

 * Before SPAs, we had html pages with smatterings of dynamic content
 * With SPAs, everything is dynamic. The whole page is rendered at once.

For a long while, I have been working with SPAs, and during that time browsers have improved in all kinds of ways...

What can we do now that ES6 has landed?
 * compact function declarations
 * Arrays are powerful - map/filter/find/sort
 * fetch, async/await
 * document.querySelector()
 * modules

### The Question

 * Question: what is the overhead of an SPA?

   In practice, I think the structure of many html pages are mostly static. They have interactive components, but don't necessarily need the full power of the SPA.

In a typical app:

   * Most elements are static 
   * Some elements are shown/hidden (style.display) - the HTML isn't recreated when you interact with the page.
   * For some elements, only content is updated (text, input values).
   * Often, the only 'DOM generation' is for a list of items.
		
 * What would it look like without React/Angular/etc? What problems need to be solved which are usually solved by the framework?

 * If it is easier - what's the threshold for creating a new app? Is there a happy medium?

## Approach

I have tried a few different approaches, but the general idea is as follows:

### UI

 * Vanilla JS - no frameworks, limited libraries
 * No transpiler; serve libraries from CDNs
 * Embrace ES6, because browser support
 * CSS: Flexbox - no need for bootstrap etc
 * Modules? ES6 has modules. Some quirks but I mostly like them. See ['js-modules' tag](https://github.com/laher/voteo/tree/superfine)
 * DOM manipulation???? Try some things:

#### DOM manipulation

This is the area I was least certain about ... I've tried a few approaches:

  * string interpolation (super easy with ES6)
  * Virtual DOM library (superfine) - see ['superfine' tag](https://github.com/laher/voteo/tree/superfine)
  * fetch rendered html from server. See ['templates' tag](https://github.com/laher/voteo/tree/templates)

I think the 'virtual DOM' approach is overkill for this case, and probably for most apps which aren't fully SPA.

I do like the string interpolation, but rendered HTML has some great advantages. The template can be delivered as part of the initial page load, and later separately. It feels snappy and lightweight.

### Backend concerns

What do we need from our backend? 
 * Persistence (DB. CRUD operations). Temporarily this is just files and a mutex
 * Auth/Identity - validate tokens, protect endpoints, identity. 3rd party service …
 * Config
 * SSL?
 * Proxying?

### Auth / Identity

Auth is non-trivial, I don't want to build my own auth for a small app like this. 

I have introduced a 3rd party tool for the time being - the Okta widget - it seems fine.

In due course Okta can provide Google/Facebook/MFA/...


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
