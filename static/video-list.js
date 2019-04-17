import { initOkta, getPersonId } from './auth-okta.js';
import {
  getVideoList,
  addVideoToList,
  postVote,
  deleteVote,
  getMetadata,
} from './api.js';
import { initAuth, cleanInput } from './app.js';

let state = {
  videoListId: 0,
  videoList: null,
  selectedItem: '',
};

const myVideoList = () => {
  if (state.videoList != null) {
    return state.videoList;
  }
  return { id: state.videoListId };
};

export const countVotes = id => {
  return (
    state.videoList.votes.filter(vote => vote.videoId == id && vote.up).length -
    state.videoList.votes.filter(vote => vote.videoId == id && !vote.up).length
  );
};

export const haveIUpvoted = id => {
  return (
    state.videoList.votes.filter(
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
    state.videoList.votes.filter(
      vote => vote.videoId == id && vote.personId == getPersonId()
    ).length > 0
  );
};

export const setSelectedItem = i => {
  history.pushState({ video: i }, 'title 1', '#v=' + i);
  console.log('history', history.state);
  loadSelectedItem({ state: history.state });
};

export const preview = () => {
  const id = cleanInput(document.getElementById('addbox').value);
  if (id) {
    show(id);
  }
};

export const add = () => {
  console.log('add');
  const videoId = cleanInput(document.getElementById('addbox').value);
  getMetadata(videoId, json => {
    let title = json.title;
    if (title.length > 30) {
      title = title.substring(0, 30) + ' ...';
    }
    const vid = { id: videoId, title: title };
    state.videoList.videos.push(vid);
    console.log('videoList after push', state.videoList);
    state.selectedItem = videoId;
    reflop();
    addVideoToList(vid, state.videoList.id, () => {});
  });
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
  const item = state.videoList.videos.find(item => item.id === i);
  let title = 'preview';
  if (item) {
    title = item.title;
  }
  document.getElementById('title').innerHTML = title;
};

export const upvote = i => {
  const item = state.videoList.videos.find(item => item.id === i);
  let vote = state.videoList.votes.find(
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
  postVote(state.videoListId, vote, json => {
    state.videoList.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
};

export const unvote = i => {
  let vote = state.videoList.votes.find(
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
  deleteVote(state.videoListId, vote, json => {
    state.videoList.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
};

export const downvote = i => {
  const item = state.videoList.videos.find(item => item.id === i);
  let vote = state.videoList.votes.find(
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
  postVote(state.videoListId, vote, json => {
    state.videoList.votes = json.Votes;
    updateVidList(json.ItemsHTML);
  });
};

const updateVidList = () => {
  document.getElementById('videoList').innerHTML = renderVideos();
  addHandlers();
};

const reflop = async () => {
  var t0 = performance.now();
  const vidList = document.getElementById('videoList');
  if (vidList) {
    updateVidList();
    //    getVideoList(state.videoListId, updateVidList);
    document.getElementById('videoCount').innerHTML =
      state.videoList.videos.length;
  }
  show(state.selectedItem);
  var t1 = performance.now();
  console.log('Call to reflop took ' + (t1 - t0) + ' milliseconds.');
  console.log('videoList after reflop', state.videoList);
};

const renderVideos = () => {
  if (!state.videoList.videos) {
    console.log('state bad: ', state);
  }
  return state.videoList.videos
    .map(
      i => `<li>
          <div class="video-item" data-source-id="${i.sourceId}" data-id="${
        i.id
      }">
            <div class="title">
              ${i.title} 
            </div>
            <div>
              <span class="loadingSpan" style="display:none">
                <img src="https://img.icons8.com/material-two-tone/24/000000/loading.png" style="visibility:hidden">
                <img src="https://img.icons8.com/material-two-tone/24/000000/loading.png">
              </span>
              <div class="votingSpan" style="padding-right: 5px">
                <em>
                  <abbr title='Current vote count'>
                    ${
                      countVotes(i.id) > 0
                        ? '+'
                        : countVotes(i.id) < 0
                        ? '-'
                        : '&nbsp;'
                    }${countVotes(i.ID)}
                  </abbr>
                </em>
              </div>
              <span class="votingSpan" ${
                getPersonId() ? '' : 'style="display:none"'
              }>
                <span class="unvoteSpan" ${
                  haveIVoted(i.id) ? '' : 'style="display:none"'
                }>
                  <img src="https://img.icons8.com/material-two-tone/24/000000/undo.png" style="visibility:hidden">
                  ${imgButton('unvote', 'undo', 'Undo Vote')}
                </span>
                <span class="voteSpan" ${
                  haveIVoted(i.id) ? 'style="display:none"' : ''
                }>
                  ${imgButton('upvote', 'like', 'Upvote')}
                  ${imgButton('downvote', 'dislike', 'Downvote')}
                </span>
              </span>
              <abbr ${
                getPersonId() ? 'style="display:none"' : ''
              } title='log in to vote'>
                  ${imgButton('loginToVote', 'novel', 'log in to vote')}
            </div>
          </div>
        </li>`
    )
    .join('');
};

const imgButton = (cls, img, title) =>
  `<img class="${cls}" src="https://img.icons8.com/material/24/000000/${img}.png" title="${title}" />`;

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
const addHandlers = () => {
  console.log('set up selectors / event listeners');
  const nodeList = document.querySelectorAll('div.video-item');
  //const nodeList = document.querySelectorAll('div.video-item');
  console.log('found ', nodeList.length, ' items');
  Array.from(nodeList).forEach(item => {
    const videoId = parseInt(item.getAttribute('data-id'), 10);
    const sourceId = item.getAttribute('data-source-id');
    item.querySelector('.title').addEventListener('click', e => {
      setSelectedItem(sourceId);
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

export const pageInit = videoList => {
  console.log('pageInit[video-list] with data: ', videoList);
  if (videoList.votes == null) {
    videoList.votes = [];
  }
  if (videoList.videos == null) {
    videoList.videos = [];
  }
  state.videoList = videoList;
  state.videoListId = new URLSearchParams(location.search).get('id');
  initAuth(reflop);
  window.setSelectedItem = setSelectedItem;
  document.getElementById('addbox').addEventListener('change', preview);
  document.getElementById('add').addEventListener('click', add);
  addHandlers();
  window.addEventListener('popstate', loadSelectedItem);
  loadSelectedItem({});
};
