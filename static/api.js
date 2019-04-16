'use strict';

import { getAccessToken, getPersonId } from './auth-okta.js';

export const getVideosHTML = (id, andThen) => {
  fetch(`/items?id=${id}`, {
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

export const getVideoList = (id, andThen) => {
  fetch(`/videoLists?id=${id}`, {
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
export const addVideoToList = (listId, video, andThen) => {};

export const putNewVideoList = (title, video, andThen) => {
  if (!getPersonId()) {
    alert('Please log in or register to create a video list');
    return;
  }
  const item = { creatorId: getPersonId(), videos: [video], title };
  fetch(`/videoLists`, {
    method: 'PUT',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + getAccessToken(),
    },
    redirect: 'follow', // manual, *follow, error
    referrer: 'no-referrer', // no-referrer, *client
    body: JSON.stringify(item), // body data type must match "Content-Type" header
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

export const putVideo = (id, item, andThen) => {
  if (!getPersonId()) {
    alert('Please log in or register to add videos');
    return;
  }
  fetch(`/videos?id=${id}`, {
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

export const postVote = (id, vote, andThen) => {
  fetch(`/vote?id=${id}`, {
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

export const deleteVote = (id, vote, andThen) => {
  fetch(`/vote?id=${id}`, {
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

export const getMetadata = (id, andThen) => {
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
