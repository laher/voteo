'use strict';

import { initOkta, getAccessToken, getPersonId } from './auth-okta.js';
import { getRender } from './sf.js';

let state = {
  items: [
    { id: 'X7hFERntlog', title: 'Fearless Org' },
    { id: 'd_HHnEROy_w', title: 'Stop managing' },
    { id: 'BCkCvay4-DQ', title: 'Foos' },
  ],
  votes: [{ videoId: 'x7hFERntlog', personId: 'am', up: true }],
  selectedItem: 'X7hFERntlog',
};

export const countVotes = id => {
  return (
    state.votes.filter(vote => vote.videoId == id && vote.up).length -
    state.votes.filter(vote => vote.videoId == id && !vote.up).length
  );
};

export const haveIUpvoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && vote.up && vote.personId == getPersonId()
    ).length > 0
  );
};

export const haveIDownvoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && !vote.up && vote.personId == getPersonId()
    ).length > 0
  );
};

export const haveIVoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && vote.personId == getPersonId()
    ).length > 0
  );
};

const itemsHTMLFetch = () => {};

const itemsHTML = () => {
  return state.items
    .sort((a, b) => countVotes(b.id) - countVotes(a.id))
    .map(
      i =>
        `
    <li><div>
    <div onclick="setSelectedItem('${i.id}')">${i.title} (${countVotes(
          i.id
        )} votes)</div>
    <div>` +
        (getPersonId()
          ? (haveIUpvoted(i.id)
              ? `<img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('${
                  i.id
                }')" />`
              : `<img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote('${
                  i.id
                }')" />`) +
            (haveIDownvoted(i.id)
              ? `<img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('${
                  i.id
                }')" />`
              : `<img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote('${
                  i.id
                }')" /> `)
          : `<abbr title='log in to vote'><img src="https://img.icons8.com/material/24/000000/question.png" onclick="alert('log in to vote')"></abbr>`) +
        `</div>
    </div></li>`
    )
    .join('');
};

export const setSelectedItem = i => {
  state.selectedItem = i;
  show(i);
};

const show = i => {
  if (!i) {
    return;
  }
  const link = 'https://www.youtube.com/embed/' + i;
  const player = document.getElementById('player');
  if (player.getAttribute('src') == link) {
    return; // no need to refresh
  }
  player.setAttribute('src', link);
  const item = state.items.find(item => item.id === i);
  let title = 'preview';
  if (item) {
    title = item.title;
  }
  document.getElementById('title').innerHTML = title;
};

export const upvote = i => {
  const item = state.items.find(item => item.id === i);
  let vote = state.votes.find(
    vote => vote.personId == getPersonId() && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: getPersonId(), videoId: i };
  }
  vote.up = true;
  postVote(vote);
};

export const unvote = i => {
  let vote = state.votes.find(
    vote => vote.personId == getPersonId() && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: getPersonId(), videoId: i };
  }
  deleteVote(vote);
};

export const downvote = i => {
  const item = state.items.find(item => item.id === i);
  let vote = state.votes.find(
    vote => vote.personId == getPersonId() && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: getPersonId(), videoId: i };
  }
  vote.up = false;
  postVote(vote);
};

const getVideosHTML = vidList => {
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
      vidList.innerHTML = body;
    });
};

const getVideos = () => {
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
      state.items = json;
      reflop();
    });
};

export const reFetch = () => {
  getVotes();
  getVideos();
};

const getVotes = () => {
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
      state.votes = json;
      reflop();
    });
};

const cleanInput = input => {
  const pref = 'https://www.youtube.com/watch?v=';
  if (input.startsWith(pref)) {
    return input.substring(pref.length);
  }
  return input;
};

const putVideos = () => {
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
    body: JSON.stringify(state.items), // body data type must match "Content-Type" header
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
      state.items = json;
      reflop();
    })
    .catch(console.error);
};

const postVote = vote => {
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
      state.votes = json;
      reflop();
    });
};

const deleteVote = vote => {
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
      state.votes = json;
      reflop();
    });
};

export const preview = () => {
  const id = cleanInput(document.getElementById('addbox').value);
  if (id) {
    show(id);
  }
};

export const add = () => {
  console.log('add');
  const id = cleanInput(document.getElementById('addbox').value);

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
      let title = json.title;
      if (title.length > 30) {
        title = title.substring(0, 30) + ' ...';
      }
      state.items.push({ id: id, title: title });
      state.selectedItem = id;
      putVideos();
    });
};
const itemRenderer = 'htmlFetch';
const reflop = async () => {
  var t0 = performance.now();
  const vidList = document.getElementById('videoList');
  if (vidList) {
    switch (itemRenderer) {
      case 'superfine':
        const render = getRender(vidList);
        render(state.items);
        break;
      case 'htmlString':
        vidList.innerHTML = itemsHTML();
        break;
      case 'htmlFetch':
        getVideosHTML(vidList);
        break;
      default:
        // no rendererer
        console.log('warning - no item renderer');
    }
    document.getElementById('videoCount').innerHTML = state.items.length;
  }
  show(state.selectedItem);

  var t1 = performance.now();
  console.log('Call to reflop took ' + (t1 - t0) + ' milliseconds.');
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

const start = () => {
  window.upvote = upvote;
  window.downvote = downvote;
  window.unvote = unvote;
  window.setSelectedItem = setSelectedItem;
  window.preview = preview;
  window.add = add;
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
        initOkta(json.okta);
      } else {
        // assume it's no-login
        reFetch();
      }
    });
};

// Listen on page load:
window.addEventListener('load', start);
/*
document.addEventListener('DOMContentLoaded', () => {
  console.log('adding event listeners');
  document.getElementById('add').addEventListener('click', () => add);
  document.getElementById('addbox').addEventListener('change', () => preview);
  document
    .getElementById('sign-in')
    .addEventListener('click', () => showSignInModal);
  console.log('done adding event listeners');
}); */
