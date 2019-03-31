{{ $personID := .PersonID }}
{{ range $index, $i := .Items }}
<li><div>
    <div onclick="setSelectedItem('{{$i.ID}}')">{{$i.Title}} </div>

    <div>
      <div style="padding-right: 5px">
        <em>
          <abbr title='Current vote count'>
            {{ if (gt (countVotes $i.ID) 0) }}+{{ countVotes $i.ID }}
            {{ else if (eq (countVotes $i.ID) 0) }}&nbsp;0
            {{ else if (lt (countVotes $i.ID) 0) }}{{ countVotes $i.ID }}{{ end }}
          </abbr>
        </em>
      </div>
      {{- if $personID }}
      {{ if haveIVoted $i.ID }}
      <img src="https://img.icons8.com/material-two-tone/24/000000/undo.png" style="visibility:hidden">
      <img src="https://img.icons8.com/material-two-tone/24/000000/undo.png" onclick="unvote('{{$i.ID}}')" title="Undo Vote" />
      {{ else }}
      <img src="https://img.icons8.com/material-two-tone/24/000000/like.png" onclick="upvote('{{$i.ID}}')" title="Upvote" />
      <img src="https://img.icons8.com/material-two-tone/24/000000/dislike.png" onclick="downvote('{{$i.ID}}')" title="Downvote" />
      {{ end }}
      {{ else }}
      <abbr title='log in to vote'><img src="https://img.icons8.com/small/24/000000/thumbs-up-down.png" onclick="alert('log in to vote')"></abbr>
      {{ end }}
    </div>
  </div>

</li>
{{ end }}
