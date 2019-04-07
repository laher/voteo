{{ template "top.tpl" . }}
        <div class="card">
          <div class="top">
                <p><em>Videos are a great way to begin a conversation - TED Talks, how-tos, movie trailers, debates ...</em></p>
                <p>With <b>Voteo</b> you can collect lists of videos, and vote on which to watch together.</p>
                <p>To start a new list, simply register and add a video below. Later, you can share your list with your friends and get voting.</p>
          </div>
          <div class="left">
            <h3>Add a video</h3>
            <em>Add a video to begin a new list</em>
            <div class="video-list">
              <div class="new-item-form">
                <input
                  id="addbox"
                  type="text"
                  placeholder="Drop a youtube video id or url here ..."
                  onchange="preview()"
                />
                <button id="add" onclick="createNewList()">
                  <img
                    src="https://img.icons8.com/material-two-tone/24/000000/plus.png"
                    title="Add the video"
                  />
                </button>
              </div>
            </div>
            <div class="info">
            </div>
          </div>
          <div class="right">
            <h3>My Lists</h3>
            {{ template "video-lists.tpl" . }}
          </div>
        </div>
{{ template "bottom.tpl" }}
