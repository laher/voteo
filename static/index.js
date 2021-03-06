'use strict';

import { initAuth, cleanInput } from './app.js';
import { putNewVideoList, getMetadata } from './api.js';

let state = {
  videoLists: [],
};

const setVideoLists = videoLists => {
  if (videoLists) {
    state.videoLists = videoLists;
  }
};

const newList = () => {
  console.log('new list');
  const id = cleanInput(document.getElementById('addbox').value);
  const videoListTitle = document.getElementById('video_list_title').value;
  getMetadata(id, json => {
    let videoTitle = json.title;
    if (videoTitle.length > 30) {
      videoTitle = videoTitle.substring(0, 30) + ' ...';
    }
    state.selectedItem = id;
    putNewVideoList(
      videoListTitle,
      { sourceId: id, title: videoTitle },
      json => {
        const videoListId = json.id;
        console.log('created video list', json);
        window.location = '/video-list?id=' + videoListId;
      }
    );
  });
};

const reflopVideoLists = () => {
  console.log(state.videoLists);
  document.getElementById('videoLists').innerHTML = state.videoLists
    .map(
      i =>
        `<li><div class="videoList-item" data-id="${i.id}">List ${i.id} (${
          i.videos.length
        } videos)</div></li>`
    )
    .join('');

  const nodeList = document.querySelectorAll('div.videoList-item');

  nodeList.forEach(item => {
    const videoListId = item.getAttribute('data-id');
    item.addEventListener('click', e => {
      console.log('going to video list', videoListId);
      window.location = '/video-list?id=' + videoListId;
    });
  });
};

export const pageInit = videoLists => {
  console.log('page init with data: ', videoLists);
  initAuth(() => {});
  setVideoLists(videoLists);
  reflopVideoLists();
  console.log('before');
  document.getElementById('add').addEventListener('click', newList);
  console.log('page init complete');
};
