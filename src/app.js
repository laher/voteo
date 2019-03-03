"use strict";

let state = {
  items: [
    { id: "X7hFERntlog", title: "Fearless Org", votes: 0 }, 
    { id: "d_HHnEROy_w", title: "Stop managing", votes: -10 }, 
    { id: "BCkCvay4-DQ", title: "Foos", votes: 1 }
  ],
  selectedItem: "X7hFERntlog"
};
const items = (id) => {
  return state.items.sort((a, b) => b.votes - a.votes ).map( i => `
    <li >
    <div>
    <div onclick="setSelectedItem('${i.id}')">
    ${i.title} (${i.votes} votes)
    </div>
    <div>
    <img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote('${i.id}')" />
    <img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote('${i.id}')" />
    </div>
    </div>
    </li>` ).join('');
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
  const vidList = document.getElementById('videoList');
  if (vidList) {
    vidList.innerHTML = items();
  }
  show(state.selectedItem);
}

// Listen on page load:
window.addEventListener('load', getVideos);
