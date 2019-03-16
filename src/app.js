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
  idToken: null,
  accessToken: null,
  oktaSignIn: null,
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
        (state.personId
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

const setSelectedItem = i => {
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
  fetch(`/videos`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization: 'Bearer ' + state.accessToken,
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
  fetch(`/vote`, {
    method: 'get',
    cache: 'no-cache',
    headers: {
      Authorization: 'Bearer ' + state.accessToken,
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
  if (!state.personId) {
    alert('Please log in or register to add videos');
    return;
  }
  fetch(`/videos`, {
    method: 'PUT',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + state.accessToken,
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
      Authorization: 'Bearer ' + state.accessToken,
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
      Authorization: 'Bearer ' + state.accessToken,
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
  hideOkta();
  document.getElementById('app-container').style.display = 'block';
  document.getElementById('name').innerHTML = state.personId;
  document.getElementById('logged-in').style.display = 'flex';
  document.getElementById('logged-out').style.display = 'none';
};

const showSignInButton = () => {
  console.log('show signin button');
  hideOkta();
  document.getElementById('app-container').style.display = 'block';
  document.getElementById('name').innerHTML = '';
  document.getElementById('logged-in').style.display = 'none';
  document.getElementById('logged-out').style.display = 'flex';
};

const showSignInModal = () => {
  console.log('show signin modal');
  if (state.oktaSignIn) {
    document.getElementById('app-container').style.display = 'none';

    state.oktaSignIn.show();
  } else {
    console.log('oops: non-okta signin not implemented');
  }
};

const start = () => {
  fetch(`/auth/settings`, {
    method: 'get',
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(JSON.stringify(json));
      if (json['type'] == 'okta') {
        state.oktaSignIn = new OktaSignIn(json.okta);
        doOkta();
      } else {
        state.personId = '';
        // assume it's no-login
        getVideos();
        getVotes();
      }
    });
};

// Listen on page load:
window.addEventListener('load', start);
