'use strict';

import { initOkta, getPersonId } from './auth-okta.js';

let state = {
  conf: null,
};

export const cleanInput = input => {
  const pref = 'https://www.youtube.com/watch?v=';
  if (input.startsWith(pref)) {
    return input.substring(pref.length);
  }
  return input;
};

export const showSignInOut = personId => {
  document.getElementById('app-container').style.display = 'block';
  document.getElementById('name').innerHTML = personId;
  if (personId) {
    console.log('show signout');
    document.getElementById('logged-in').style.display = 'flex';
    document.getElementById('logged-out').style.display = 'none';
  } else {
    console.log('show signin button');
    document.getElementById('logged-in').style.display = 'none';
    document.getElementById('logged-out').style.display = 'flex';
  }
};

// Listen on page load:

export const initAuth = andThen => {
  fetch(`/auth/settings`, {
    method: 'get',
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(JSON.stringify(json));
      state.conf = json;
      if (json['type'] == 'okta') {
        initOkta(json.okta, andThen);
      } else {
        // assume it's no-login
        // reFetch();
        // TODO
        andThen();
      }
    });
};
