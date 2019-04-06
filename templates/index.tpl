        <div class="card">
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
            <h3>Watch</h3>
            <p id="title"></p>
            <div>
              <iframe
                id="player"
                width="420"
                height="315"
                src=""
                allowfullscreen
              />
            </div>
          </div>
        </div>
