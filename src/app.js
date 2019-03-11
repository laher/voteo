"use strict";

let state = {
  userId: null,
  items: [
    { id: "X7hFERntlog", title: "Fearless Org", votes: 0 }, 
    { id: "d_HHnEROy_w", title: "Stop managing", votes: -10 }, 
    { id: "BCkCvay4-DQ", title: "Foos", votes: 1 }
  ],
  selectedItem: "X7hFERntlog"
};
const items = (id) => {
  return state.items.sort((a, b) => b.votes - a.votes ).map( i => `
    <li><div>
    <div onclick="setSelectedItem('${i.id}')">${i.title} (${i.votes} votes)</div>
    <div>
    <img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote('${i.id}')" />
    <img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote('${i.id}')" />
    </div>
    </div></li>` ).join('');
};

const setSelectedItem = (i) => {
  state.selectedItem = i;
  show(i);
};

const show = (i) => {
  if (!i) {
    return;
  }
  document.getElementById('player').setAttribute('src', 'https://www.youtube.com/embed/' + i);
  const item = state.items.find(item => item.id === i);
  let title = 'preview';
  if (item) {
    title = item.title;
  }
  document.getElementById('title').innerHTML = title;
};

const upvote = (i) => {
  const item = state.items.find(item => item.id === i);
  item.votes++;
  putVideos();
}

const downvote = (i) => {
  const item = state.items.find(item => item.id === i);
  item.votes--;
  putVideos();
}

const getVideos = () => {
  fetch(`/videos`, { 
    method: 'get',
    cache: "no-cache",
  })
    .then(function(response) {
      console.log(response);
      return response.json();
    })
    .then(function(json) {
      console.log(json);
      state.items = json;
      //push({id: id, title: json.title});
      reflop();
    });
}

const cleanInput = (input) => {
  const pref = "https://www.youtube.com/watch?v=";
  if (input.startsWith(pref)) {
    return input.substring(pref.length);
  }
  return input;
}

const putVideos = () => {
  fetch(`/videos`, { 
    method: 'PUT',
    cache: "no-cache",
    headers: {
      "Content-Type": "application/json",
    },
    redirect: "follow", // manual, *follow, error
    referrer: "no-referrer", // no-referrer, *client
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
        title = title.substring(0, 30) + " ..."; 
      }
      state.items.push({id: id, title: title});
      state.selectedItem = id;
      putVideos();
    });
};

// The router code. Takes a URL, checks against the list of supported routes and then renders the corresponding content page.
const reflop = async () => {
  document.getElementById('logged-in').style.display = 'block';
  const vidList = document.getElementById('videoList');
  if (vidList) {
    vidList.innerHTML = items();
    document.getElementById('videoCount').innerHTML = state.items.length;
  }
  show(state.selectedItem);
}

const showSignOut = () => {
  console.log("show signout");
  signIn.hide();
  document.getElementById('logged-in').style.display = 'block';
}

const showSignIn = () => {
  console.log("show signin");
  document.getElementById('logged-in').style.display = 'none';
  signIn.renderEl({ el: '#widget-container' }, (res) => {
    if (res.status === 'SUCCESS') {
      console.log("signin success", res);
      signIn.tokenManager.add('id_token', res[0]);
      signIn.tokenManager.add('access_token', res[1]);
      console.log("signin success. tokenManager:", signIn.tokenManager);
      document.getElementById('name').innerHTML = res[0].claims.email;
      showSignOut();
      getVideos();
    }
  });
}

let signIn = null;
const start = () => {
  document.getElementById('sign-out').addEventListener('click', (event) => {
    event.preventDefault();

    console.log("signout clicked");
    signIn.session.close((err) => {
      if (err) {
        alert(`Error: ${err}`)
      }
      showSignIn()
    })
  });
  signIn = new OktaSignIn({
    baseUrl: 'https://dev-343286.okta.com',
    clientId: '0oabsbm6ga3Sy1tIf356',
    redirectUri: 'http://localhost:3000/auth/callback/login',
    authParams: {
      issuer: 'default',
      responseType: ['id_token','token']
    }
  });
  init();
}

const init = () => {
  signIn.session.get(async (res) => {
      if (res.status === 'ACTIVE') {
        getVideos();
        console.log('login already active', res);
        document.getElementById('name').innerHTML = res.login;

        console.log('tokenManager', signIn.tokenManager);

        showSignOut();
      } else {
        console.log('not signed in');
        showSignIn();
      }
    })

}

// Listen on page load:
window.addEventListener('load', start);
