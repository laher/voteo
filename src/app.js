'use strict';

let state = {
  personId: null,
  items: [
    { id: 'X7hFERntlog', title: 'Fearless Org' },
    { id: 'd_HHnEROy_w', title: 'Stop managing' },
    { id: 'BCkCvay4-DQ', title: 'Foos' },
  ],
  votes: [{ videoId: 'x7hFERntlog', personId: 'am', up: true }],
  selectedItem: 'X7hFERntlog',
};

const countVotes = id => {
  return (
    state.votes.filter(vote => vote.videoId == id && vote.up).length -
    state.votes.filter(vote => vote.videoId == id && !vote.up).length
  );
};

const haveIUpvoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && vote.up && vote.personId == state.personId
    ).length > 0
  );
};

const haveIDownvoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && !vote.up && vote.personId == state.personId
    ).length > 0
  );
};

const haveIVoted = id => {
  return (
    state.votes.filter(
      vote => vote.videoId == id && vote.personId == state.personId
    ).length > 0
  );
};

const items = id => {
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
        (haveIUpvoted(i.id)
          ? `<img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('${
              i.id
            }')" />`
          : `
    <img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote('${
      i.id
    }')" />`) +
        (haveIDownvoted(i.id)
          ? `<img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('${
              i.id
            }')" />`
          : `<img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote('${
              i.id
            }')" /> `) +
        `</div>
    </div></li>`
    )
    .join('');
};

const setSelectedItem = i => {
  state.selectedItem = i;
  show(i);
};

const show = i => {
  if (!i) {
    return;
  }
  document
    .getElementById('player')
    .setAttribute('src', 'https://www.youtube.com/embed/' + i);
  const item = state.items.find(item => item.id === i);
  let title = 'preview';
  if (item) {
    title = item.title;
  }
  document.getElementById('title').innerHTML = title;
};

const upvote = i => {
  const item = state.items.find(item => item.id === i);
  let vote = state.votes.find(
    vote => vote.personId == state.personId && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: state.personId, videoId: i };
  }
  vote.up = true;
  postVote(vote);
};

const unvote = i => {
  const vote = state.votes.find(
    vote => vote.personId == state.personId && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: state.personId, videoId: i };
  }
  deleteVote(vote);
};

const downvote = i => {
  const item = state.items.find(item => item.id === i);
  const vote = state.votes.find(
    vote => vote.personId == state.personId && vote.videoId == i
  );
  if (!vote) {
    vote = { personId: state.personId, videoId: i };
  }
  vote.up = false;
  postVote(vote);
};

const getVideos = () => {
  console.log('token', signIn.tokenManager.get('access_token'));
  fetch(`/videos`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization:
        'Bearer ' + signIn.tokenManager.get('access_token').accessToken,
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

const getVotes = () => {
  console.log('token', signIn.tokenManager.get('access_token'));
  fetch(`/vote`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization:
        'Bearer ' + signIn.tokenManager.get('access_token').accessToken,
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
  fetch(`/videos`, {
    method: 'PUT',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization:
        'Bearer ' + signIn.tokenManager.get('access_token').accessToken,
    },
    redirect: 'follow', // manual, *follow, error
    referrer: 'no-referrer', // no-referrer, *client
    body: JSON.stringify(state.items), // body data type must match "Content-Type" header
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

const postVote = vote => {
  fetch(`/vote`, {
    method: 'POST',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization:
        'Bearer ' + signIn.tokenManager.get('access_token').accessToken,
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
      Authorization:
        'Bearer ' + signIn.tokenManager.get('access_token').accessToken,
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

const preview = () => {
  const id = cleanInput(document.getElementById('addbox').value);
  if (id) {
    show(id);
  }
};

const add = () => {
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

const reflop = async () => {
  var t0 = performance.now();
  document.getElementById('logged-in').style.display = 'block';
  const vidList = document.getElementById('videoList');
  if (vidList) {
    vidList.innerHTML = items();
    document.getElementById('videoCount').innerHTML = state.items.length;
  }
  show(state.selectedItem);

  var t1 = performance.now();
  console.log('Call to reflop took ' + (t1 - t0) + ' milliseconds.');
};

const showSignOut = () => {
  console.log('show signout');
  signIn.hide();
  document.getElementById('name').innerHTML = state.personId;
  document.getElementById('logged-in').style.display = 'block';
};

const showSignIn = () => {
  console.log('show signin');
  document.getElementById('logged-in').style.display = 'none';
  signIn.renderEl({ el: '#widget-container' }, res => {
    if (res.status === 'SUCCESS') {
      console.log('signin success', res);
      signIn.tokenManager.add('id_token', res[0]);
      signIn.tokenManager.add('access_token', res[1]);
      console.log('signin success. tokenManager:', signIn.tokenManager);
      state.personId = res[0].claims.email;
      showSignOut();
      getVideos();
      getVotes();
    }
  });
};

let signIn = null;
const start = () => {
  document.getElementById('sign-out').addEventListener('click', event => {
    event.preventDefault();

    console.log('signout clicked');
    signIn.session.close(err => {
      if (err) {
        alert(`Error: ${err}`);
      }
      showSignIn();
    });
  });
  signIn = new OktaSignIn({
    baseUrl: 'https://dev-343286.okta.com',
    clientId: '0oabsbm6ga3Sy1tIf356',
    redirectUri: 'http://localhost:3000/auth/callback/login',
    authParams: {
      issuer: 'default',
      responseType: ['id_token', 'token'],
    },
    idps: [{ type: 'GOOGLE', id: '0oack0yq3VXGwi171356' }],
  });
  init();
};

const init = () => {
  signIn.session.get(async res => {
    if (res.status === 'ACTIVE') {
      console.log('login already active', res);
      state.personId = res.login;
      showSignOut();
      getVideos();
      getVotes();
    } else {
      console.log('not signed in');
      showSignIn();
    }
  });
};

// Listen on page load:
window.addEventListener('load', start);
