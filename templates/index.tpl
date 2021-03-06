{{ template "top.tpl" . }}
        <div class="card">
          <div class="top">
                <p><em>Videos are a great way to begin a conversation - TED Talks, how-tos, movie trailers, debates ...</em></p>
                <p>With <b>Voteo</b> you can collect lists of videos, and vote on which to watch together.</p>
                <p>To start a new list, simply register and add a video below. Later, you can share your list with your friends and get voting.</p>
          </div>
          <div class="left">
            <h3>Add a video list</h3>
            <em>Set a title, and optionally a video, to begin a new list</em>
            <div class="video-list">
              <div class="video-list-title">
                <label for="video_list_title">Title</label>
                <input
                  id="video_list_title"
                  type="text"
                  placeholder="Set a video list name ..."
                />
              </div>
              <div class="new-item-form">
                <label for="addbox">Video ID</label>
                <input
                  id="addbox"
                  type="text"
                  placeholder="Drop a youtube video id or url here ..."
                />
              </div>
              <div class="video-list-title">
                <button id="add">
                  New video list
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
            <ul id="videoLists" class="list">
            </ul>
          </div>
        </div>

<script>
  // note: html/template doesn't properly escape '<script type="modules">' in the same way as '<script>'
window.addEventListener('load', () => {
  pageInit({{ .VideoLists }});
});
</script>
{{ template "bottom.tpl" }}
