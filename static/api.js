'use strict';

import { getAccessToken, getPersonId } from './auth-okta.js';

export const getVideosHTML = andThen => {
  fetch(`/items`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization: 'Bearer ' + getAccessToken(),
    },
  })
    .then(function(response) {
      console.log(response);
      return response.text();
    })
    .then(function(body) {
      andThen(body);
    });
};

export const getVideos = andThen => {
  fetch(`/videos`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization: 'Bearer ' + getAccessToken(),
    },
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(json);
      andThen(json);
    });
};

export const getVotes = andThen => {
  fetch(`/vote`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization: 'Bearer ' + getAccessToken(),
    },
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(json);
      andThen(json);
    });
};

export const createNewList = andThen => {
  if (!getPersonId()) {
    alert('Please log in or register to create a video list');
    return;
  }
};

export const putVideos = (items, andThen) => {
  if (!getPersonId()) {
    alert('Please log in or register to add videos');
    return;
  }
  fetch(`/videos`, {
    method: 'PUT',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + getAccessToken(),
    },
    redirect: 'follow', // manual, *follow, error
    referrer: 'no-referrer', // no-referrer, *client
    body: JSON.stringify(items), // body data type must match "Content-Type" header
  })
    .then(function(response) {
      console.log(response);
      if (response.ok) {
        return response.json();
      } else {
        throw Error(`Request rejected with status ${response.status}`);
      }
    })
    .then(function(json) {
      console.log(json);
      andThen(json);
    })
    .catch(console.error);
};

export const postVote = (vote, andThen) => {
  fetch(`/vote`, {
    method: 'POST',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + getAccessToken(),
    },
    redirect: 'follow', // manual, *follow, error
    referrer: 'no-referrer', // no-referrer, *client
    body: JSON.stringify(vote), // body data type must match "Content-Type" header
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(json);
      andThen(json);
    });
};

export const deleteVote = vote => {
  fetch(`/vote`, {
    method: 'DELETE',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + getAccessToken(),
    },
    redirect: 'follow', // manual, *follow, error
    referrer: 'no-referrer', // no-referrer, *client
    body: JSON.stringify(vote), // body data type must match "Content-Type" header
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(json);
      andThen(json);
    });
};

export const getMetadataAndThen = (id, andThen) => {
  fetch(`/yt/data?id=${id}`, {
    method: 'get',
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(JSON.stringify(json));
      console.log(json.title);
      andThen(json);
    });
};
