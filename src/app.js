"use strict";

let state = {
  items: [
    { id: "X7hFERntlog", title: "Fearless Org" }, 
    { id: "d_HHnEROy_w", title: "Stop managing" }, 
    { id: "BCkCvay4-DQ", title: "Foos" }
  ],
};
const items = () => {
  return state.items.map( i => `
    <li onclick="show('${i.id}')">
      <div>
        <div>
          ${i.title} 
        </div>
        <div>
          <img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote(${i.id})" />
          <img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote(${i.id})" />
        </div>
      </div>
    </li>` ).join('');
};

const show = (i) => {
  document.getElementById('player').setAttribute('src', 'https://www.youtube.com/embed/' + i);
  const item = state.items.find(item => item.id === i);
  document.getElementById('title').innerHTML = item.title;
};

const upvote = (i) => {
}

const downvote = (i) => {
}

const add = () => {
  const id = document.getElementById('addbox').value;
  fetch(`/yt/data/${id}`, { 
      method: 'get',
    })
  .then(function(response) {
    var json = response.json();
    state.items.push({id: id, title: json.title});
    reflop();
  })
  .then(function(myJson) {
    console.log(JSON.stringify(myJson));
  });
};

// The router code. Takes a URL, checks against the list of supported routes and then renders the corresponding content page.
const reflop = async () => {
  const content = document.getElementById('page_content')
  const vidList = document.getElementById('videoList');
  
  if (content) {
    content.innerHTML = '<h2>app</h2>';
  }
  if (vidList) {
    vidList.innerHTML = items();
  }
}

// Listen on page load:
window.addEventListener('load', reflop);
