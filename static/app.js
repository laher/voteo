'use strict';

import { initOkta, getPersonId } from './auth-okta.js';
import {
  getVideosHTML,
  getVideos,
  getVotes,
  putVideos,
  postVote,
  deleteVote,
  getMetadataAndThen,
  createNewList,
} from './api.js';

let state = {
  items: [
    { id: 'X7hFERntlog', title: 'Fearless Org' },
    { id: 'd_HHnEROy_w', title: 'Stop managing' },
    { id: 'BCkCvay4-DQ', title: 'Foos' },
  ],
  votes: [{ videoId: 'x7hFERntlog', personId: 'am', up: true }],
  selectedItem: '',
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

export const setSelectedItem = i => {
  history.pushState({ video: i }, 'title 1', '#v=' + i);
  console.log('history', history.state);
  loadSelectedItem({ state: history.state });
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
  const d = document.querySelector(`div[data-id='${i}']`);
  d.querySelectorAll('.votingSpan').forEach(item => {
    item.style.display = 'none';
  });
  d.querySelector('.loadingSpan').style.display = 'block';
  document.querySelectorAll('div.video-item').forEach(item => {
    item.classList.add('hidden');
  });
  postVote(vote, json => {
    state.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
};

export const unvote = i => {
  let vote = state.votes.find(
    vote => vote.personId == getPersonId() && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: getPersonId(), videoId: i };
  }
  const d = document.querySelector(`div[data-id='${i}']`);
  d.querySelectorAll('.votingSpan').forEach(item => {
    item.style.display = 'none';
  });
  d.querySelector('.loadingSpan').style.display = 'block';
  document.querySelectorAll('div.video-item').forEach(item => {
    item.classList.add('hidden');
  });
  deleteVote(vote, json => {
    state.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
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
  const d = document.querySelector(`div[data-id='${i}']`);
  d.querySelectorAll('.votingSpan').forEach(item => {
    item.style.display = 'none';
  });
  d.querySelector('.loadingSpan').style.display = 'block';
  document.querySelectorAll('div.video-item').forEach(item => {
    item.classList.add('hidden');
  });
  postVote(vote, json => {
    state.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
};

const updateVidList = body => {
  document.getElementById('videoList').innerHTML = body;
  addHandlers();
};

export const reFetch = () => {
  getVotes(() => {
    state.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
  getVideos(() => {
    state.items = json;
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

export const preview = () => {
  const id = cleanInput(document.getElementById('addbox').value);
  if (id) {
    show(id);
  }
};

export const newList = () => {
  console.log('new list');
  const id = cleanInput(document.getElementById('addbox').value);
  getMetadataAndThen(id, json => {
    let title = json.title;
    if (title.length > 30) {
      title = title.substring(0, 30) + ' ...';
    }
    state.items.push({ id: id, title: title });
    state.selectedItem = id;
    createNewList(state.items);
  });
};

export const add = () => {
  console.log('add');
  const id = cleanInput(document.getElementById('addbox').value);
  getMetadataAndThen(id, json => {
    let title = json.title;
    if (title.length > 30) {
      title = title.substring(0, 30) + ' ...';
    }
    state.items.push({ id: id, title: title });
    state.selectedItem = id;
    putVideos(state.items, () => {
      state.items = json;
      reflop();
    });
  });
};

const reflop = async () => {
  var t0 = performance.now();
  const vidList = document.getElementById('videoList');
  if (vidList) {
    getVideosHTML(updateVidList);
    document.getElementById('videoCount').innerHTML = state.items.length;
  }
  show(state.selectedItem);
  var t1 = performance.now();
  console.log('Call to reflop took ' + (t1 - t0) + ' milliseconds.');
};

const addHandlers = () => {
  console.log('set up selectors / event listeners');
  const nodeList = document.querySelectorAll('div.video-item');
  //const nodeList = document.querySelectorAll('div.video-item');
  console.log('found ', nodeList.length, ' items');
  Array.from(nodeList).forEach(item => {
    const videoId = item.getAttribute('data-id');
    item.querySelector('.title').addEventListener('click', e => {
      setSelectedItem(videoId);
    });
    item.querySelector('.upvote').addEventListener('click', e => {
      upvote(videoId);
    });
    item.querySelector('.downvote').addEventListener('click', e => {
      downvote(videoId);
    });
    item.querySelector('.unvote').addEventListener('click', e => {
      unvote(videoId);
    });
    item.querySelector('.loginToVote').addEventListener('click', e => {
      alert('log in to vote');
    });
  });
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
  window.createNewList = createNewList;
  window.newList = newList;
  addHandlers();
  window.addEventListener('popstate', loadSelectedItem);
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
  loadSelectedItem({});
};

const loadSelectedItem = event => {
  console.log(
    'location: ' + document.location + ', state: ' + JSON.stringify(event.state)
  );

  let id = '';
  if (event && event.state && event.state.video) {
    id = event.state.video;
  } else {
    const u = new URL(document.location);
    if (u.hash.startsWith('#v=')) {
      id = u.hash.substring(3);
    }
  }
  state.selectedItem = id;
  show(id);
};
// Listen on page load:
window.addEventListener('load', start);
