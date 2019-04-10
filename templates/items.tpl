{{ $personID := .PersonID }}
{{ range $index, $i := .VideoList.Videos }}
<li>
  <div class="video-item" x="{{rand}}" data-id="{{$i.ID}}">
    <div class="title">
      {{$i.Title}} 
    </div>
    <div>

      <span class="loadingSpan" style="display:none">
        <img src="https://img.icons8.com/material-two-tone/24/000000/loading.png" style="visibility:hidden">
        <img src="https://img.icons8.com/material-two-tone/24/000000/loading.png">
      </span>
      <div class="votingSpan" style="padding-right: 5px">
        <em>
          <abbr title='Current vote count'>
            {{ if (gt (countVotes $i.ID) 0) }}+{{ countVotes $i.ID }}
            {{ else if (eq (countVotes $i.ID) 0) }}&nbsp;0
            {{ else if (lt (countVotes $i.ID) 0) }}{{ countVotes $i.ID }}{{ end }}
          </abbr>
        </em>
      </div>
      <span class="votingSpan" {{ if not $personID }}style="display:none"{{end}}>
        <span class="unvoteSpan" {{ if not (haveIVoted $i.ID) }}style="display:none"{{ end }}>
          <img src="https://img.icons8.com/material-two-tone/24/000000/undo.png" style="visibility:hidden">
          <img class="unvote" src="https://img.icons8.com/material-two-tone/24/000000/undo.png" title="Undo Vote" />
        </span>
        <span class="voteSpan" {{ if haveIVoted $i.ID }}style="display:none"{{ end }}>
          <img class="upvote" src="https://img.icons8.com/material-two-tone/24/000000/like.png" title="Upvote" />
          <img class="downvote" src="https://img.icons8.com/material-two-tone/24/000000/dislike.png" title="Downvote" />
        </span>
      </span>
      <abbr {{ if $personID }}style="display:none"{{ end }} title='log in to vote'><img class="loginToVote" src="https://img.icons8.com/small/24/000000/thumbs-up-down.png"></abbr>
    </div>
  </div>
</li>
{{ end }}
